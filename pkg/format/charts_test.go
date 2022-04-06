package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestCharts(t *testing.T) {
	cases := []struct {
		name      string
		formatter *Charts
		decimal   string
		csv       string
		plain     string
		html      string
	}{
		{
			"empty charts",
			&Charts{
				Charts:     []charts.Charts{charts.FromMap(map[string][]float64{})},
				Numbered:   false,
				Precision:  0,
				Percentage: false,
			},
			".",
			"\"Name\";\"Value\"\n",
			"",
			"<table></table>",
		},
		{
			"1",
			&Charts{
				Charts: []charts.Charts{charts.FromMap(
					map[string][]float64{
						"ABC": {123.4},
						"X":   {1.238},
					})},
				Numbered:   false,
				Precision:  2,
				Percentage: false,
			}, ",",
			"\"Name\";\"Value\"\n\"ABC\";123,40\n\"X\";1,24\n",
			"ABC - 123.40\nX   -   1.24\n",
			"<table><tr><td>ABC</td><td>123.40</td></tr><tr><td>X</td><td>1.24</td></tr></table>",
		},
		{
			"percentage",
			&Charts{
				Charts: []charts.Charts{charts.FromMap(
					map[string][]float64{
						"a": {.75},
						"b": {.25},
					})},
				Numbered:   false,
				Precision:  0,
				Percentage: true,
			}, ".",
			"\"Name\";\"Value\"\n\"a\";75%\n\"b\";25%\n",
			"a - 75%\nb - 25%\n",
			"<table><tr><td>a</td><td>75%</td></tr><tr><td>b</td><td>25%</td></tr></table>",
		},
		{
			"comma for decimals",
			&Charts{
				Charts: []charts.Charts{charts.FromMap(
					map[string][]float64{
						"a": {12.1},
						"b": {4},
					})},
				Numbered:   true,
				Precision:  1,
				Percentage: false,
			}, ",",
			"\"#\";\"Name\";\"Value\"\n1;\"a\";12,1\n2;\"b\";4,0\n",
			"1: a - 12.1\n2: b -  4.0\n",
			"<table><tr><td>1</td><td>a</td><td>12.1</td></tr><tr><td>2</td><td>b</td><td>4.0</td></tr></table>",
		},
		{
			"numbered",
			&Charts{
				Charts: []charts.Charts{charts.FromMap(
					map[string][]float64{
						"A": {10}, "B": {9}, "C": {8},
						"D": {7}, "E": {6}, "F": {5},
						"G": {4}, "H": {3}, "I": {2},
						"J": {1},
					})},
				Numbered:   true,
				Precision:  0,
				Percentage: false,
			}, ".",
			"\"#\";\"Name\";\"Value\"\n1;\"A\";10\n2;\"B\";9\n3;\"C\";8\n4;\"D\";7\n5;\"E\";6\n6;\"F\";5\n7;\"G\";4\n8;\"H\";3\n9;\"I\";2\n10;\"J\";1\n",
			" 1: A - 10\n 2: B -  9\n 3: C -  8\n 4: D -  7\n 5: E -  6\n 6: F -  5\n 7: G -  4\n 8: H -  3\n 9: I -  2\n10: J -  1\n",
			"<table><tr><td>1</td><td>A</td><td>10</td></tr><tr><td>2</td><td>B</td><td>9</td></tr><tr><td>3</td><td>C</td><td>8</td></tr><tr><td>4</td><td>D</td><td>7</td></tr><tr><td>5</td><td>E</td><td>6</td></tr><tr><td>6</td><td>F</td><td>5</td></tr><tr><td>7</td><td>G</td><td>4</td></tr><tr><td>8</td><td>H</td><td>3</td></tr><tr><td>9</td><td>I</td><td>2</td></tr><tr><td>10</td><td>J</td><td>1</td></tr></table>",
		},
		{
			"multiple charts",
			&Charts{
				Charts: []charts.Charts{
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
				Numbered:   true,
				Precision:  1,
				Percentage: false,
			}, ".",
			"\"#\";\"Name\";\"Value\";\"Name\";\"Value\"\n1;\"a\";12.1;\"X\";5.0\n2;\"b\";4.0;\"b\";4.0\n",
			"1: a - 12.1\tX - 5.0\n2: b -  4.0\tb - 4.0\n",
			"<table><tr><td>1</td><td>a</td><td>12.1</td><td>X</td><td>5.0</td></tr><tr><td>2</td><td>b</td><td>4.0</td><td>b</td><td>4.0</td></tr></table>",
		},
		{
			"multiple charts with range",
			&Charts{
				Charts: []charts.Charts{
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
				Ranges:     charts.ParseRangesTrusted("1y", rsrc.ParseDay("2012-01-01"), 366),
				Numbered:   true,
				Precision:  1,
				Percentage: false,
			}, ".",
			"\"#\";\"Name\";\"Value\";\"Name\";\"Value\"\n1;\"a\";12.1;\"X\";5.0\n2;\"b\";4.0;\"b\";4.0\n",
			"#  2012-01-01       \t2013-01-01      \n1: a          - 12.1\tX          - 5.0\n2: b          -  4.0\tb          - 4.0\n",
			"<table><tr><td>#</td><td>2012-01-01</td><td></td><td>2013-01-01</td><td></td></tr><tr><td>1</td><td>a</td><td>12.1</td><td>X</td><td>5.0</td></tr><tr><td>2</td><td>b</td><td>4.0</td><td>b</td><td>4.0</td></tr></table>",
		},
		{
			"percentage with sum total",
			&Charts{
				Charts: []charts.Charts{charts.InOrder([]charts.Pair{
					{Title: charts.ArtistTitle("Týrs"), Values: []float64{.6}},
					{Title: charts.ArtistTitle("AB"), Values: []float64{.4}},
				})},
				Numbered:   false,
				Precision:  1,
				Percentage: true}, ",",
			"\"Name\";\"Value\"\n\"Týrs\";60,0%\n\"AB\";40,0%\n",
			"Týrs - 60.0%\nAB   - 40.0%\n",
			"<table><tr><td>Týrs</td><td>60.0%</td></tr><tr><td>AB</td><td>40.0%</td></tr></table>",
		},
		{
			"zero percentage",
			&Charts{
				Charts: []charts.Charts{charts.FromMap(
					map[string][]float64{
						"AB": {0},
					})},
				Numbered:   false,
				Precision:  0,
				Percentage: true,
			}, ".",
			"\"Name\";\"Value\"\n\"AB\";0%\n",
			"AB - 0%\n",
			"<table><tr><td>AB</td><td>0%</td></tr></table>",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			{
				buf := new(bytes.Buffer)
				c.formatter.CSV(buf, c.decimal)

				str := buf.String()
				if str != c.csv {
					t.Errorf("false CSV formatting:\nhas:\n%v\nwant:\n%v", str, c.csv)
				}
			}
			{
				buf := new(bytes.Buffer)
				c.formatter.Plain(buf)

				str := buf.String()
				if str != c.plain {
					t.Errorf("false Plain formatting:\nhas:\n%v\nwant:\n%v", str, c.plain)
				}
			}
			{
				buf := new(bytes.Buffer)
				c.formatter.HTML(buf)

				str := buf.String()
				if str != c.html {
					t.Errorf("false HTML formatting:\nhas:\n%v\nwant:\n%v", str, c.html)
				}
			}
		})
	}
}
