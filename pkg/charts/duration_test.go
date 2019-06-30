package charts

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestChartsPeriod(t *testing.T) {
	cases := []struct {
		cha        Charts
		period     string
		registered rsrc.Day
		col        Column
		ok         bool
	}{
		{
			Charts{
				"a": []float64{3, 4, 5},
				"b": []float64{2, 3, 6},
			},
			"2009", rsrc.ParseDay("2009-12-30"),
			Column{{"a", 7}, {"b", 5}}, true,
		},
		{
			Charts{
				"a": []float64{3, 4, 5},
				"b": []float64{2, 3, 6},
			},
			"2010", rsrc.ParseDay("2009-12-30"),
			Column{{"b", 6}, {"a", 5}}, true,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"42", rsrc.ParseDay("2009-12-30"),
			nil, false,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"-300", rsrc.ParseDay("2009-12-30"),
			nil, false,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"xxxx", rsrc.ParseDay("2009-12-30"),
			nil, false,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"2008", rsrc.ParseDay("2009-12-30"),
			Column{}, true,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"2011", rsrc.ParseDay("2009-12-30"),
			Column{}, true,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"2009-03", rsrc.ParseDay("2009-03-31"),
			Column{{"a", 3}}, true,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"2009-12", rsrc.ParseDay("2009-12-31"),
			Column{{"a", 3}}, true,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"1e+5-12", rsrc.ParseDay("2009-12-31"),
			nil, false,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"2009-xx", rsrc.ParseDay("2009-12-31"),
			nil, false,
		},
		{
			Charts{"a": []float64{3, 4, 5}},
			"1999012", rsrc.ParseDay("2009-12-31"),
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
				col := c.cha.Sum().Interval(period, c.registered)

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
			Charts{},
			"M", true, rsrc.ParseDay("2011-10-11"),
			Charts{},
		},
		{
			Charts{"a": []float64{12, 33, 10}},
			"d", true, rsrc.ParseDay("2011-10-11"),
			Charts{"a": []float64{12, 33, 10}},
		},
		{
			Charts{"a": iotaF(30), "b": repeat(1, 30)},
			"M", true, rsrc.ParseDay("2011-10-11"),
			Charts{"a": []float64{210, 225}, "b": []float64{21, 9}},
		},
		{
			Charts{"a": repeat(2, 400)},
			"y", true, rsrc.ParseDay("2011-12-25"),
			Charts{"a": []float64{14, 732, 54}},
		},
		{
			Charts{"a": repeat(2, 400)},
			"xx", false, rsrc.ParseDay("2011-12-25"),
			nil,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sum := c.charts.Sum()
			intervals, err := sum.ToIntervals(c.descr, c.registered)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}

			if c.ok {
				charts := sum.Intervals(intervals, c.registered)

				if !reflect.DeepEqual(charts, c.intervals) {
					t.Errorf("wrong data:\nhas:  %v\nwant: %v", charts, c.intervals)
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
		t.Run("", func(t *testing.T) {
			days := []float64{}
			for i := 0; i < c.n; i++ {
				days = append(days, 0)
			}
			cha := Charts{"x": days}

			intervals, err := cha.ToIntervals(c.descr, c.registered)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if c.ok {
				if !reflect.DeepEqual(c.intervals, intervals) {
					t.Errorf("%v != %v", c.intervals, intervals)
				}
			}
		})
	}
}
