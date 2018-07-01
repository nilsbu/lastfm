package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

func TestChartsPlain(t *testing.T) {
	cases := []struct {
		name     string
		charts   charts.Charts
		col      int
		n        int
		numbered bool
		prec     int
		str      string
	}{
		{
			"empty charts",
			charts.Charts{},
			-1, 3, false, 0,
			"",
		},
		{
			"simple one-liner",
			charts.Charts{
				"ABC": []float64{1, 2, 3},
			},
			-1, 3, false, 0,
			"ABC - 3\n",
		},
		{
			"alignment correct",
			charts.Charts{
				"AKSLJDHLJKH": []float64{1},
				"AB":          []float64{3},
				"Týrs":        []float64{12},
			},
			0, 2, false, 0,
			"Týrs - 12\nAB   -  3\n",
		},
		{
			"correct precision",
			charts.Charts{
				"ABC": []float64{123.4},
				"X":   []float64{1.238},
			},
			-1, 2, false, 2,
			"ABC - 123.40\nX   -   1.24\n",
		},
		{
			"numbered",
			charts.Charts{
				"A": []float64{10}, "B": []float64{9}, "C": []float64{8},
				"D": []float64{7}, "E": []float64{6}, "F": []float64{5},
				"G": []float64{4}, "H": []float64{3}, "I": []float64{2},
				"J": []float64{1},
			},
			-1, 10, true, 0,
			" 1: A - 10\n 2: B -  9\n 3: C -  8\n 4: D -  7\n 5: E -  6\n 6: F -  5\n 7: G -  4\n 8: H -  3\n 9: I -  2\n10: J -  1\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Charts{c.charts, c.col, c.n, c.numbered, c.prec}
			formatter.Plain(buf)

			str := buf.String()
			if str != c.str {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
			}
		})
	}
}
