package command

import (
	"reflect"
	"testing"
)

func TestResolve(t *testing.T) {
	cases := []struct {
		args []string
		cmd  command
		ok   bool
	}{
		{
			[]string{},
			nil, false},
		{
			[]string{"grep"},
			nil, false,
		},
		{
			[]string{"lastfm"},
			help{}, true,
		},
		{
			[]string{"lastfm", "asjkdfh"},
			help{}, false,
		},
		{
			[]string{"lastfm", "help"},
			help{}, true,
		},
		// TODO help for commands
		{
			[]string{"lastfm", "session"},
			sessionInfo{}, true,
		},
		{
			[]string{"lastfm", "session", "info"},
			sessionInfo{}, true,
		},
		{
			[]string{"lastfm", "session", "info", "tim"},
			nil, false,
		},
		{
			[]string{"lastfm", "session", "asd"},
			nil, false,
		},
		{
			[]string{"lastfm", "session", "start"},
			nil, false,
		},
		{
			[]string{"lastfm", "session", "start", "tim"},
			sessionStart{user: "tim"}, true,
		},
		{
			[]string{"lastfm", "session", "start", "tim", "xs"},
			nil, false,
		},
		{
			[]string{"lastfm", "session", "stop"},
			sessionStop{}, true,
		},
		{
			[]string{"lastfm", "session", "stop", "tim"},
			nil, false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			cmd, err := resolve(c.args)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				if !reflect.DeepEqual(cmd, c.cmd) {
					t.Errorf("resolve() returned wrong command:\nhas:      %v (%v)\nexpected: %v (%v)",
						cmd, reflect.TypeOf(cmd), c.cmd, reflect.TypeOf(c.cmd))
				}
			}
		})
	}

}
