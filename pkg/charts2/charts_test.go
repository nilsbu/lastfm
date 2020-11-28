package charts2

import (
	"reflect"
	"testing"
)

func TestCharts(t *testing.T) {
	// The purpose here is to ensure the content of the charts is correct, the
	// API is tested in TestLazyCharts in greater detail.
	for _, c := range []struct {
		name   string
		charts LazyCharts
		titles []Title
		lines  [][]float64
	}{
		{
			"Artists",
			Artists([][]Song{
				{Song{Artist: "A", Duration: 1},
					Song{Artist: "B", Duration: 2},
				},
				{Song{Artist: "C", Duration: 1},
					Song{Artist: "A", Duration: 1},
				},
			}),
			[]Title{ArtistTitle("A"), ArtistTitle("B"), ArtistTitle("C")},
			[][]float64{
				{1, 1}, {1, 0}, {0, 1},
			},
		},
		{
			"ArtistsDuration",
			ArtistsDuration([][]Song{
				{Song{Artist: "A", Duration: 1},
					Song{Artist: "B", Duration: 2},
				},
				{Song{Artist: "C", Duration: 1},
					Song{Artist: "A", Duration: 1},
				},
			}),
			[]Title{ArtistTitle("A"), ArtistTitle("B"), ArtistTitle("C")},
			[][]float64{
				{1, 1}, {2, 0}, {0, 1},
			},
		},
		{
			"Songs",
			Songs([][]Song{
				{
					Song{Artist: "A", Title: "b", Duration: 1},
					Song{Artist: "B", Title: "b", Duration: 2},
					Song{Artist: "A", Title: "a", Duration: 1},
				},
				{
					Song{Artist: "C", Title: "b", Duration: 1},
					Song{Artist: "A", Title: "b", Duration: 1},
				},
			}),
			[]Title{
				SongTitle(Song{Artist: "A", Title: "a"}), SongTitle(Song{Artist: "A", Title: "b"}),
				SongTitle(Song{Artist: "B", Title: "b"}),
				SongTitle(Song{Artist: "C", Title: "b"}),
			},
			[][]float64{
				{1, 0}, {1, 1},
				{1, 0},
				{0, 1},
			},
		},
		{
			"SongsDuration",
			SongsDuration([][]Song{
				{
					Song{Artist: "A", Title: "b", Duration: 1},
					Song{Artist: "B", Title: "b", Duration: 2},
					Song{Artist: "A", Title: "a", Duration: 1},
				},
				{
					Song{Artist: "C", Title: "b", Duration: 1},
					Song{Artist: "A", Title: "b", Duration: 1},
				},
			}),
			[]Title{
				SongTitle(Song{Artist: "A", Title: "a"}), SongTitle(Song{Artist: "A", Title: "b"}),
				SongTitle(Song{Artist: "B", Title: "b"}),
				SongTitle(Song{Artist: "C", Title: "b"}),
			},
			[][]float64{
				{1, 0}, {1, 1},
				{2, 0},
				{0, 1},
			},
		},
		{
			"single column normalizer",
			NormalizeColumn(Artists([][]Song{
				{
					Song{Artist: "A"}, Song{Artist: "A"},
					Song{Artist: "B"},
					Song{Artist: "C"},
				},
				{
					Song{Artist: "B"}, Song{Artist: "B"}, Song{Artist: "B"},
					Song{Artist: "C"}, Song{Artist: "C"}, Song{Artist: "C"},
				},
				{},
			})),
			[]Title{ArtistTitle("A"), ArtistTitle("B"), ArtistTitle("C")},
			[][]float64{
				{.5, 0, 0}, {.25, .5, 0}, {.25, .5, 0},
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
				row := c.charts.Row(title, 0, c.charts.Len())
				if !reflect.DeepEqual(c.lines[i], row) {
					t.Errorf("row, '%v': %v != %v", title, c.lines[i], row)
				}

				if !reflect.DeepEqual(c.lines[i], data[title.Key()].Line) {
					t.Errorf("data, '%v': %v != %v", title, c.lines[i], data[title.Key()].Line)
				}
			}

			for i := 0; i < c.charts.Len(); i++ {
				col := c.charts.Column(c.titles, i)
				for j, title := range c.titles {
					if c.lines[j][i] != col[title.Key()].Value {
						t.Errorf("col %v, %v: %v != %v",
							title, i,
							c.lines[j][i], col[title.Key()].Value)
					}
				}
			}

		})
	}
}