package charts

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Song struct {
	Artist string
	Title  string
	Album  string
}

type Key interface {
	fmt.Stringer
	Artist() string
	FullTitle() string
}

type simpleKey string

func (s simpleKey) String() string {
	return string(s)
}

func (s simpleKey) Artist() string {
	return string(s)
}

func (s simpleKey) FullTitle() string {
	return string(s)
}

type tagKey string

func (s tagKey) String() string {
	return string(s)
}

func (s tagKey) Artist() string {
	return ""
}

func (s tagKey) FullTitle() string {
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

	return Charts{
		Headers: Days(registered, registered.AddDate(0, 0, size)),
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
