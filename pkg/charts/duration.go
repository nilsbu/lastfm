package charts

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Interval is a time span. Begin is the first moment, Before is the moment
// immediately after the interval ends.
type Interval struct {
	Begin  rsrc.Day
	Before rsrc.Day
}

type intervalIterator interface {
	HasNext() bool
	Next() Interval
}

type iIterator struct {
	next   Interval
	before rsrc.Day
	step   int
}

type dayIterator struct {
	iIterator
}

type monthIterator struct {
	iIterator
}

type yearIterator struct {
	iIterator
}

func newYearIterator(
	step int,
	from rsrc.Day,
	before rsrc.Day) intervalIterator {
	return &yearIterator{iIterator{
		next:   yearPeriod(from, step),
		before: before,
		step:   step,
	}}
}

func newMonthIterator(
	step int,
	from rsrc.Day,
	before rsrc.Day) intervalIterator {
	return &monthIterator{iIterator{
		next:   monthPeriod(from, step),
		before: before,
		step:   step,
	}}
}

func newDayIterator(
	step int,
	from rsrc.Day,
	before rsrc.Day) intervalIterator {
	return &dayIterator{iIterator{
		next:   dayPeriod(from, step),
		before: before,
		step:   step,
	}}
}

func (ii *iIterator) HasNext() bool {
	return ii.next.Begin.Midnight() < ii.before.Midnight()
}

func (ii *dayIterator) Next() Interval {
	next := ii.next
	ii.next = dayPeriod(ii.next.Before, ii.step)

	return next
}

func (ii *monthIterator) Next() Interval {
	next := ii.next
	ii.next = monthPeriod(ii.next.Before, ii.step)

	return next
}

func (ii *yearIterator) Next() Interval {
	next := ii.next
	ii.next = yearPeriod(ii.next.Before, ii.step)

	return next
}

func dayPeriod(day rsrc.Day, step int) Interval {
	t := day.Time()
	y, m, d := t.Year(), t.Month(), t.Day()
	return Interval{
		Begin:  rsrc.DayFromTime(time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)),
		Before: rsrc.DayFromTime(time.Date(y, time.Month(m), d+step, 0, 0, 0, 0, time.UTC)),
	}
}

func monthPeriod(day rsrc.Day, step int) Interval {
	t := day.Time()
	y := t.Year()
	m := int(int(t.Month())-1)/step*step + 1
	return Interval{
		Begin:  rsrc.DayFromTime(time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)),
		Before: rsrc.DayFromTime(time.Date(y, time.Month(m+step), 1, 0, 0, 0, 0, time.UTC)),
	}
}

func yearPeriod(day rsrc.Day, step int) Interval {
	t := day.Time()
	y := int(t.Year()/step) * step
	return Interval{
		Begin:  rsrc.DayFromTime(time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)),
		Before: rsrc.DayFromTime(time.Date(y+step, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}
}

// Period parses a string describing a period and returns the corresponding
// interval. The descriptor is either a year in the format 'yyyy' or a month
// in the format 'yyyy-MM'.
func Period(descr string) (Interval, error) {
	switch len(descr) {
	case 4:
		begin, err := time.Parse("2006", descr)
		if err != nil {
			return Interval{}, err
		}
		return yearPeriod(rsrc.DayFromTime(begin), 1), nil
	case 7:
		begin, err := time.Parse("2006-01", descr)
		if err != nil {
			return Interval{}, err
		}
		return monthPeriod(rsrc.DayFromTime(begin), 1), nil
	default:
		return Interval{}, fmt.Errorf("interval format '%v' not supported", descr)
	}
}

// Interval returns a Column that sums an interval of the charts. The charts
// have to be a sum.
func (c Charts) Interval(i Interval, registered rsrc.Day) Column {
	size := c.Len()

	from := Index(i.Begin, registered)
	to := Index(i.Before, registered)
	if to < 0 {
		return Column{}
	} else if to >= size {
		to = size - 1
	}

	column := Column{}

	if from >= size {
		return Column{}
	} else if from < 0 {
		for i, name := range c.Keys {
			column = append(column, Score{name.String(), c.Values[i][to]})
		}
	} else {
		for i, name := range c.Keys {
			column = append(column, Score{name.String(), c.Values[i][to] - c.Values[i][from]})
		}
	}
	sort.Sort(column)
	return column
}

func (c Charts) Intervals(intervals Intervals, registered rsrc.Day) Charts {
	icharts := []map[string]float64{}
	for i := 0; i < intervals.Len(); i++ {
		col := c.Interval(intervals.At(i), registered)

		if len(col) == 0 {
			continue
		}

		cha := map[string]float64{}
		for _, x := range col {
			cha[x.Name] = x.Score
		}

		icharts = append(icharts, cha)
	}

	ncha := CompileArtists(icharts, registered)
	ncha.Headers = intervals
	return ncha
}

// Index calculates an column index based on registration date and searcherd
// date.
func Index(t rsrc.Day, registered rsrc.Day) int { // TODO is obsolete
	return int((t.Midnight()-registered.Midnight())/86400 - 1)
}

// TODO change name
func ToIntervals(
	descr string, begin, end rsrc.Day,
) (Intervals, error) {

	re := regexp.MustCompile("^\\d*[yMd]$")
	if !re.MatchString(descr) {
		return nil, fmt.Errorf("interval descriptor '%v' invalid", descr)
	}

	n, err := strconv.Atoi(descr[:len(descr)-1])
	if err != nil {
		n = 1
	}

	// var ii intervalIterator
	switch descr[len(descr)-1] {
	case 'y':
		return Years(begin, end, n), nil
	case 'M':
		return Months(begin, end, n), nil
	default:
		return MultiDays(begin, end, n), nil
	}
}
