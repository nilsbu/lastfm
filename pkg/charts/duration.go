package charts

import (
	"fmt"
	"math"
	"sort"
	"strconv"
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
	y, err := parse(descr)
	if err != nil {
		return nil, err
	}

	begin, _ := time.Parse("2006-01-02", fmt.Sprintf("%04v-01-01", y))
	before, _ := time.Parse("2006-01-02", fmt.Sprintf("%04v-01-01", y+1))

	return &period{begin, before}, nil
}

func month(descr string) (Interval, error) {
	y, err := parse(descr[:4])
	if err != nil {
		return nil, err
	}
	m, err := parse(descr[5:7])
	if err != nil {
		return nil, err
	}
	if descr[4] != '-' {
		return nil, fmt.Errorf("'%v' does not obey 'YYYY-MM' format", descr)
	}

	begin, _ := time.Parse("2006-01-02", fmt.Sprintf("%04v-%02v-01", y, m))

	if m == 12 {
		y++
		m = 1
	} else {
		m++
	}
	before, _ := time.Parse("2006-01-02", fmt.Sprintf("%04v-%02v-01", y, m))

	return &period{begin, before}, nil
}

func parse(str string) (int, error) {
	x, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	} else if x < 0 {
		// upper bound is ensured by # of digits
		return 0, fmt.Errorf("'%v' is invalid, must be in range [0-%v]",
			x, int(math.Pow(float64(x), float64(len(str)))))
	}
	return x, nil
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
