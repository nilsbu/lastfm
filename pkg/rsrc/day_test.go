package rsrc_test

import (
	"testing"
	"time"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestToDay(t *testing.T) {
	cases := []struct {
		day      rsrc.Day
		midnight int64
	}{
		{rsrc.ToDay(917740800), 917740800}, // same time
		{rsrc.ToDay(917741200), 917740800}, // some hours later
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			midnight := c.day.Midnight()

			if midnight != c.midnight {
				t.Errorf("got midnight '%v', expected '%v'",
					midnight, c.midnight)
			}
		})
	}
}

func TestParseDay(t *testing.T) {
	cases := []struct {
		str      string
		midnight int64
		ok       bool
	}{
		{"2017-04-04", 1491264000, true},
		{"2017-04-04T", 0, false},
		{"2017-02-30", 0, false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			day := rsrc.ParseDay(c.str)

			if day == nil && c.ok {
				t.Errorf("valid result was expected but 'nil' was returned")
			} else if day != nil && !c.ok {
				t.Error("expected error but none occurred")
			}
			if c.ok {
				midnight := day.Midnight()

				if midnight != c.midnight {
					t.Errorf("got midnight '%v', expected '%v'",
						midnight, c.midnight)
				}
			}
		})
	}
}

func TestDayFromTime(t *testing.T) {
	cases := []struct {
		time     time.Time
		midnight int64
	}{
		{time.Date(2017, 4, 4, 0, 0, 0, 0, time.UTC), 1491264000},
		{time.Date(2017, 4, 4, 12, 6, 2, 0, time.UTC), 1491264000},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			day := rsrc.DayFromTime(c.time)

			if day == nil {
				t.Errorf("valid result was expected but 'nil' was returned")
			}

			midnight := day.Midnight()

			if midnight != c.midnight {
				t.Errorf("got midnight '%v', expected '%v'",
					midnight, c.midnight)
			}
		})
	}
}

func TestDayTime(t *testing.T) {
	cases := []struct {
		day  rsrc.Day
		time time.Time
	}{
		{rsrc.ToDay(917740800), time.Unix(917740800, 0).UTC()}, // same time
		{rsrc.ToDay(917741200), time.Unix(917740800, 0).UTC()}, // some hours later
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			time := c.day.Time()

			if time != c.time {
				t.Errorf("got '%v', expected '%v'",
					time, c.time)
			}
		})
	}
}

func TestDayAddDate(t *testing.T) {
	cases := []struct {
		base      rsrc.Day
		addYears  int
		addMonths int
		addDays   int
		sum       rsrc.Day
	}{
		{
			rsrc.ParseDay("1992-11-12"),
			0, 0, 0,
			rsrc.ParseDay("1992-11-12"),
		},
		{
			rsrc.ParseDay("1992-11-12"),
			1, 2, 1,
			rsrc.ParseDay("1994-01-13"),
		},
		{
			rsrc.ParseDay("1992-11-12"),
			1, -3, 0,
			rsrc.ParseDay("1993-08-12"),
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sum := c.base.AddDate(
				c.addYears, c.addMonths, c.addDays)

			if sum.Midnight() != c.sum.Midnight() {
				t.Errorf("got '%v', expected '%v'",
					sum, c.sum)
			}
		})
	}
}

func TestDayString(t *testing.T) {
	cases := []struct {
		day rsrc.Day
		str string
	}{
		{rsrc.ToDay(1491264000), "2017-04-04"},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			str := c.day.String()

			if str != c.str {
				t.Errorf("got '%v', expected '%v'",
					str, c.str)
			}
		})
	}
}
