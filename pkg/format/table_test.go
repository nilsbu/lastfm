package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestTableCSV(t *testing.T) {
	cases := []struct {
		charts  charts.Charts
		date    rsrc.Day
		step    int
		count   int
		decimal string
		ok      bool
		str     string
	}{
		{
			charts.Charts{},
			rsrc.ParseDay("2012-01-01"),
			1, 2, ",", true,
			"",
		},
		{
			charts.Charts{"X": []float64{}},
			rsrc.ParseDay("2012-01-01"),
			1, 2, ",", true,
			"",
		},
		{
			charts.Charts{
				"ABC": []float64{1.25, 2},
				"X":   []float64{2, 3},
			},
			rsrc.ParseDay("2012-01-01"),
			1, 2, ",", true,
			"\"name\";2012-01-01;2012-01-02\n\"X\";2;3\n\"ABC\";1,25;2\n",
		},
		{
			charts.Charts{
				"A": []float64{1, 2, 3, 4, 5, 6, 7},
			},
			rsrc.ParseDay("2012-01-01"),
			3, 1, ".", true,
			"\"name\";2012-01-01;2012-01-04;2012-01-07\n\"A\";1;4;7\n",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			f := &Table{
				Charts: c.charts,
				First:  c.date,
				Step:   c.step,
				Count:  c.count,
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
		date   rsrc.Day
		step   int
		count  int
		ok     bool
		str    string
	}{
		{
			charts.Charts{},
			rsrc.ParseDay("2012-01-01"),
			1, 2, true,
			"",
		},
		{
			charts.Charts{
				"A": []float64{1.33, 2, 3, 4, 5, 6, 7},
			},
			rsrc.ParseDay("2012-01-01"),
			3, 1, true,
			"A: 1.33, 4, 7\n",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			f := &Table{
				Charts: c.charts,
				First:  c.date,
				Step:   c.step,
				Count:  c.count,
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
