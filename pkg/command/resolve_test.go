package command

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
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
			[]string{"lastfm", "session", "config", "n", "42"},
			&unpack.SessionInfo{User: "tom"},
			sessionConfig{option: "n", value: "42"}, true,
		},
		{
			[]string{"lastfm", "session", "config", "n"},
			&unpack.SessionInfo{User: "tom"},
			nil, false,
		},
		{
			[]string{"lastfm", "session", "config", "n", "42", "n"},
			&unpack.SessionInfo{User: "tom"},
			nil, false,
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
			printTotal{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-n=25"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 25}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-%=TRUE"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10, percentage: true}}, true,
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
			printTotal{printCharts: printCharts{keys: "artist", by: "super", name: "", n: 25}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=super", "-normalized", "-date=2018-02-01"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "super", name: "", normalized: true, n: 10}, date: rsrc.ParseDay("2018-02-01")}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=year"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "year", name: "", n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=year", "-name=2018"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "year", name: "2018", n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=year", "-entry=60"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "year", entry: 60, n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=year", "-entry=60", "-keys=song"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "song", by: "year", entry: 60, n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=country"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "country", n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-by=total"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "total", name: "", n: 10}}, true,
		},
		{
			[]string{"lastfm", "print", "total", "-duration"},
			&unpack.SessionInfo{User: "user"},
			printTotal{printCharts: printCharts{keys: "artist", by: "all", duration: true, n: 10}}, true,
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
			printFade{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, hl: 30.25}, true,
		},
		{
			[]string{"lastfm", "print", "fade", "30.25", "-name=DYD"},
			&unpack.SessionInfo{User: "user"},
			printFade{printCharts: printCharts{keys: "artist", by: "all", name: "DYD", n: 10}, hl: 30.25}, true,
		},
		{
			[]string{"lastfm", "print", "fade", "30.25", "-name"},
			&unpack.SessionInfo{User: "user"}, nil, false,
		},
		{
			[]string{"lastfm", "print", "fade", "10", "-%"},
			&unpack.SessionInfo{User: "user"},
			printFade{printCharts: printCharts{keys: "artist", by: "all", n: 10, percentage: true}, hl: 10}, true,
		},
		{
			[]string{"lastfm", "print", "fade", "10", "-normalized=True", "-date=2000-01-01"},
			&unpack.SessionInfo{User: "user"},
			printFade{printCharts: printCharts{keys: "artist", by: "all", n: 10, normalized: true}, hl: 10, date: rsrc.ParseDay("2000-01-01")}, true,
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
			printPeriod{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, period: "2015"}, true,
		},
		{
			[]string{"lastfm", "print", "period", "2015", "-%=1"},
			&unpack.SessionInfo{User: "user"},
			printPeriod{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10, percentage: true}, period: "2015"}, true,
		},
		{
			[]string{"lastfm", "print", "period", "2015", "-normalized=t"},
			&unpack.SessionInfo{User: "user"},
			printPeriod{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10, normalized: true}, period: "2015"}, true,
		},
		{
			[]string{"lastfm", "print", "interval", "2007-01-01", "2018-12-24"},
			&unpack.SessionInfo{User: "user"},
			printInterval{
				printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10, normalized: false},
				begin:       rsrc.ParseDay("2007-01-01"),
				end:         rsrc.ParseDay("2018-12-24")}, true,
		},
		{
			[]string{"lastfm", "print", "fademax", "66"},
			&unpack.SessionInfo{User: "user"},
			printFadeMax{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10, percentage: false}, hl: 66}, true,
		},
		{
			[]string{"lastfm", "print", "periods", "3y"},
			&unpack.SessionInfo{User: "user"},
			printPeriods{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10, percentage: false}, period: "3y", begin: rsrc.ParseDay("0001-01-01"), end: rsrc.ParseDay("9999-12-31")}, true,
		},
		{
			[]string{"lastfm", "print", "periods", "3y", "-keys=song", "-end=2022-01-01"},
			&unpack.SessionInfo{User: "user"},
			printPeriods{printCharts: printCharts{keys: "song", by: "all", name: "", n: 10, percentage: false}, period: "3y", begin: rsrc.ParseDay("0001-01-01"), end: rsrc.ParseDay("2022-01-01")}, true,
		},
		{
			[]string{"lastfm", "print", "fades", "333", "3y", "-keys=song", "-end=2022-01-01"},
			&unpack.SessionInfo{User: "user"},
			printFades{printCharts: printCharts{keys: "song", by: "all", name: "", n: 10, percentage: false}, hl: 333, period: "3y", begin: rsrc.ParseDay("0001-01-01"), end: rsrc.ParseDay("2022-01-01")}, true,
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
			printTotal{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}}, true,
		},
		{
			[]string{"lastfm", "table", "total"},
			&unpack.SessionInfo{User: "user"},
			tableTotal{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, step: 1}, true,
		},
		{
			[]string{"lastfm", "table", "total", "-step=200"},
			&unpack.SessionInfo{User: "user"},
			tableTotal{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, step: 200}, true,
		},
		{
			[]string{"lastfm", "table", "period", "1y"},
			&unpack.SessionInfo{User: "user"},
			tablePeriods{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, period: "1y"}, true,
		},
		// {
		// 	[]string{"lastfm", "timeline", "-before=2008-01-23", "-from=2000-11-03"},
		// 	&unpack.SessionInfo{User: "user"},
		// 	printTimeline{
		// 		from:   time.Date(2000, time.Month(11), 03, 0, 0, 0, 0, time.UTC),
		// 		before: time.Date(2008, time.Month(1), 23, 0, 0, 0, 0, time.UTC),
		// 	}, true,
		// },
		{
			[]string{"lastfm-csv", "table", "fade", "10"},
			&unpack.SessionInfo{User: "user"},
			tableFade{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, hl: 10, step: 1}, true,
		},
		{
			// relevant option stored
			[]string{"lastfm-csv", "table", "fade", "10"},
			&unpack.SessionInfo{User: "user", Options: map[string]string{"step": "30"}},
			tableFade{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, hl: 10, step: 30}, true,
		},
		{
			// irrelevant option stored
			[]string{"lastfm", "print", "total"},
			&unpack.SessionInfo{User: "user", Options: map[string]string{"step": "30"}},
			printTotal{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}}, true,
		},
		{
			// explicit parameter overrides stored parameter
			[]string{"lastfm-csv", "table", "fade", "10", "-step=25"},
			&unpack.SessionInfo{User: "user", Options: map[string]string{"step": "30"}},
			tableFade{printCharts: printCharts{keys: "artist", by: "all", name: "", n: 10}, hl: 10, step: 25}, true,
		},
		{
			// raw steps
			[]string{"lastfm", "print", "raw", "artistsduration", "sum", "top,69"},
			&unpack.SessionInfo{User: "user"},
			printRaw{precision: 2, steps: []string{"artistsduration", "sum", "top,69"}}, true,
		},
	}

	for i, c := range cases {
		str := fmt.Sprintf("%v - %v", i, strings.Join(c.args, " "))
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
