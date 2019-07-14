package charts

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestIntervals(t *testing.T) {
	cases := []struct {
		descr     string
		in        Intervals
		intervals []Interval
		testDates map[string]int
	}{
		{
			"1d",
			Days(rsrc.ParseDay("2001-09-11"), rsrc.ParseDay("2001-09-14")),
			[]Interval{
				{rsrc.ParseDay("2001-09-11"), rsrc.ParseDay("2001-09-12")},
				{rsrc.ParseDay("2001-09-12"), rsrc.ParseDay("2001-09-13")},
				{rsrc.ParseDay("2001-09-13"), rsrc.ParseDay("2001-09-14")},
			},
			map[string]int{
				"2001-09-12": 1,
			},
		},
		{
			"5d",
			MultiDays(rsrc.ParseDay("2001-09-11"), rsrc.ParseDay("2001-09-21"), 5),
			[]Interval{
				{rsrc.ParseDay("2001-09-11"), rsrc.ParseDay("2001-09-16")},
				{rsrc.ParseDay("2001-09-16"), rsrc.ParseDay("2001-09-21")},
			},
			map[string]int{
				"2001-09-20": 1,
			},
		},
		{
			"1M",
			Months(rsrc.ParseDay("2001-09-07"), rsrc.ParseDay("2002-03-01"), 1),
			[]Interval{
				{rsrc.ParseDay("2001-09-01"), rsrc.ParseDay("2001-10-01")},
				{rsrc.ParseDay("2001-10-01"), rsrc.ParseDay("2001-11-01")},
				{rsrc.ParseDay("2001-11-01"), rsrc.ParseDay("2001-12-01")},
				{rsrc.ParseDay("2001-12-01"), rsrc.ParseDay("2002-01-01")},
				{rsrc.ParseDay("2002-01-01"), rsrc.ParseDay("2002-02-01")},
				{rsrc.ParseDay("2002-02-01"), rsrc.ParseDay("2002-03-01")},
			},
			map[string]int{
				"2001-09-12": 0,
				"2001-12-31": 3,
				"2002-02-01": 5,
			},
		},
		{
			"7M",
			Months(rsrc.ParseDay("2001-09-11"), rsrc.ParseDay("2002-09-25"), 7),
			[]Interval{
				{rsrc.ParseDay("2001-08-01"), rsrc.ParseDay("2002-03-01")},
				{rsrc.ParseDay("2002-03-01"), rsrc.ParseDay("2002-10-01")},
			},
			map[string]int{
				"2002-02-12": 0,
				"2002-03-12": 1,
			},
		},
		{
			"1y",
			Years(rsrc.ParseDay("2001-01-01"), rsrc.ParseDay("2003-01-01"), 1),
			[]Interval{
				{rsrc.ParseDay("2001-01-01"), rsrc.ParseDay("2002-01-01")},
				{rsrc.ParseDay("2002-01-01"), rsrc.ParseDay("2003-01-01")},
			},
			map[string]int{
				"2002-09-12": 1,
			},
		},
		{
			"10y",
			Years(rsrc.ParseDay("2001-01-01"), rsrc.ParseDay("2016-02-29"), 10),
			[]Interval{
				{rsrc.ParseDay("2000-01-01"), rsrc.ParseDay("2010-01-01")},
				{rsrc.ParseDay("2010-01-01"), rsrc.ParseDay("2020-01-01")},
			},
			map[string]int{
				"2010-09-12": 1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.descr, func(t *testing.T) {
			if len(c.intervals) != c.in.Len() {
				t.Fatalf("expected length: '%v', actual: '%v'",
					len(c.intervals), c.in.Len())
			}

			for i, interval := range c.intervals {
				has := c.in.At(i)
				if interval.Begin.Midnight() != has.Begin.Midnight() {
					t.Errorf("expected begin '%v' but has '%v' at index %v",
						interval.Begin, has.Begin, i)
				}
				if interval.Before.Midnight() != has.Before.Midnight() {
					t.Errorf("expected before '%v' but has '%v' at index %v",
						interval.Before, has.Before, i)
				}
			}

			for dateStr, index := range c.testDates {
				hasIndex := c.in.Index(rsrc.ParseDay(dateStr))
				if index != hasIndex {
					t.Errorf("expected index %v for '%v' but got '%v'",
						index, dateStr, hasIndex)
				}
			}
		})
	}
}
