package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func trustRanges(s string, registered rsrc.Day, l int) charts.Ranges {
	ranges, _ := charts.ParseRanges(s, registered, l)
	return ranges
}

func TestTableCSV(t *testing.T) {
	cases := []struct {
		charts  charts.Charts
		ranges  charts.Ranges
		decimal string
		ok      bool
		str     string
	}{
		{
			charts.FromMap(map[string][]float64{}),
			trustRanges("1d", rsrc.ParseDay("2012-01-01"), 1),
			",", true,
			"\"name\";\n",
		},
		{
			charts.InOrder([]charts.Pair{
				{Title: charts.StringTitle("X"), Values: []float64{2, 3}},
				{Title: charts.StringTitle("ABC"), Values: []float64{1.25, 2}},
			}),
			trustRanges("1d", rsrc.ParseDay("2012-01-01"), 1),
			",", true,
			"\"name\";2012-01-01;2012-01-02\n\"X\";2;3\n\"ABC\";1,25;2\n",
		},
		{
			charts.FromMap(map[string][]float64{
				"A": {1, 4, 7},
			}),
			trustRanges("3d", rsrc.ParseDay("2012-01-01"), 8),
			".", true,
			"\"name\";2012-01-01;2012-01-04;2012-01-07\n\"A\";1;4;7\n",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			f := &Table{
				Charts: c.charts,
				Ranges: c.ranges,
			}
			err := f.CSV(buf, c.decimal)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if err == nil {
				str := buf.String()
				if str != c.str {
					t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
				}
			}
		})
	}
}

func TestTablePlain(t *testing.T) {
	cases := []struct {
		charts charts.Charts
		ranges charts.Ranges
		ok     bool
		str    string
	}{
		{
			charts.FromMap(map[string][]float64{}),
			trustRanges("1d", rsrc.ParseDay("2012-01-01"), 1),
			true,
			"",
		},
		{
			charts.FromMap(map[string][]float64{
				"A": {1.33, 4, 7},
			}),
			trustRanges("3d", rsrc.ParseDay("2012-01-01"), 8),
			true,
			"A: 1.33, 4, 7\n",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			f := &Table{
				Charts: c.charts,
				Ranges: c.ranges,
			}
			err := f.Plain(buf)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if err == nil {
				str := buf.String()
				if str != c.str {
					t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
				}
			}
		})
	}
}

func TestTableHTML(t *testing.T) {
	cases := []struct {
		charts charts.Charts
		ranges charts.Ranges
		ok     bool
		str    string
	}{
		{
			charts.FromMap(map[string][]float64{}),
			trustRanges("1d", rsrc.ParseDay("2012-01-01"), 1),
			true,
			"<table></table>",
		},
		{ // TODO tables HTML needs headers with dates
			charts.FromMap(map[string][]float64{
				"A": {1.33, 4, 7},
			}),
			trustRanges("3d", rsrc.ParseDay("2012-01-01"), 8),
			true,
			"<table><tr><td>A</td><td>1.33</td><td>4</td><td>7</td></tr></table>",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			f := &Table{
				Charts: c.charts,
				Ranges: c.ranges,
			}
			err := f.HTML(buf)
			if err != nil && c.ok {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && !c.ok {
				t.Fatalf("expected error but none occurred")
			}

			if err == nil {
				str := buf.String()
				if str != c.str {
					t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
				}
			}
		})
	}
}
