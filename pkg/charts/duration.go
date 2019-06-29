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

// TODO mix between int64 and time.Time
type iIterator struct {
	step     Step
	interval Interval
	before   int64
}

func newIntervalIterator(step Step, from rsrc.Day, before int64) intervalIterator {
	var interval Interval

	switch step {
	case Day:
		interval = dayPeriod(from)
	case Month:
		interval = monthPeriod(from)
	default:
		interval = yearPeriod(from)
	}

	return &iIterator{
		step:     step,
		interval: interval,
		before:   before,
	}
}

func (ii *iIterator) HasNext() bool {
	return ii.interval.Begin.Midnight() < ii.before
}

func (ii *iIterator) Next() Interval {
	interval := ii.interval

	before := ii.interval.Before
	switch ii.step {
	case Day:
		ii.interval = dayPeriod(before)
	case Month:
		ii.interval = monthPeriod(before)
	case Year:
		ii.interval = yearPeriod(before)
	default:
	}

	return interval
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
		reg+int64(86400*c.Len()))

	intervals := []Interval{}
	for ii.HasNext() {
		current := ii.Next()

		intervals = append(intervals, current)
	}

	return intervals
}
