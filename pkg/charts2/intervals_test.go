package charts2_test

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts2"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestParseRange(t *testing.T) {
	for _, c := range []struct {
		name   string
		str    string
		l      int
		result charts2.Range
		ok     bool
	}{
		{
			"year",
			"2022",
			5 * 365,
			charts2.Range{
				rsrc.ParseDay("2022-01-01"),
				rsrc.ParseDay("2023-01-01"),
				rsrc.ParseDay("2019-01-01"),
			},
			true,
		},
		{
			"month",
			"2022-02",
			5 * 365,
			charts2.Range{
				rsrc.ParseDay("2022-02-01"),
				rsrc.ParseDay("2022-03-01"),
				rsrc.ParseDay("2019-01-01"),
			},
			true,
		},
		{
			"day",
			"2022-02-28",
			5 * 365,
			charts2.Range{
				rsrc.ParseDay("2022-02-28"),
				rsrc.ParseDay("2022-03-01"),
				rsrc.ParseDay("2019-01-01"),
			},
			true,
		},
		{
			"registered in the middle of the year",
			"2022",
			5 * 365,
			charts2.Range{
				rsrc.ParseDay("2022-07-01"),
				rsrc.ParseDay("2023-01-01"),
				rsrc.ParseDay("2022-07-01"),
			},
			true,
		},
		{
			"len shorter than year",
			"2022",
			31,
			charts2.Range{
				rsrc.ParseDay("2022-01-01"),
				rsrc.ParseDay("2022-02-01"),
				rsrc.ParseDay("2022-01-01"),
			},
			true,
		},
		{
			"string bs",
			"202",
			31,
			charts2.Range{
				rsrc.ParseDay("2022-01-01"),
				rsrc.ParseDay("2022-02-01"),
				rsrc.ParseDay("2022-01-01"),
			},
			false,
		},
		{
			"string bs 2",
			"20x2",
			31,
			charts2.Range{
				rsrc.ParseDay("2022-01-01"),
				rsrc.ParseDay("2022-02-01"),
				rsrc.ParseDay("2022-01-01"),
			},
			false,
		},
		{
			"begin after end of data",
			"2024",
			31,
			charts2.Range{
				rsrc.ParseDay("2024-01-01"),
				rsrc.ParseDay("2024-02-01"),
				rsrc.ParseDay("2022-01-01"),
			},
			false,
		},
		{
			"end before registered",
			"2022",
			31,
			charts2.Range{
				rsrc.ParseDay("2022-01-01"),
				rsrc.ParseDay("2022-02-01"),
				rsrc.ParseDay("2023-01-01"),
			},
			false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			r, err := charts2.ParseRange(c.str, c.result.Registered, c.l)
			if (err == nil) != c.ok {
				t.Fatalf("error: %v, ok status expected: %v",
					err, c.ok)
			}
			if err == nil {
				if r.Begin != c.result.Begin {
					t.Errorf("begin is wrong: %v != %v",
						r.Begin, c.result.Begin)
				}

				if r.End != c.result.End {
					t.Errorf("end is wrong: %v != %v",
						r.End, c.result.End)
				}

				if r.Registered != c.result.Registered {
					t.Errorf("registered is wrong: %v != %v",
						r.Registered, c.result.Registered)
				}
			}
		})
	}
}
