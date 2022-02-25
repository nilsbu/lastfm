package charts

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Key interface {
	fmt.Stringer
	ArtistName() string
	FullTitle() string
}

type simpleKey string

func (s simpleKey) String() string {
	return string(s)
}

func (s simpleKey) ArtistName() string {
	return string(s)
}

func (s simpleKey) FullTitle() string {
	return string(s)
}

type tagKey string

func (s tagKey) String() string {
	return string(s)
}

func (s tagKey) ArtistName() string {
	return ""
}

func (s tagKey) FullTitle() string {
	return string(s)
}

type Song struct {
	Artist   string
	Title    string
	Album    string
	Duration float64
}

func (s Song) String() string {
	return s.Artist + " - " + s.Title
}

func (s Song) ArtistName() string {
	return s.Artist
}

func (s Song) FullTitle() string {
	return s.Artist + " - " + s.Title
}

type Charts struct {
	Headers Intervals
	Keys    []Key
	Values  [][]float64
}

func CompileArtists(
	days []map[string]float64,
	registered rsrc.Day) Charts {
	return compile(days, registered, func(s string) Key { return simpleKey(s) })
}

func CompileTags(
	days []map[string]float64,
	registered rsrc.Day) Charts {
	return compile(days, registered, func(s string) Key { return tagKey(s) })
}

func compile(
	days []map[string]float64,
	registered rsrc.Day,
	toKey func(string) Key) Charts {
	size := len(days)

	keys := []Key{}
	values := [][]float64{}

	charts := make(map[string]int)
	for i, day := range days {
		for name, plays := range day {
			if _, ok := charts[name]; !ok {
				charts[name] = len(values)
				keys = append(keys, toKey(name))
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

// CompileSongs creates charts of individual songs from songs split by days.
func CompileSongs(
	days [][]Song,
	registered rsrc.Day) Charts {
	return chartsFromSongs(
		days,
		registered,
		func(s Song) Key { return s })
}

// ArtistsFromSongs creates artist charts from songs split by days.
func ArtistsFromSongs(
	days [][]Song,
	registered rsrc.Day,
) Charts {
	return chartsFromSongs(
		days,
		registered,
		func(s Song) Key { return simpleKey(s.Artist) })
}

func chartsFromSongs(
	days [][]Song,
	registered rsrc.Day,
	getKey func(s Song) Key,
) Charts {
	size := len(days)

	keys := []Key{}
	values := [][]float64{}

	charts := make(map[string]int)
	for i, day := range days {
		for _, song := range day {
			k := getKey(song)
			key := k.String()
			if _, ok := charts[key]; !ok {
				charts[key] = len(values)
				keys = append(keys, k)
				values = append(values, make([]float64, size))
				values[len(values)-1][i] = 1
			} else {
				values[charts[key]][i]++
			}
		}
	}

	return Charts{
		Headers: Days(registered, registered.AddDate(0, 0, size)),
		Keys:    keys,
		Values:  values,
	}
}

// UnravelDays takes Charts and disassembles it into single day plays. It acts
// as an inverse to CompileArtists().
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

// UnravelSongs takes Charts and disassembles it into songs. It acts as an
// inverse to CompileSongs().
func (c Charts) UnravelSongs() [][]Song {
	songs := make([][]Song, c.Len())
	for d := range songs {
		day := []Song{}
		for k, key := range c.Keys {
			for n := 0; n < int(c.Values[k][d]); n++ {
				if song, ok := key.(Song); ok {
					day = append(day, song)
				} else {
					day = append(day, Song{
						Artist: key.ArtistName(),
						Title:  "",
						Album:  ""})
				}
			}
		}
		songs[d] = day
	}
	return songs
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
