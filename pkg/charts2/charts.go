package charts2

import "sort"

type LazyCharts interface {
	Column(titles []Title, index int) []float64
	Data(titles []Title, begin, end int) [][]float64

	Titles() []Title
	Len() int
}

type charts struct {
	titles []Title
	values map[string][]float64
}

// Song contains basic information about a song.
type Song struct {
	Artist, Title, Album string
	Duration             float64
}

// Artists compiles LazyCharts in which all songs by an artist are grouped.
func Artists(songs [][]Song) LazyCharts {
	return compileCharts(
		songs,
		func(s Song) Title { return ArtistTitle(s.Artist) },
		func(s Song) float64 { return 1.0 },
	)
}

// ArtistsDuration compiles LazyCharts in which all songs by an artist are
// grouped. The songs are weighted by duration before they are summed up.
func ArtistsDuration(songs [][]Song) LazyCharts {
	return compileCharts(
		songs,
		func(s Song) Title { return ArtistTitle(s.Artist) },
		func(s Song) float64 { return s.Duration },
	)
}

// Songs compiles LazyCharts in which all songs are listed separately.
func Songs(songs [][]Song) LazyCharts {
	return compileCharts(
		songs,
		func(s Song) Title { return SongTitle(s) },
		func(s Song) float64 { return 1.0 },
	)
}

// SongsDuration compiles LazyCharts in which all songs are listed separately.
// The songs are weighted by duration.
func SongsDuration(songs [][]Song) LazyCharts {
	return compileCharts(
		songs,
		func(s Song) Title { return SongTitle(s) },
		func(s Song) float64 { return s.Duration },
	)
}

func compileCharts(
	songs [][]Song,
	key func(Song) Title,
	value func(Song) float64) *charts {
	charts := &charts{
		values: map[string][]float64{},
		titles: []Title{},
	}

	for d, day := range songs {
		for _, song := range day {
			k := key(song)
			if line, ok := charts.values[k.Key()]; ok {
				line[d] += value(song)
			} else {
				charts.titles = append(charts.titles, k)
				charts.values[k.Key()] = make([]float64, len(songs))
				charts.values[k.Key()][d] = value(song)
			}
		}
	}

	return charts
}

func FromMap(data map[string][]float64) LazyCharts {
	titles := make([]Title, len(data))
	i := 0
	for t := range data {
		titles[i] = KeyTitle(t)
		i++
	}
	sort.Slice(titles, func(i, j int) bool { return titles[i].Key() < titles[j].Key() })

	charts := &charts{
		values: data,
		titles: titles,
	}

	return charts
}

func (l *charts) Column(titles []Title, index int) []float64 {
	col := make([]float64, len(titles))
	for i, t := range titles {
		col[i] = l.values[t.Key()][index]
	}
	return col
}

func (l *charts) Data(titles []Title, begin, end int) [][]float64 {
	data := make([][]float64, len(titles))
	for i, t := range titles {
		data[i] = l.values[t.Key()][begin:end]
	}
	return data
}

func (l *charts) Titles() []Title {
	// assumption: noone touches the return value
	return l.titles
}

func (l *charts) Len() int {
	for _, line := range l.values {
		return len(line)
	}
	return -1
}
