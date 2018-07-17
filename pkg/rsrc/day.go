package rsrc

import "time"

// Day represents a day from midnight to midnight at Greenwich.
// Days before 1970-01-01 are considered undefined.
//
// Midnight returns the beginning of the day as a Unix time stamp.
// The value is ok, if the Day actually corresponds to a valid day with
// non-negative Unix time stamp (not before 1970).
type Day interface {
	Midnight() (unix int64, ok bool)
}

// Date is a representation of time. It implements Day.
type Date time.Time

// ToDay converts a Unix timestamp into a Day. The day is only valid if the
// time stamp is non-negative.
func ToDay(timestamp int64) Day {
	return Date(time.Unix(timestamp, 0).UTC())
}

func ParseDay(date string) Day {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return NoDay()
	}

	return Date(t)
}

// NoDay returns an invalid Day.
func NoDay() Day {
	return ToDay(-1)
}

// Midnight returns the Unix timestamp of a date's midnight.
func (d Date) Midnight() (unix int64, ok bool) {
	t := time.Time(d).Unix()
	if t < 0 {
		return -1, false
	}
	return t - t%86400, true
}
