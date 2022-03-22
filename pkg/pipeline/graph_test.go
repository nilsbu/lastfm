package pipeline

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

func TestCache(t *testing.T) {
	// the semantics of the charts are irrelevant for this test,
	// all that matters is that pointers are recognizable.
	cs := make([]charts.Charts, 100)
	for i := range cs {
		cs[i] = charts.FromMap(map[string][]float64{})
	}

	type setget struct {
		method string
		steps  []string
		charts charts.Charts
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
			}},
		},
		{
			"set 1st step",
			1,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}},
		},
		{
			"read 1st step",
			1,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"get",
				[]string{"sum"},
				cs[0],
			}},
		},
		{
			"2 steps",
			1,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"sum", "id"},
				cs[1],
			}, {
				"get",
				[]string{"sum", "id"},
				cs[1],
			}},
		},
		{
			"2 roots",
			2,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"id"},
				cs[1],
			}, {
				"get",
				[]string{"id"},
				cs[1],
			}, {
				"get",
				[]string{"sum"},
				cs[0],
			}},
		},
		{
			"1st root pruned by cache",
			2,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"get",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"id"},
				cs[1],
			}, {
				"get",
				[]string{"id"},
				cs[1],
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2],
			}, {
				"get",
				[]string{"id", "cache"},
				cs[2],
			}, {
				"get",
				[]string{"sum"},
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
			}, {
				"get",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"id"},
				cs[1],
			}, {
				"get",
				[]string{"id"},
				cs[1],
			}, {
				"get",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2],
			}, {
				"get",
				[]string{"id", "cache"},
				cs[2],
			}},
		},
		{
			"prune children first",
			3,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"sum", "cache"},
				cs[3],
			}, {
				"set",
				[]string{"id"},
				cs[1],
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2],
			}, {
				"get",
				[]string{"sum"},
				cs[0],
			}, {
				"get",
				[]string{"sum", "cache"},
				nil,
			}},
		},
		{
			"overwrite can be done",
			3,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"sum"},
				cs[1],
			}, {
				"get",
				[]string{"sum"},
				cs[1],
			}},
		},
		{
			"overwrite keeps cache list intact",
			3,
			[]setget{{
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"sum", "cache"},
				cs[3],
			}, {
				"set",
				[]string{"sum"},
				cs[0],
			}, {
				"set",
				[]string{"id"},
				cs[1],
			}, {
				"set",
				[]string{"id", "cache"},
				cs[2],
			}, {
				"get",
				[]string{"sum", "cache"},
				nil,
			}},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cache := newGraph(c.limit)

			for i, cmd := range c.cmds {
				var result charts.Charts
				switch cmd.method {
				case "get":
					result = cache.get(cmd.steps)

				case "set":
					result = cache.set(cmd.steps, cmd.charts)

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
			}
		})
	}
}
