package command

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nilsbu/lastfm/pkg/organize"
)

func TestResolve(t *testing.T) {
	cases := []struct {
		args []string
		sid  organize.SessionID
		cmd  command
		ok   bool
	}{
		{
			[]string{},
			"", nil, false,
		},
		{
			[]string{"grep"},
			"", nil, false,
		},
		{
			[]string{"lastfm"},
			"", help{}, true,
		},
		{
			[]string{"lastfm", "asjkdfh"},
			"", help{}, false,
		},
		{
			[]string{"lastfm", "help"},
			"", help{}, true,
		},
		// TODO help for commands
		{
			[]string{"lastfm", "session"},
			"", sessionInfo{}, true,
		},
		{
			[]string{"lastfm", "session", "info"},
			"", sessionInfo{}, true,
		},
		{
			[]string{"lastfm", "session", "info"},
			"xs", sessionInfo{"xs"}, true,
		},
		{
			[]string{"lastfm", "session", "info", "tim"},
			"", nil, false,
		},
		{
			[]string{"lastfm", "session", "asd"},
			"", nil, false,
		},
		{
			[]string{"lastfm", "session", "start"},
			"", nil, false,
		},
		{
			[]string{"lastfm", "session", "start", "tim"},
			"tom", sessionStart{sid: "tom", user: "tim"}, true,
		},
		{
			[]string{"lastfm", "session", "start", "tim", "xs"},
			"", nil, false,
		},
		{
			[]string{"lastfm", "session", "stop"},
			"", sessionStop{""}, true,
		},
		{
			[]string{"lastfm", "session", "stop", "tim"},
			"", nil, false,
		},
		{
			[]string{"lastfm", "update"},
			"", nil, false,
		},
		{
			[]string{"lastfm", "update"},
			"user", updateHistory{"user"}, true,
		},
		{
			[]string{"lastfm", "update", "aargh!"},
			"user", nil, false,
		},
	}

	for _, c := range cases {
		str := strings.Join(c.args, " ")
		if c.sid != "" {
			str += " (" + string(c.sid) + ")"
		}

		t.Run(str, func(t *testing.T) {
			cmd, err := resolve(c.args, c.sid)

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
