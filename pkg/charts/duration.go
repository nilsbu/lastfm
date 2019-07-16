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
		day := rsrc.DayFromTime(begin)
		return Years(day, day.AddDate(0, 0, 1), 1).At(0), nil
	case 7:
		begin, err := time.Parse("2006-01", descr)
		if err != nil {
			return Interval{}, err
		}
		day := rsrc.DayFromTime(begin)
		return Months(day, day.AddDate(0, 0, 1), 1).At(0), nil
	default:
		return Interval{}, fmt.Errorf("interval format '%v' not supported", descr)
	}
}

// Interval returns a Column that sums an interval of the charts. The charts
// have to be a sum.
func (c Charts) Interval(i Interval) Column {
	size := c.Len()

	// A day is subtacted here because the entry at begin & and already the plays
	// from that day.
	from := c.Headers.Index(i.Begin) - 1
	to := c.Headers.Index(i.Before) - 1

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

func (c Charts) Intervals(intervals Intervals) Charts {
	icharts := []map[string]float64{}
	for i := 0; i < intervals.Len(); i++ {
		col := c.Interval(intervals.At(i))

		if len(col) == 0 {
			continue
		}

		cha := map[string]float64{}
		for _, x := range col {
			cha[x.Name] = x.Score
		}

		icharts = append(icharts, cha)
	}

	ncha := CompileArtists(icharts, c.Headers.At(0).Begin)
	ncha.Headers = intervals
	return ncha
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
