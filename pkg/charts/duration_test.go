package charts

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestChartsPeriod(t *testing.T) {
	cases := []struct {
		cha    Charts
		period string
		col    Column
		ok     bool
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-30"), rsrc.ParseDay("2010-01-02")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{3, 4, 5}, {2, 3, 6}}},
			"2009",
			Column{{"a", 7}, {"b", 5}}, true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-30"), rsrc.ParseDay("2010-01-02")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{3, 4, 5}, {2, 3, 6}}},
			"2010",
			Column{{"b", 6}, {"a", 5}}, true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-30"), rsrc.ParseDay("2010-01-02")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"42",
			nil, false,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-30"), rsrc.ParseDay("2010-01-02")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"-300",
			nil, false,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-30"), rsrc.ParseDay("2010-01-02")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"xxxx",
			nil, false,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-30"), rsrc.ParseDay("2010-01-02")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"2008",
			Column{}, true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-30"), rsrc.ParseDay("2000-01-04")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"2011",
			Column{}, true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-03-31"), rsrc.ParseDay("2010-04-02")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"2009-03",
			Column{{"a", 3}}, true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-31"), rsrc.ParseDay("2010-01-03")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"2009-12",
			Column{{"a", 3}}, true,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-31"), rsrc.ParseDay("2010-01-03")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"1e+5-12",
			nil, false,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-31"), rsrc.ParseDay("2010-01-03")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"2009-xx",
			nil, false,
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2009-12-31"), rsrc.ParseDay("2010-01-03")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{3, 4, 5}}},
			"1999012",
			nil, false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			period, err := Period(c.period)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				col := c.cha.Sum().Interval(period)

				if !reflect.DeepEqual(col, c.col) {
					t.Errorf("wrong data:\nhas:  %v\nwant: %v", col, c.col)
				}
			}
		})
	}
}

func iotaF(n int) []float64 {
	nums := make([]float64, n)
	for i := range nums {
		nums[i] = float64(i)
	}

	return nums
}

func repeat(x, n int) []float64 {
	nums := make([]float64, n)
	for i := range nums {
		nums[i] = float64(x)
	}

	return nums
}

func TestChartsSumIntervals(t *testing.T) {
	cases := []struct {
		charts     Charts
		descr      string
		ok         bool
		registered rsrc.Day
		intervals  Charts
	}{
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2011-10-11"), rsrc.ParseDay("2011-10-14")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			"M", true, rsrc.ParseDay("2011-10-11"),
			Charts{
				Headers: Months(rsrc.ParseDay("2011-10-01"), rsrc.ParseDay("2011-11-01"), 1),
				Keys:    []Key{},
				Values:  [][]float64{}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2011-10-11"), rsrc.ParseDay("2011-10-14")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{12, 33, 10}}},
			"d", true, rsrc.ParseDay("2011-10-11"),
			Charts{
				Headers: Days(rsrc.ParseDay("2011-10-11"), rsrc.ParseDay("2011-10-14")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{12, 33, 10}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2011-10-11"), rsrc.ParseDay("2011-11-10")),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{iotaF(30), repeat(1, 30)}},
			"M", true, rsrc.ParseDay("2011-10-11"),
			Charts{
				Headers: Months(rsrc.ParseDay("2011-10-01"), rsrc.ParseDay("2011-12-01"), 1),
				Keys:    []Key{simpleKey("a"), simpleKey("b")},
				Values:  [][]float64{{210, 225}, {21, 9}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2011-12-25"), rsrc.ParseDay("2013-01-28")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{repeat(2, 400)}},
			"y", true, rsrc.ParseDay("2011-12-25"),
			Charts{
				Headers: Years(rsrc.ParseDay("2011-01-01"), rsrc.ParseDay("2014-01-01"), 1),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{{14, 732, 54}}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2011-12-25"), rsrc.ParseDay("2013-01-28")),
				Keys:    []Key{simpleKey("a")},
				Values:  [][]float64{repeat(2, 400)}},
			"xx", false, rsrc.ParseDay("2011-12-25"),
			Charts{},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sum := c.charts.Sum()

			intervals, err := ToIntervals(
				c.descr,
				sum.Headers.At(0).Begin,
				sum.Headers.At(sum.Headers.Len()).Begin)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}

			if c.ok {
				charts := sum.Intervals(intervals)

				if err = c.intervals.AssertEqual(charts); err != nil {
					fmt.Println(c.charts.Headers.At(0), c.charts.Headers.Len())
					fmt.Println(c.intervals.Headers.At(0), c.intervals.Headers.Len())
					fmt.Println(intervals.At(0), intervals.Len())
					fmt.Println(charts.Headers.At(0), charts.Headers.Len())
					t.Errorf("charts are wrong: %v", err)
				}
			}
		})
	}
}

func TestChartsToIntervals(t *testing.T) {
	cases := []struct {
		descr      string
		registered rsrc.Day
		n          int
		intervals  []Interval
		ok         bool
	}{
		{
			"y", rsrc.ParseDay("2007-01-01"), 600,
			[]Interval{
				{Begin: rsrc.ParseDay("2007-01-01"), Before: rsrc.ParseDay("2008-01-01")},
				{Begin: rsrc.ParseDay("2008-01-01"), Before: rsrc.ParseDay("2009-01-01")},
			}, true,
		},
		{
			"y", rsrc.ParseDay("2007-02-01"), 3,
			[]Interval{
				{Begin: rsrc.ParseDay("2007-01-01"), Before: rsrc.ParseDay("2008-01-01")},
			}, true,
		},
		{
			"M", rsrc.ParseDay("2007-02-01"), 30,
			[]Interval{
				{Begin: rsrc.ParseDay("2007-02-01"), Before: rsrc.ParseDay("2007-03-01")},
				{Begin: rsrc.ParseDay("2007-03-01"), Before: rsrc.ParseDay("2007-04-01")},
			}, true,
		},
		{
			"asdasd", rsrc.ParseDay("2007-02-01"), 30,
			[]Interval{},
			false,
		},
		{
			"2M", rsrc.ParseDay("2007-01-01"), 110,
			[]Interval{
				{Begin: rsrc.ParseDay("2007-01-01"), Before: rsrc.ParseDay("2007-03-01")},
				{Begin: rsrc.ParseDay("2007-03-01"), Before: rsrc.ParseDay("2007-05-01")},
			}, true,
		},
		{
			"6M", rsrc.ParseDay("2007-03-01"), 365,
			[]Interval{
				{Begin: rsrc.ParseDay("2007-01-01"), Before: rsrc.ParseDay("2007-07-01")},
				{Begin: rsrc.ParseDay("2007-07-01"), Before: rsrc.ParseDay("2008-01-01")},
				{Begin: rsrc.ParseDay("2008-01-01"), Before: rsrc.ParseDay("2008-07-01")},
			}, true,
		},
		{
			"10y", rsrc.ParseDay("2007-03-01"), 3653,
			[]Interval{
				{Begin: rsrc.ParseDay("2000-01-01"), Before: rsrc.ParseDay("2010-01-01")},
				{Begin: rsrc.ParseDay("2010-01-01"), Before: rsrc.ParseDay("2020-01-01")},
			}, true,
		},
		{
			"3y", rsrc.ParseDay("2008-03-01"), 2,
			[]Interval{
				{Begin: rsrc.ParseDay("2007-01-01"), Before: rsrc.ParseDay("2010-01-01")},
			}, true,
		},
	}

	for _, c := range cases {
		t.Run(c.descr, func(t *testing.T) {
			days := []float64{}
			for i := 0; i < c.n; i++ {
				days = append(days, 0)
			}
			cha := Charts{
				Headers: Days(c.registered, c.registered.AddDate(0, 0, c.n)),
				Keys:    []Key{simpleKey("x")},
				Values:  [][]float64{days}}

			intervals, err := ToIntervals(
				c.descr,
				c.registered,
				c.registered.AddDate(0, 0, cha.Len()))
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if c.ok {
				if len(c.intervals) != intervals.Len() {
					t.Fatalf("expected length: '%v', actual: '%v'",
						len(c.intervals), intervals.Len())
				}

				for i, interval := range c.intervals {
					has := intervals.At(i)
					if interval.Begin.Midnight() != has.Begin.Midnight() {
						t.Errorf("expected begin '%v' but has '%v' at index %v",
							interval.Begin, has.Begin, i)
					}
					if interval.Before.Midnight() != has.Before.Midnight() {
						t.Errorf("expected before '%v' but has '%v' at index %v",
							interval.Before, has.Before, i)
					}
				}
			}
		})
	}
}
