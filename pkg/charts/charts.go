package charts

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Key interface {
	fmt.Stringer
}

type simpleKey string

func (s simpleKey) String() string {
	return string(s)
}

type Charts struct {
	Headers Intervals
	Keys    []Key
	Values  [][]float64
}

func CompileArtists(
	days []map[string]float64,
	registered rsrc.Day) Charts {
	size := len(days)

	keys := []Key{}
	values := [][]float64{}

	charts := make(map[string]int)
	for i, day := range days {
		for name, plays := range day {
			if _, ok := charts[name]; !ok {
				charts[name] = len(values)
				keys = append(keys, simpleKey(name))
				values = append(values, make([]float64, size))
			}
			values[charts[name]][i] = plays
		}
	}

	end := rsrc.ToDay(registered.Midnight() + int64(86400*size))
	return Charts{
		Headers: Days(registered, end),
		Keys:    keys,
		Values:  values,
	}
}

// UnravelDays takes Charts and disassembles it into single day plays. It acts
// as an inverse to Compile().
func (c Charts) UnravelDays() []map[string]float64 {
	days := []map[string]float64{}
	for i := 0; i < c.Len(); i++ {
		day := map[string]float64{}

		for j, line := range c.Values {
			if line[i] != 0 {
				day[c.Keys[j].String()] = line[i]
			}
		}

		days = append(days, day)
	}

	return days
}

func (c Charts) Len() int {
	if len(c.Values) == 0 {
		return 0
	}

	return len(c.Values[0])
}

// GetKeys returns the keys of the charts.
func (c Charts) GetKeys() []string {
	keys := []string{}
	for _, key := range c.Keys {
		keys = append(keys, key.String())
	}
	return keys
}
