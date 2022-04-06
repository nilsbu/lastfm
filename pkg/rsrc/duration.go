package rsrc

import (
	"fmt"
	"time"
)

type Duration interface {
	fmt.Stringer
	Days() int
}

type duration struct {
	a, b time.Time
}

func Between(a, b Day) Duration {
	return &duration{a: a.Time(), b: b.Time()}
}

func (d *duration) Days() int {
	return int(d.b.Sub(d.a).Hours() / 24)
}

func (d *duration) String() string {
	inverted, years, months, days, _, _, _, _ := Elapsed(d.a, d.b)

	if years != 0 {
		if inverted {
			return fmt.Sprintf("-(%vy %vM %vd)", years, months, days)
		} else {
			return fmt.Sprintf("%vy %vM %vd", years, months, days)
		}
	} else if months != 0 {
		if inverted {
			return fmt.Sprintf("-(%vM %vd)", months, days)
		} else {
			return fmt.Sprintf("%vM %vd", months, days)
		}
	} else if days != 0 {
		if inverted {
			return fmt.Sprintf("-%vd", days)
		} else {
			return fmt.Sprintf("%vd", days)
		}
	}
	return "0d"
}

func DaysIn(year int, month time.Month) int {
	return time.Date(year, month, 0, 0, 0, 0, 0, time.UTC).Day()
}

func Elapsed(from, to time.Time) (inverted bool, years, months, days, hours, minutes, seconds, nanoseconds int) {
	if from.Location() != to.Location() {
		to = to.In(to.Location())
	}

	inverted = false
	if from.After(to) {
		inverted = true
		from, to = to, from
	}

	y1, M1, d1 := from.Date()
	y2, M2, d2 := to.Date()

	h1, m1, s1 := from.Clock()
	h2, m2, s2 := to.Clock()

	ns1, ns2 := from.Nanosecond(), to.Nanosecond()

	years = y2 - y1
	months = int(M2 - M1)
	days = d2 - d1

	hours = h2 - h1
	minutes = m2 - m1
	seconds = s2 - s1
	nanoseconds = ns2 - ns1

	if nanoseconds < 0 {
		nanoseconds += 1e9
		seconds--
	}
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		days += DaysIn(y2, M2-1)
		months--
	}
	if days < 0 {
		days += DaysIn(y2, M2)
		months--
	}
	if months < 0 {
		months += 12
		years--
	}
	return
}
