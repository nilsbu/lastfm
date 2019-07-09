package rsrc

import (
	"fmt"
	"time"
)

// Day represents a day from midnight to midnight at Greenwich.
// Days before 1970-01-01 are considered undefined.
//
// Midnight returns the beginning of the day as a Unix time stamp.
// Time converts the Day to a time.Time object.
type Day interface {
	fmt.Stringer
	Midnight() (unix int64)
	Time() time.Time
}

// date is a representation of time. It implements Day.
type date time.Time

// ToDay converts a Unix timestamp into a Day. The day is only valid if the
// time stamp is non-negative.
func ToDay(timestamp int64) Day {
	midnight := timestamp - timestamp%86400
	return date(time.Unix(midnight, 0).UTC())
}

// ParseDay parses a date from a string in the format YYYY-MM-DD. It returns nil
// if the string is not valid.
func ParseDay(str string) Day {
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return nil
	}

	return date(t)
}

// DayFromTime converts a Time into a Day.
func DayFromTime(t time.Time) Day {
	return date(time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0, time.UTC))
}

// Midnight returns the Unix timestamp of a date's midnight.
func (d date) Midnight() (unix int64) {
	return time.Time(d).Unix()
}

// Time converts a Date to a time.Time object.
func (d date) Time() time.Time {
	return time.Time(d)
}

func (d date) String() string {
	return time.Time(d).Format("2006-01-02")
}
