package rsrc_test

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestDuration(t *testing.T) {
	for _, c := range []struct {
		name string
		a, b rsrc.Day
		days int
	}{
		{
			"same date",
			rsrc.ParseDay("2022-04-06"),
			rsrc.ParseDay("2022-04-06"),
			0,
		},
		{
			"later",
			rsrc.ParseDay("2022-04-06"),
			rsrc.ParseDay("2022-04-08"),
			2,
		},
		{
			"earlier",
			rsrc.ParseDay("2022-04-08"),
			rsrc.ParseDay("2022-04-01"),
			-7,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			days := rsrc.Between(c.a, c.b).Days()
			if c.days != days {
				t.Errorf("expected %v days but got %v", c.days, days)
			}
		})
	}
}

func TestDurationDays(t *testing.T) {
	for _, c := range []struct {
		name     string
		duration rsrc.Duration
		days     int
	}{
		{
			"1d",
			rsrc.Days(1),
			1,
		},
		{
			"7d",
			rsrc.Days(7),
			7,
		},
		{
			"-700d",
			rsrc.Days(-700),
			-700,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			days := c.duration.Days()
			if c.days != days {
				t.Errorf("expected %v days but got %v", c.days, days)
			}
		})
	}
}

func TestDurationString(t *testing.T) {
	for _, c := range []struct {
		name string
		a, b rsrc.Day
		str  string
	}{
		{
			"same date",
			rsrc.ParseDay("2022-04-06"),
			rsrc.ParseDay("2022-04-06"),
			"0d",
		},
		{
			"later",
			rsrc.ParseDay("2022-04-06"),
			rsrc.ParseDay("2022-04-08"),
			"2d",
		},
		{
			"earlier",
			rsrc.ParseDay("2022-04-08"),
			rsrc.ParseDay("2022-04-01"),
			"-7d",
		},
		{
			"later months",
			rsrc.ParseDay("2022-03-06"),
			rsrc.ParseDay("2022-04-08"),
			"1M 2d",
		},
		{
			"later years",
			rsrc.ParseDay("2012-05-09"),
			rsrc.ParseDay("2022-04-08"),
			"9y 10M 27d",
		},
		{
			"earlier months",
			rsrc.ParseDay("2022-04-08"),
			rsrc.ParseDay("2022-03-06"),
			"-(1M 2d)",
		},
		{
			"earlier years",
			rsrc.ParseDay("2022-04-08"),
			rsrc.ParseDay("2012-05-09"),
			"-(9y 10M 27d)",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			str := rsrc.Between(c.a, c.b).String()
			if c.str != str {
				t.Errorf("expected '%v' but got '%v'", c.str, str)
			}
		})
	}
}
