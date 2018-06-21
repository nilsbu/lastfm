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

type date time.Time

// ToDay converts a Unix timestamp into a Day. The day is only valid if the
// time stamp is non-negative.
func ToDay(timestamp int64) Day {
	return date(time.Unix(timestamp, 0).UTC())
}

// NoDay returns an invalid Day.
func NoDay() Day {
	return ToDay(-1)
}

func (d date) Midnight() (unix int64, ok bool) {
	t := time.Time(d).Unix()
	if t < 0 {
		return -1, false
	}
	return t - t%86400, true
}
