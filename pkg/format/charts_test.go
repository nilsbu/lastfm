package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

func TestChartsCSV(t *testing.T) {
	cases := []struct {
		name       string
		charts     charts.Charts
		numbered   bool
		precision  int
		percentage bool
		decimal    string
		str        string
	}{
		{
			"empty charts",
			charts.FromMap(map[string][]float64{}),
			false, 0, false, ".",
			"\"Name\";\"Value\"\n",
		},
		{
			"1",
			charts.FromMap(
				map[string][]float64{
					"ABC": {123.4},
					"X":   {1.238},
				}),
			false, 2, false, ",",
			"\"Name\";\"Value\"\n\"ABC\";123,40\n\"X\";  1,24\n",
		},
		{
			"percentage",
			charts.FromMap(
				map[string][]float64{
					"a": {.75},
					"b": {.25},
				}),
			false, 0, true, ".",
			"\"Name\";\"Value\"\n\"a\";75%\n\"b\";25%\n",
		},
		{
			"comma for decimals",
			charts.FromMap(
				map[string][]float64{
					"a": {12.1},
					"b": {4},
				}),
			true, 1, false, ",",
			"\"#\";\"Name\";\"Value\"\n1;\"a\";12,1\n2;\"b\"; 4,0\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Charts{
				Charts:     c.charts,
				Numbered:   c.numbered,
				Precision:  c.precision,
				Percentage: c.percentage,
			}
			formatter.CSV(buf, c.decimal)

			str := buf.String()
			if str != c.str {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
			}
		})
	}
}

func TestChartsPlain(t *testing.T) {
	cases := []struct {
		name       string
		charts     charts.Charts
		numbered   bool
		precision  int
		percentage bool
		str        string
	}{
		{
			"empty charts",
			charts.FromMap(map[string][]float64{}),
			false, 0, false,
			"",
		},
		{
			"simple one-liner",
			charts.FromMap(
				map[string][]float64{
					"ABC": {1, 2, 3},
				}),
			false, 0, false,
			"ABC - 3\n",
		},
		{
			"alignment correct",
			charts.InOrder([]charts.Pair{
				{Title: charts.ArtistTitle("Týrs"), Values: []float64{12}},
				{Title: charts.ArtistTitle("AB"), Values: []float64{3}},
			}),
			false, 0, false,
			"Týrs - 12\nAB   -  3\n",
		},
		{
			"correct precision",
			charts.FromMap(
				map[string][]float64{
					"ABC": {123.4},
					"X":   {1.238},
				}),
			false, 2, false,
			"ABC - 123.40\nX   -   1.24\n",
		},
		{
			"numbered",
			charts.FromMap(
				map[string][]float64{
					"A": {10}, "B": {9}, "C": {8},
					"D": {7}, "E": {6}, "F": {5},
					"G": {4}, "H": {3}, "I": {2},
					"J": {1},
				}),
			true, 0, false,
			" 1: A - 10\n 2: B -  9\n 3: C -  8\n 4: D -  7\n 5: E -  6\n 6: F -  5\n 7: G -  4\n 8: H -  3\n 9: I -  2\n10: J -  1\n",
		},
		{
			"percentage with sum total",
			charts.InOrder([]charts.Pair{
				{Title: charts.ArtistTitle("Týrs"), Values: []float64{.6}},
				{Title: charts.ArtistTitle("AB"), Values: []float64{.4}},
			}),
			false, 1, true,
			"Týrs - 60.0%\nAB   - 40.0%\n",
		},
		{
			"zero percentage",
			charts.FromMap(
				map[string][]float64{
					"AB": {0},
				}),
			true, 0, true,
			"1: AB - 0%\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Charts{
				Charts:     c.charts,
				Numbered:   c.numbered,
				Precision:  c.precision,
				Percentage: c.percentage,
			}
			formatter.Plain(buf)

			str := buf.String()
			if str != c.str {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
			}
		})
	}
}

func TestChartsHTML(t *testing.T) {
	cases := []struct {
		name       string
		charts     charts.Charts
		n          int
		numbered   bool
		precision  int
		percentage bool
		str        string
	}{
		{
			"alignment correct",
			charts.InOrder([]charts.Pair{
				{Title: charts.ArtistTitle("Týrs"), Values: []float64{12}},
				{Title: charts.ArtistTitle("AB"), Values: []float64{3}},
			}),
			2, false, 0, false,
			"<table><tr><td>Týrs</td><td>12</td></tr><tr><td>AB</td><td> 3</td></tr></table>", // TODO: numbers are aligned by length by don't neet to be
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Charts{
				Charts:     c.charts,
				Numbered:   c.numbered,
				Precision:  c.precision,
				Percentage: c.percentage,
			}
			formatter.HTML(buf)

			str := buf.String()
			if str != c.str {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
			}
		})
	}
}
