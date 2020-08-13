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
	} {
		t.Run(c.name, func(t *testing.T) {
			if !areTitlesSame(c.titles, c.charts.Titles()) {
				t.Fatalf("titles are not equal: %v != %v",
					c.titles, c.charts.Titles())
			}

			for i, title := range c.titles {
				row := c.charts.Row(title, 0, c.charts.Len())
				if !reflect.DeepEqual(c.lines[i], row) {
					t.Errorf("'%v': %v != %v", title, c.lines[i], row)
				}
			}
		})
	}
}
