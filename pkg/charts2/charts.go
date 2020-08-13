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
		songs, func(s Song) string { return s.Artist })
}

func compileCharts(songs [][]Song, key func(Song) string) *charts {
	charts := &charts{
		values: map[string][]float64{},
		titles: []Title{},
	}

	for d, day := range songs {
		for _, song := range day {
			k := key(song)
			if line, ok := charts.values[k]; ok {
				line[d]++
			} else {
				charts.titles = append(charts.titles, ArtistTitle(song.Artist))
				charts.values[k] = make([]float64, len(songs))
				charts.values[k][d] = 1
			}
		}
	}

	return charts
}
