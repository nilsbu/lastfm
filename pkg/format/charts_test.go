package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

func TestChartsCSV(t *testing.T) {
	cases := []struct {
		name       string
		charts     []charts.Charts
		numbered   bool
		precision  int
		percentage bool
		decimal    string
		str        string
	}{
		{
			"empty charts",
			[]charts.Charts{charts.FromMap(map[string][]float64{})},
			false, 0, false, ".",
			"\"Name\";\"Value\"\n",
		},
		{
			"1",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"ABC": {123.4},
					"X":   {1.238},
				})},
			false, 2, false, ",",
			"\"Name\";\"Value\"\n\"ABC\";123,40\n\"X\";  1,24\n",
		},
		{
			"percentage",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"a": {.75},
					"b": {.25},
				})},
			false, 0, true, ".",
			"\"Name\";\"Value\"\n\"a\";75%\n\"b\";25%\n",
		},
		{
			"comma for decimals",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"a": {12.1},
					"b": {4},
				})},
			true, 1, false, ",",
			"\"#\";\"Name\";\"Value\"\n1;\"a\";12,1\n2;\"b\"; 4,0\n",
		},
		{
			"multiple charts",
			[]charts.Charts{
				charts.FromMap(
					map[string][]float64{
						"a": {12.1},
						"b": {4},
					}),
				charts.FromMap(
					map[string][]float64{
						"X": {5},
						"b": {4},
					}),
			},
			true, 1, false, ".",
			"\"#\";\"Name\";\"Value\";\"#\";\"Name\";\"Value\"\n1;\"a\";12.1;\"X\";5.0\n2;\"b\"; 4.0;\"b\";4.0\n",
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
		charts     []charts.Charts
		numbered   bool
		precision  int
		percentage bool
		str        string
	}{
		{
			"empty charts",
			[]charts.Charts{charts.FromMap(map[string][]float64{})},
			false, 0, false,
			"",
		},
		{
			"simple one-liner",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"ABC": {1, 2, 3},
				})},
			false, 0, false,
			"ABC - 3\n",
		},
		{
			"alignment correct",
			[]charts.Charts{charts.InOrder([]charts.Pair{
				{Title: charts.ArtistTitle("Týrs"), Values: []float64{12}},
				{Title: charts.ArtistTitle("AB"), Values: []float64{3}},
			})},
			false, 0, false,
			"Týrs - 12\nAB   -  3\n",
		},
		{
			"correct precision",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"ABC": {123.4},
					"X":   {1.238},
				})},
			false, 2, false,
			"ABC - 123.40\nX   -   1.24\n",
		},
		{
			"numbered",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"A": {10}, "B": {9}, "C": {8},
					"D": {7}, "E": {6}, "F": {5},
					"G": {4}, "H": {3}, "I": {2},
					"J": {1},
				})},
			true, 0, false,
			" 1: A - 10\n 2: B -  9\n 3: C -  8\n 4: D -  7\n 5: E -  6\n 6: F -  5\n 7: G -  4\n 8: H -  3\n 9: I -  2\n10: J -  1\n",
		},
		{
			"percentage with sum total",
			[]charts.Charts{charts.InOrder([]charts.Pair{
				{Title: charts.ArtistTitle("Týrs"), Values: []float64{.6}},
				{Title: charts.ArtistTitle("AB"), Values: []float64{.4}},
			})},
			false, 1, true,
			"Týrs - 60.0%\nAB   - 40.0%\n",
		},
		{
			"zero percentage",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"AB": {0},
				})},
			true, 0, true,
			"1: AB - 0%\n",
		},
		{
			"2 charts numbered",
			[]charts.Charts{charts.FromMap(
				map[string][]float64{
					"A": {10}, "B": {9}, "C": {8},
					"D": {7}, "E": {6}, "F": {5},
					"G": {4}, "H": {3}, "I": {2},
					"J": {1},
				}),
				charts.FromMap(
					map[string][]float64{
						"A": {10}, "B": {9}, "C": {8},
						"D": {7}, "E": {6}, "F": {5},
						"G": {4}, "H": {3}, "I": {2},
						"J": {1},
					})},
			true, 0, false,
			" 1: A - 10\tA - 10\n 2: B -  9\tB -  9\n 3: C -  8\tC -  8\n 4: D -  7\tD -  7\n 5: E -  6\tE -  6\n 6: F -  5\tF -  5\n 7: G -  4\tG -  4\n 8: H -  3\tH -  3\n 9: I -  2\tI -  2\n10: J -  1\tJ -  1\n",
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
		charts     []charts.Charts
		n          int
		numbered   bool
		precision  int
		percentage bool
		str        string
	}{
		{
			"1 charts",
			[]charts.Charts{charts.InOrder([]charts.Pair{
				{Title: charts.ArtistTitle("Týrs"), Values: []float64{12}},
				{Title: charts.ArtistTitle("AB"), Values: []float64{3}},
			})},
			2, false, 0, false,
			"<table><tr><td>Týrs</td><td>12</td></tr><tr><td>AB</td><td> 3</td></tr></table>", // TODO: numbers are aligned by length by don't neet to be
		},
		{
			"2 charts",
			[]charts.Charts{charts.InOrder([]charts.Pair{
				{Title: charts.ArtistTitle("Týrs"), Values: []float64{12}},
				{Title: charts.ArtistTitle("AB"), Values: []float64{3}},
			}),
				charts.InOrder([]charts.Pair{
					{Title: charts.ArtistTitle("X"), Values: []float64{12}},
					{Title: charts.ArtistTitle("AB"), Values: []float64{3}},
				})},
			2, true, 0, false,
			"<table><tr><td>1</td><td>Týrs</td><td>12</td><td>X</td><td>12</td></tr><tr><td>2</td><td>AB</td><td> 3</td><td>AB</td><td> 3</td></tr></table>",
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
