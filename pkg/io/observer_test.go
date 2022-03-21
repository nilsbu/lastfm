package io

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

type event struct {
	t   rune
	loc rsrc.Locator
}

func TestObserver(t *testing.T) {
	cases := []struct {
		name   string
		events []event
		msgs   []format.Formatter
	}{
		{
			"no notify",
			[]event{},
			[]format.Formatter{},
		},
		{
			"1 write",
			[]event{
				{'w', rsrc.APIKey()},
				{'W', rsrc.APIKey()},
			},
			[]format.Formatter{
				&format.Message{Msg: "r: 0/0, w: 0/1, rm: 0/0"},
				&format.Message{Msg: "r: 0/0, w: 1/1, rm: 0/0"},
			},
		},
		{
			"read, remove",
			[]event{
				{'r', rsrc.APIKey()},
				{'d', rsrc.APIKey()},
				{'D', rsrc.APIKey()},
				{'R', rsrc.APIKey()},
			},
			[]format.Formatter{
				&format.Message{Msg: "r: 0/1, w: 0/0, rm: 0/0"},
				&format.Message{Msg: "r: 0/1, w: 0/0, rm: 0/1"},
				&format.Message{Msg: "r: 0/1, w: 0/0, rm: 1/1"},
				&format.Message{Msg: "r: 1/1, w: 0/0, rm: 1/1"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fChan := make(chan format.Formatter)
			quit := make(chan bool)
			d := mock.NewDisplay()
			go func(d display.Display) {
				for msg := range fChan {
					d.Display(msg)
				}
				quit <- true
			}(d)

			o := newObserver(fChan)
			for _, e := range c.events {
				switch e.t {
				case 'r':
					o.RequestRead(e.loc)
				case 'R':
					o.NotifyRead(e.loc)
				case 'w':
					o.RequestWrite(e.loc)
				case 'W':
					o.NotifyWrite(e.loc)
				case 'd':
					o.RequestRemove(e.loc)
				case 'D':
					o.NotifyRemove(e.loc)
				}
			}
			close(fChan)
			<-quit

			if len(c.msgs) != len(d.Msgs) {
				t.Fatalf("expected %v messages but got %v", len(c.msgs), len(d.Msgs))
			}
			for i, expect := range c.msgs {
				if !reflect.DeepEqual(expect, d.Msgs[i]) {
					t.Errorf("expect %v, got %v", expect, d.Msgs[i])
				}
			}
		})
	}
}
