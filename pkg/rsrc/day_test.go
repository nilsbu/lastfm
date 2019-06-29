package rsrc

import (
	"testing"
	"time"
)

func TestToDay(t *testing.T) {
	cases := []struct {
		day      Day
		midnight int64
	}{
		{ToDay(917740800), 917740800}, // same time
		{ToDay(917741200), 917740800}, // some hours later
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
			day := ParseDay(c.str)

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

func TestDayTime(t *testing.T) {
	cases := []struct {
		day  Day
		time time.Time
	}{
		{ToDay(917740800), time.Unix(917740800, 0).UTC()}, // same time
		{ToDay(917741200), time.Unix(917740800, 0).UTC()}, // some hours later
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
