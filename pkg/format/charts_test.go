package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

func TestChartsCSV(t *testing.T) {
	cases := []struct {
		charts     charts.Charts
		col        int
		n          int
		numbered   bool
		precision  int
		percentage bool
		decimal    string
		str        string
	}{
		{
			charts.Charts{},
			-1, 3, false, 0, false, ".",
			"",
		},
		{
			charts.Charts{
				"ABC": []float64{123.4},
				"X":   []float64{1.238},
			},
			-1, 2, false, 2, false, ",",
			"\"Name\";\"Value\"\n\"ABC\";123,40\n\"X\";  1,24\n",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Charts{
				Charts:     c.charts,
				Column:     c.col,
				Count:      c.n,
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
		col        int
		n          int
		numbered   bool
		precision  int
		percentage bool
		str        string
	}{
		{
			"empty charts",
			charts.Charts{},
			-1, 3, false, 0, false,
			"",
		},
		{
			"simple one-liner",
			charts.Charts{
				"ABC": []float64{1, 2, 3},
			},
			-1, 3, false, 0, false,
			"ABC - 3\n",
		},
		{
			"alignment correct",
			charts.Charts{
				"AKSLJDHLJKH": []float64{1},
				"AB":          []float64{3},
				"Týrs":        []float64{12},
			},
			0, 2, false, 0, false,
			"Týrs - 12\nAB   -  3\n",
		},
		{
			"correct precision",
			charts.Charts{
				"ABC": []float64{123.4},
				"X":   []float64{1.238},
			},
			-1, 2, false, 2, false,
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
			-1, 0, true, 0, false,
			" 1: A - 10\n 2: B -  9\n 3: C -  8\n 4: D -  7\n 5: E -  6\n 6: F -  5\n 7: G -  4\n 8: H -  3\n 9: I -  2\n10: J -  1\n",
		},
		{
			"percentage with sum total",
			charts.Charts{
				"AKSLJDHLJKH": []float64{1},
				"AB":          []float64{3},
				"Týrs":        []float64{4},
			},
			0, 2, false, 1, true,
			"Týrs - 50.0%\nAB   - 37.5%\n",
		},
		{
			"zero percentage",
			charts.Charts{
				"AB": []float64{0},
			},
			0, 2, true, 0, true,
			"1: AB - 0%\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Charts{
				Charts:     c.charts,
				Column:     c.col,
				Count:      c.n,
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

func TestColumnCSV(t *testing.T) {
	cases := []struct {
		name       string
		col        charts.Column
		numbered   bool
		precision  int
		percentage bool
		sumTotal   float64
		decimal    string
		str        string
	}{
		{
			"empty column",
			charts.Column{},
			false, 3, false, 0, ".",
			"",
		},
		{
			"percentage with no total",
			charts.Column{{Name: "a", Score: 12}, {Name: "b", Score: 4}},
			false, 0, true, 0, ".",
			"\"Name\";\"Value\"\n\"a\";75%\n\"b\";25%\n",
		},
		{
			"percentage with no total",
			charts.Column{{Name: "a", Score: 12.1}, {Name: "b", Score: 4}},
			true, 1, false, 0, ",",
			"\"#\";\"Name\";\"Value\"\n1;\"a\";12,1\n2;\"b\"; 4,0\n",
		},
		// rest is covered by TestChartsCSV
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Column{
				Column:     c.col,
				Numbered:   c.numbered,
				Precision:  c.precision,
				Percentage: c.percentage,
				SumTotal:   c.sumTotal,
			}
			formatter.CSV(buf, c.decimal)

			str := buf.String()
			if str != c.str {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
			}
		})
	}
}

func TestColumnPlain(t *testing.T) {
	cases := []struct {
		name       string
		col        charts.Column
		numbered   bool
		precision  int
		percentage bool
		sumTotal   float64
		str        string
	}{
		{
			"empty column",
			charts.Column{},
			false, 3, false, 0,
			"",
		},
		{
			"percentage with no total",
			charts.Column{{Name: "a", Score: 12}, {Name: "b", Score: 4}},
			false, 0, true, 0,
			"a - 75%\nb - 25%\n",
		},
		// rest is covered by TestChartsPlain
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Column{
				Column:     c.col,
				Numbered:   c.numbered,
				Precision:  c.precision,
				Percentage: c.percentage,
				SumTotal:   c.sumTotal,
			}
			formatter.Plain(buf)

			str := buf.String()
			if str != c.str {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", str, c.str)
			}
		})
	}
}
