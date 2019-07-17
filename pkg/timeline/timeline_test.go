package timeline

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestNumberOne(t *testing.T) {
	cases := []struct {
		descr    string
		charts   charts.Charts
		interval charts.Interval
		events   []Event
	}{
		{
			"empty charts",
			charts.CompileArtists(
				[]map[string]float64{},
				rsrc.ParseDay("2010-01-01")),
			charts.Interval{
				Begin:  rsrc.ParseDay("2010-01-01"),
				Before: rsrc.ParseDay("2010-01-02")},
			[]Event{},
		},
		{
			"end before registration date",
			charts.CompileArtists(
				[]map[string]float64{
					{"A": 1},
					{"A": 1},
				},
				rsrc.ParseDay("2010-01-01")),
			charts.Interval{
				Begin:  rsrc.ParseDay("2009-01-01"),
				Before: rsrc.ParseDay("2009-12-02")},
			[]Event{},
		},
		{
			"only entry date",
			charts.CompileArtists(
				[]map[string]float64{
					{"A": 1},
					{"A": 1},
				},
				rsrc.ParseDay("2010-01-01")),
			charts.Interval{
				Begin:  rsrc.ParseDay("2010-01-01"),
				Before: rsrc.ParseDay("2010-01-02")},
			[]Event{
				{rsrc.ParseDay("2010-01-01"), "top at begin is 'A'"},
			},
		},
		{
			"with changes",
			charts.CompileArtists(
				[]map[string]float64{
					{"A": 1, "B": 0, "C": 0},
					{"A": 1, "B": 2, "C": 0},
					{"A": 1, "B": 2, "C": 1},
					{"A": 1, "B": 2, "C": 4},
					{"A": 9, "B": 5, "C": 4},
					{"A": 9, "B": 17, "C": 4},
					{"A": 20, "B": 17, "C": 4},
				},
				rsrc.ParseDay("2010-01-01")),
			charts.Interval{
				Begin:  rsrc.ParseDay("2010-01-02"),
				Before: rsrc.ParseDay("2010-01-06")},
			[]Event{
				{rsrc.ParseDay("2010-01-02"), "'A' -> 'B' (1d)"},
				{rsrc.ParseDay("2010-01-04"), "'B' -> 'C' (2d)"},
				{rsrc.ParseDay("2010-01-05"), "'C' -> 'A' (1d)"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.descr, func(t *testing.T) {
			events := CompileNumberOne(c.charts, c.interval)

			if len(c.events) != len(events) {
				t.Fatalf("expect %v events but got %v",
					len(c.events), len(events))
			}

			for i, event := range c.events {
				if event.Date.Midnight() != events[i].Date.Midnight() {
					t.Errorf("at position %v: expected date %v but got %v",
						i, event.Date, events[i].Date)
				}
				if event.Message != events[i].Message {
					t.Errorf("at position %v: expected message '%v' but got '%v'",
						i, event.Message, events[i].Message)
				}
			}
		})
	}
}
