package charts2

type charts struct {
	// Headers charts.Interval
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
