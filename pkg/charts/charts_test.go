package charts_test

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestCharts(t *testing.T) {
	// The purpose here is to ensure the content of the charts is correct, the
	// API is tested in TestLazyCharts in greater detail.
	for _, c := range []struct {
		name   string
		charts charts.LazyCharts
		titles []charts.Title
		lines  [][]float64
	}{
		{
			"Artists",
			charts.Artists([][]charts.Song{
				{charts.Song{Artist: "A", Duration: 1},
					charts.Song{Artist: "B", Duration: 2},
				},
				{charts.Song{Artist: "C", Duration: 1},
					charts.Song{Artist: "A", Duration: 1},
				},
			}),
			[]charts.Title{charts.ArtistTitle("A"), charts.ArtistTitle("B"), charts.ArtistTitle("C")},
			[][]float64{
				{1, 1}, {1, 0}, {0, 1},
			},
		},
		{
			"ArtistsDuration",
			charts.ArtistsDuration([][]charts.Song{
				{charts.Song{Artist: "A", Duration: 1},
					charts.Song{Artist: "B", Duration: 2},
				},
				{charts.Song{Artist: "C", Duration: 1},
					charts.Song{Artist: "A", Duration: 1},
				},
			}),
			[]charts.Title{charts.ArtistTitle("A"), charts.ArtistTitle("B"), charts.ArtistTitle("C")},
			[][]float64{
				{1, 1}, {2, 0}, {0, 1},
			},
		},
		{
			"Songs",
			charts.Songs([][]charts.Song{
				{
					charts.Song{Artist: "A", Title: "b", Duration: 1},
					charts.Song{Artist: "B", Title: "b", Duration: 2},
					charts.Song{Artist: "A", Title: "a", Duration: 1},
				},
				{
					charts.Song{Artist: "C", Title: "b", Duration: 1},
					charts.Song{Artist: "A", Title: "b", Duration: 1},
				},
			}),
			[]charts.Title{
				charts.SongTitle(charts.Song{Artist: "A", Title: "a"}), charts.SongTitle(charts.Song{Artist: "A", Title: "b"}),
				charts.SongTitle(charts.Song{Artist: "B", Title: "b"}),
				charts.SongTitle(charts.Song{Artist: "C", Title: "b"}),
			},
			[][]float64{
				{1, 0}, {1, 1},
				{1, 0},
				{0, 1},
			},
		},
		{
			"SongsDuration",
			charts.SongsDuration([][]charts.Song{
				{
					charts.Song{Artist: "A", Title: "b", Duration: 1},
					charts.Song{Artist: "B", Title: "b", Duration: 2},
					charts.Song{Artist: "A", Title: "a", Duration: 1},
				},
				{
					charts.Song{Artist: "C", Title: "b", Duration: 1},
					charts.Song{Artist: "A", Title: "b", Duration: 1},
				},
			}),
			[]charts.Title{
				charts.SongTitle(charts.Song{Artist: "A", Title: "a"}), charts.SongTitle(charts.Song{Artist: "A", Title: "b"}),
				charts.SongTitle(charts.Song{Artist: "B", Title: "b"}),
				charts.SongTitle(charts.Song{Artist: "C", Title: "b"}),
			},
			[][]float64{
				{1, 0}, {1, 1},
				{2, 0},
				{0, 1},
			},
		},
		{
			"single column normalizer",
			charts.NormalizeColumn(charts.Artists([][]charts.Song{
				{
					charts.Song{Artist: "A"}, charts.Song{Artist: "A"},
					charts.Song{Artist: "B"},
					charts.Song{Artist: "C"},
				},
				{
					charts.Song{Artist: "B"}, charts.Song{Artist: "B"}, charts.Song{Artist: "B"},
					charts.Song{Artist: "C"}, charts.Song{Artist: "C"}, charts.Song{Artist: "C"},
				},
				{},
			})),
			[]charts.Title{charts.ArtistTitle("A"), charts.ArtistTitle("B"), charts.ArtistTitle("C")},
			[][]float64{
				{.5, 0, 0}, {.25, .5, 0}, {.25, .5, 0},
			},
		},
		{
			"charts.FromMap",
			charts.FromMap(map[string][]float64{
				"A": {1, 1},
				"B": {1, 0},
				"C": {0, 1},
			}),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B"), charts.KeyTitle("C")},
			[][]float64{
				{1, 1}, {1, 0}, {0, 1},
			},
		},
		{
			"Only",
			charts.Only(charts.FromMap(map[string][]float64{
				"A": {1, 1},
				"B": {1, 0},
				"C": {0, 1},
			}), []charts.Title{charts.KeyTitle("A"), charts.KeyTitle("C")}),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("C")},
			[][]float64{
				{1, 1}, {0, 1},
			},
		},
		{
			"Intervals with Sum",
			charts.Intervals(charts.FromMap(map[string][]float64{
				"A": {1, 1, 0, 1, 3, 3, 2, 0},
				"B": {1, 0, 1, 0, 0, 0, 0, 5},
				"C": {0, 1, 0, 9, 0, 2, 0, 0},
			}), charts.Ranges{
				Delims: []rsrc.Day{
					rsrc.ParseDay("2022-01-01"),
					rsrc.ParseDay("2022-01-03"),
					rsrc.ParseDay("2022-01-05"),
					rsrc.ParseDay("2022-01-08")},
				Registered: rsrc.ParseDay("2022-01-01"),
			}, charts.Sum),
			[]charts.Title{charts.KeyTitle("A"), charts.KeyTitle("B"), charts.KeyTitle("C")},
			[][]float64{
				{2, 1, 8}, {1, 1, 0}, {1, 9, 2},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if !areTitlesSame(c.titles, c.charts.Titles()) {
				t.Fatalf("titles are not equal: %v != %v",
					c.titles, c.charts.Titles())
			}

			data := c.charts.Data(c.titles, 0, c.charts.Len())
			for i, title := range c.titles {
				row := c.charts.Data([]charts.Title{title}, 0, c.charts.Len())[0]
				if !reflect.DeepEqual(c.lines[i], row) {
					t.Errorf("row, '%v': %v != %v", title, c.lines[i], row)
				}

				if !reflect.DeepEqual(c.lines[i], data[i]) {
					t.Errorf("data, '%v': %v != %v", title, c.lines[i], data[i])
				}
			}

			for i := 0; i < c.charts.Len(); i++ {
				col := c.charts.Column(c.titles, i)
				for j, title := range c.titles {
					if c.lines[j][i] != col[j] {
						t.Errorf("col %v, %v: %v != %v",
							title, i,
							c.lines[j][i], col[j])
					}
				}
			}
		})
	}
}
