package charts

import (
	"fmt"
	"sort"
	"time"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Interval interface {
	Begin() int64
	Before() int64
}

type period struct {
	begin  time.Time
	before time.Time
}

func (p period) Begin() int64 {
	return p.begin.Unix()
}

func (p period) Before() int64 {
	return p.before.Unix()
}

func Period(descr string) (Interval, error) {
	switch len(descr) {
	case 4:
		return year(descr)
	case 7:
		return month(descr)
	default:
		return nil, fmt.Errorf("interval format '%v' not supported", descr)
	}
}

func year(descr string) (Interval, error) {
	begin, err := time.Parse("2006", descr)
	if err != nil {
		return nil, err
	}

	return yearPeriod(begin), nil
}

func month(descr string) (Interval, error) {
	begin, err := time.Parse("2006-01", descr)
	if err != nil {
		return nil, err
	}
	return monthPeriod(begin), nil
}

func yearPeriod(t time.Time) Interval {
	y := t.Year()
	return &period{
		begin:  time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC),
		before: time.Date(y+1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

func monthPeriod(t time.Time) Interval {
	y, m := t.Year(), t.Month()
	return &period{
		begin:  time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC),
		before: time.Date(y, time.Month(m+1), 1, 0, 0, 0, 0, time.UTC),
	}
}

func (c Charts) Interval(i Interval, registered rsrc.Day) Column {
	var size int
	for _, line := range c {
		size = len(line)
		break
	}

	offset, _ := registered.Midnight()
	offset /= 86400
	from := int(i.Begin()/86400 - offset - 1)

	to := int(i.Before()/86400 - offset - 1)
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
