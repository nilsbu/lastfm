package command

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nilsbu/lastfm/pkg/unpack"
)

func TestResolve(t *testing.T) {
	cases := []struct {
		args    []string
		session *unpack.SessionInfo
		cmd     command
		ok      bool
	}{
		{
			[]string{},
			nil, nil, false,
		},
		{
			[]string{"grep"},
			nil, nil, false,
		},
		{
			[]string{"lastfm"},
			nil, help{}, true,
		},
		{
			[]string{"lastfm", "asjkdfh"},
			nil, help{}, false,
		},
		{
			[]string{"lastfm", "help"},
			nil, help{}, true,
		},
		// TODO help for commands
		{
			[]string{"lastfm", "session"},
			nil, sessionInfo{}, true,
		},
		{
			[]string{"lastfm", "session", "info"},
			nil, sessionInfo{}, true,
		},
		{
			[]string{"lastfm", "session", "info", "tim"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "session", "asd"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "session", "start"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "session", "start", "tim"},
			&unpack.SessionInfo{User: "tom"},
			sessionStart{user: "tim"}, true,
		},
		{
			[]string{"lastfm", "session", "start", "tim", "xs"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "session", "stop"},
			nil, sessionStop{}, true,
		},
		{
			[]string{"lastfm", "session", "stop", "tim"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "update"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "update"},
			&unpack.SessionInfo{User: "user"},
			updateHistory{}, true,
		},
		{
			[]string{"lastfm", "update", "aargh!"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "print", "total"},
			nil, nil, false,
		},
		{
			[]string{"lastfm", "print", "total"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{by: "all", name: "", n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-n=25"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{by: "all", name: "", n: 25}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-%=TRUE"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{by: "all", name: "", n: 10, percentage: true}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-n=k25"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "total", "-k25"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "total", "-n=10", "-bo=x"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "total", "-by=super", "-n=25"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{by: "super", name: "", n: 25}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=super", "-normalized", "-date=2018-02-01"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{by: "super", name: "", normalized: true, n: 10}, date: time.Date(2018, time.February, 1, 0, 0, 0, 0, time.UTC)}, true,
		},
		{
			[]string{"lastfm", "print", "asdf"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "fade"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "fade", "30.25"},
			&unpack.SessionInfo{User: "user"},
			printFade{printCharts: printCharts{by: "all", name: "", n: 10}, hl: 30.25}, true,
		},
		{
			[]string{"lastfm", "print", "fade", "30.25", "-name=DYD"},
			&unpack.SessionInfo{User: "user"},
			printFade{printCharts: printCharts{by: "all", name: "DYD", n: 10}, hl: 30.25}, true,
		},
		{
			[]string{"lastfm", "print", "fade", "30.25", "-name"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "fade", "10", "-%"},
			&unpack.SessionInfo{User: "user"},
			printFade{printCharts: printCharts{by: "all", n: 10, percentage: true}, hl: 10}, true,
		},
		{
			[]string{"lastfm", "print", "fade", "10", "-normalized=True", "-date=2000-01-01"},
			&unpack.SessionInfo{User: "user"},
			printFade{printCharts: printCharts{by: "all", n: 10, normalized: true}, hl: 10, date: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)}, true,
		},
		{
			[]string{"lastfm", "print", "fade"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "fade", "30", "10", "too many"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "fade", "..."},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "fade", "2", "x"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "period", "2015"},
			&unpack.SessionInfo{User: "user"},
			printPeriod{printCharts: printCharts{by: "all", name: "", n: 10}, period: "2015"}, true,
		},
		{
			[]string{"lastfm", "print", "period", "2015", "-%=1"},
			&unpack.SessionInfo{User: "user"},
			printPeriod{printCharts: printCharts{by: "all", name: "", n: 10, percentage: true}, period: "2015"}, true,
		},
		{
			[]string{"lastfm", "print", "period", "2015", "-normalized=t"},
			&unpack.SessionInfo{User: "user"},
			printPeriod{printCharts: printCharts{by: "all", name: "", n: 10, normalized: true}, period: "2015"}, true,
		},
		{
			[]string{"lastfm", "print", "tags", "Add"},
			&unpack.SessionInfo{User: "user"}, printTags{"Add"}, true,
		},
		{
			[]string{"lastfm", "print", "tags", "Add", "xx"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm-csv", "print", "total"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{by: "all", name: "", n: 10}}, true,
		},
		{
			[]string{"lastfm", "table", "total"},
			&unpack.SessionInfo{User: "user"},
			tableTotal{printCharts: printCharts{by: "all", name: "", n: 10}, step: 1}, true,
		},
	}

	for _, c := range cases {
		str := strings.Join(c.args, " ")
		if c.session != nil {
			str += " (" + string(c.session.User) + ")"
		}

		t.Run(str, func(t *testing.T) {
			cmd, err := resolve(c.args, c.session)

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
