package charts

import (
	"fmt"
	"sort"
	"time"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Interval is a time span. Begin is the first moment, Before is the moment
// immediately after the interval ends.
type Interval struct {
	Begin  rsrc.Day
	Before rsrc.Day
}

type Step int

const (
	Day Step = iota
	Month
	Year
)

type intervalIterator interface {
	HasNext() bool
	Next() Interval
}

type iIterator struct {
	next   Interval
	before rsrc.Day
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

func newIntervalIterator(
	step Step,
	from rsrc.Day,
	before rsrc.Day) intervalIterator {
	switch step {
	case Day:
		return &dayIterator{iIterator{
			next:   dayPeriod(from),
			before: before,
		}}
	case Month:
		return &monthIterator{iIterator{
			next:   monthPeriod(from),
			before: before,
		}}
	default:
		return &yearIterator{iIterator{
			next:   yearPeriod(from),
			before: before,
		}}
	}
}

func (ii *iIterator) HasNext() bool {
	return ii.next.Begin.Midnight() < ii.before.Midnight()
}

func (ii *dayIterator) Next() Interval {
	next := ii.next
	ii.next = dayPeriod(ii.next.Before)

	return next
}

func (ii *monthIterator) Next() Interval {
	next := ii.next
	ii.next = monthPeriod(ii.next.Before)

	return next
}

func (ii *yearIterator) Next() Interval {
	next := ii.next
	ii.next = yearPeriod(ii.next.Before)

	return next
}

func dayPeriod(day rsrc.Day) Interval {
	t := day.Time()
	y, m, d := t.Year(), t.Month(), t.Day()
	return Interval{
		Begin:  rsrc.ToDay(time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC).Unix()),
		Before: rsrc.ToDay(time.Date(y, time.Month(m), d+1, 0, 0, 0, 0, time.UTC).Unix()),
	}
}

func monthPeriod(day rsrc.Day) Interval {
	t := day.Time()
	y, m := t.Year(), t.Month()
	return Interval{
		Begin:  rsrc.ToDay(time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC).Unix()),
		Before: rsrc.ToDay(time.Date(y, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC).Unix()),
	}
}

func yearPeriod(day rsrc.Day) Interval {
	t := day.Time()
	y := t.Year()
	return Interval{
		Begin:  rsrc.ToDay(time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()),
		Before: rsrc.ToDay(time.Date(y+1, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()),
	}
}

func Period(descr string) (Interval, error) {
	switch len(descr) {
	case 4:
		begin, err := time.Parse("2006", descr)
		if err != nil {
			return Interval{}, err
		}
		return yearPeriod(rsrc.ToDay(begin.Unix())), nil
	case 7:
		begin, err := time.Parse("2006-01", descr)
		if err != nil {
			return Interval{}, err
		}
		return monthPeriod(rsrc.ToDay(begin.Unix())), nil
	default:
		return Interval{}, fmt.Errorf("interval format '%v' not supported", descr)
	}
}

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
		for name, line := range c {
			column = append(column, Score{name, line[to]})
		}
	} else {
		for name, line := range c {
			column = append(column, Score{name, line[to] - line[from]})
		}
	}
	sort.Sort(column)
	return column
}

func (c Charts) Intervals(intervals []Interval, registered rsrc.Day) Charts {
	icharts := []Charts{}
	for _, i := range intervals {
		col := c.Interval(i, registered)

		if len(col) == 0 {
			continue
		}

		cha := Charts{}
		for _, x := range col {
			cha[x.Name] = []float64{x.Score}
		}

		icharts = append(icharts, cha)
	}

	return Compile(icharts)
}

func Index(t rsrc.Day, registered rsrc.Day) int {
	return int((t.Midnight()-registered.Midnight())/86400 - 1)
}

// TODO use in table formatting & change name
func (c Charts) ToIntervals(step Step, registered rsrc.Day) []Interval {
	reg := registered.Midnight()
	ii := newIntervalIterator(
		step,
		registered,
		rsrc.ToDay(reg+int64(86400*c.Len())))

	intervals := []Interval{}
	for ii.HasNext() {
		current := ii.Next()

		intervals = append(intervals, current)
	}

	return intervals
}
