package format

import (
	"bytes"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

func TestDiffCharts(t *testing.T) {
	cases := []struct {
		name      string
		formatter *DiffCharts
		decimal   string
		csv       string
		plain     string
		html      string
		json      string
	}{
		{
			"empty charts",
			&DiffCharts{
				Charts:     []charts.DiffCharts{charts.NewDiffCharts(charts.FromMap(map[string][]float64{}), 0)},
				Numbered:   false,
				Precision:  0,
				Percentage: false,
			},
			".",
			"\"Name\";\"Value\"\n",
			"",
			"<table></table>",
			"{\"chart\":{\"data\":null}}",
		},
		{
			"1",
			&DiffCharts{
				Charts: []charts.DiffCharts{charts.NewDiffCharts(charts.FromMap(
					map[string][]float64{
						"ABC": {123.4},
						"X":   {1.238},
					}), 0)},
				Numbered:   false,
				Precision:  2,
				Percentage: false,
			}, ",",
			"\"Name\";\"Value\"\n\"ABC\";123,40\n\"X\";1,24\n",
			"ABC - 123.40\nX   -   1.24\n",
			"<table><tr><td>ABC</td><td>123.40</td></tr><tr><td>X</td><td>1.24</td></tr></table>",
			"{\"chart\":{\"data\":[{\"title\":\"ABC\",\"value\":123.4,\"prevPos\":0,\"prevValue\":123.4},{\"title\":\"X\",\"value\":1.238,\"prevPos\":1,\"prevValue\":1.238}]}}",
		},
		{
			"percentage",
			&DiffCharts{
				Charts: []charts.DiffCharts{charts.NewDiffCharts(charts.FromMap(
					map[string][]float64{
						"a": {.75},
						"b": {.25},
					}), 0)},
				Numbered:   false,
				Precision:  0,
				Percentage: true,
			}, ".",
			"\"Name\";\"Value\"\n\"a\";75%\n\"b\";25%\n",
			"a - 75%\nb - 25%\n",
			"<table><tr><td>a</td><td>75%</td></tr><tr><td>b</td><td>25%</td></tr></table>",
			"{\"chart\":{\"data\":[{\"title\":\"a\",\"value\":0.75,\"prevPos\":0,\"prevValue\":0.75},{\"title\":\"b\",\"value\":0.25,\"prevPos\":1,\"prevValue\":0.25}]}}",
		},
		{
			"comma for decimals",
			&DiffCharts{
				Charts: []charts.DiffCharts{charts.NewDiffCharts(charts.FromMap(
					map[string][]float64{
						"a": {2, 6},
						"b": {4, 5},
					}), 0)},
				Numbered:   true,
				Precision:  1,
				Percentage: false,
			}, ",",
			"\"#\";\"Name\";\"Value\"\n1;\"a\";6,0\n2;\"b\";5,0\n",
			"1: a - 6.0\n2: b - 5.0\n",
			"<table><tr><td>1</td><td>a</td><td>6.0</td></tr><tr><td>2</td><td>b</td><td>5.0</td></tr></table>",
			"{\"chart\":{\"data\":[{\"title\":\"a\",\"value\":6,\"prevPos\":1,\"prevValue\":2},{\"title\":\"b\",\"value\":5,\"prevPos\":0,\"prevValue\":4}]}}",
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
			// check JSON
			{
				buf := new(bytes.Buffer)
				c.formatter.JSON(buf)

				str := buf.String()
				if str != c.json {
					t.Errorf("false JSON formatting:\nhas:\n%v\nwant:\n%v", str, c.json)
				}
			}
		})
	}
}
