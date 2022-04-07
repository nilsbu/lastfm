package pipeline

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestCache(t *testing.T) {
	// the semantics of the charts are irrelevant for this test,
	// all that matters is that pointers are recognizable.
	cs := make([]charts.Charts, 100)
	for i := range cs {
		cs[i] = charts.FromMap(map[string][]float64{})
	}

	type setget struct {
		method     string
		steps      []string
		charts     charts.Charts
		registered rsrc.Day
	}

	for _, c := range []struct {
		name  string
		limit int
		cmds  []setget
	}{
		{
			"no steps",
			1,
			[]setget{},
		},
		{
			"get without set",
			1,
			[]setget{{
				"get",
				[]string{"sum"},
				nil,
				nil,
			}},
		},
		{
			"set 1st step",
			1,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}},
		},
		{
			"read 1st step",
			1,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}},
		},
		{
			"2 steps",
			1,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"sum", "id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum", "id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}},
		},
		{
			"2 roots",
			2,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}},
		},
		{
			"1st root pruned by cache",
			2,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"id", "cache"},
				cs[2],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				nil,
				nil,
			}},
		},
		{
			"read is saved from pruning",
			2,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"id"},
				cs[1],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				cs[0],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2],
				rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"id", "cache"},
				cs[2],
				rsrc.ParseDay("2000-01-01"),
			}},
		},
		{
			"prune children first",
			3,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"sum", "cache"},
				cs[3], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id"},
				cs[1], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2], rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				cs[0], rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum", "cache"},
				nil, nil,
			}},
		},
		{
			"overwrite can be done",
			3,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"sum"},
				cs[1], rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum"},
				cs[1], rsrc.ParseDay("2000-01-01"),
			}},
		},
		{
			"overwrite keeps cache list intact",
			3,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"sum", "cache"},
				cs[3], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"sum"},
				cs[0], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id"},
				cs[1], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2], rsrc.ParseDay("2000-01-01"),
			}, {
				"get",
				[]string{"sum", "cache"},
				nil, nil,
			}},
		},
		{
			"registered changes",
			3,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0], rsrc.ParseDay("2000-01-01"),
			}, {
				"set",
				[]string{"sum", "sum"},
				cs[1], rsrc.ParseDay("2001-01-01"),
			}, {
				"get",
				[]string{"sum", "sum"},
				cs[1], rsrc.ParseDay("2001-01-01"),
			}},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cache := newGraph(c.limit)

			for i, cmd := range c.cmds {
				var result charts.Charts
				registered := cmd.registered
				switch cmd.method {
				case "get":
					result, registered = cache.get(cmd.steps)

				case "set":
					result = cache.set(cmd.steps, cmd.charts, cmd.registered)

				default:
					t.Fatalf("method '%v' not supported", cmd.method)
				}

				if cmd.charts != result {
					if cmd.charts == nil {
						t.Fatalf("expected result of step %v (%v %v) to be nil but wasn't",
							i, cmd.method, cmd.steps)
					} else if result == nil {
						t.Fatalf("expected result of step %v (%v %v) to not be nil but was",
							i, cmd.method, cmd.steps)
					} else {
						t.Fatalf("result of step %v (%v %v) is wrong but not nil",
							i, cmd.method, cmd.steps)
					}
				}
				if cmd.registered != nil && registered == nil {
					t.Fatal("registered must not be nil")
				} else if cmd.registered == nil && registered != nil {
					t.Fatal("registered must be nil but is", registered)
				} else if cmd.registered != nil && rsrc.Between(cmd.registered, registered).Days() != 0 {
					t.Fatalf("expected registered %v but was %v", cmd.registered, registered)
				}
			}
		})
	}
}
