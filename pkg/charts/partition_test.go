package charts_test

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type partitionTitles struct {
	partition charts.Title
	titles    []charts.Title
}

func TestPartiton(t *testing.T) {
	for _, c := range []struct {
		name            string
		partition       charts.Partition
		partitionTitles []partitionTitles
		partitions      []charts.Title
	}{
		// TODO check errors
		{
			"empty key partition",
			charts.KeyPartition([][2]charts.Title{}),
			[]partitionTitles{
				{charts.KeyTitle("x"), []charts.Title{}},
			},
			[]charts.Title{},
		},
		{
			"key partition",
			charts.KeyPartition([][2]charts.Title{
				{charts.KeyTitle("a"), charts.KeyTitle("l")},
				{charts.KeyTitle("A"), charts.KeyTitle("u")},
				{charts.KeyTitle("b"), charts.KeyTitle("l")},
				{charts.KeyTitle("C"), charts.KeyTitle("u")},
			}),
			[]partitionTitles{
				{charts.KeyTitle("l"), []charts.Title{charts.KeyTitle("a"), charts.KeyTitle("b")}},
				{charts.KeyTitle("u"), []charts.Title{charts.KeyTitle("A"), charts.KeyTitle("C")}},
				{charts.KeyTitle("x"), []charts.Title{}},
			},
			[]charts.Title{charts.KeyTitle("l"), charts.KeyTitle("u")},
		},
		{
			"totalPartition",
			charts.TotalPartition([]charts.Title{charts.KeyTitle("a"), charts.KeyTitle("b")}),
			[]partitionTitles{
				{charts.StringTitle("total"), []charts.Title{charts.KeyTitle("a"), charts.KeyTitle("b")}},
				{charts.KeyTitle("n"), []charts.Title{}},
			},
			[]charts.Title{charts.StringTitle("total")},
		},
		{
			"empty first key partition",
			charts.FirstTagPartition(
				map[string][]info.Tag{},
				map[string]string{},
				nil,
			),
			[]partitionTitles{
				{charts.KeyTitle("x"), []charts.Title{}},
			},
			[]charts.Title{},
		},
		{
			"first key partition without correction",
			charts.FirstTagPartition(
				map[string][]info.Tag{
					"A": {{Name: "a", Weight: 100}, {Name: "c", Weight: 25}},
					"B": {{Name: "b", Weight: 25}, {Name: "c", Weight: 100}}, // Ignore Weight
					"C": {{Name: "-", Weight: 100}, {Name: "c", Weight: 50}},
				},
				map[string]string{
					"a": "vowel", "b": "consonant", "c": "consonant",
				},
				nil,
			),
			[]partitionTitles{
				{charts.KeyTitle("consonant"), []charts.Title{charts.ArtistTitle("B"), charts.ArtistTitle("C")}},
				{charts.KeyTitle("vowel"), []charts.Title{charts.ArtistTitle("A")}},
			},
			[]charts.Title{charts.KeyTitle("vowel"), charts.KeyTitle("consonant")},
		},
		{
			"first key partition with correction",
			charts.FirstTagPartition(
				map[string][]info.Tag{
					"A": {{Name: "a", Weight: 100}, {Name: "c", Weight: 25}},
					"Y": {{Name: "b", Weight: 25}, {Name: "y", Weight: 100}}, // Ignore Weight
					"Ü": {{Name: "-", Weight: 100}, {Name: "ü", Weight: 50}},
				},
				map[string]string{
					"a": "vowel", "y": "consonant", "ü": "vowel",
				},
				map[string]string{
					"Y": "vowel", "Ü": "umlaut",
				},
			),
			[]partitionTitles{
				{charts.KeyTitle("consonant"), []charts.Title{}},
				{charts.KeyTitle("vowel"), []charts.Title{charts.ArtistTitle("A"), charts.ArtistTitle("Y")}},
				{charts.KeyTitle("umlaut"), []charts.Title{charts.ArtistTitle("Ü")}},
			},
			[]charts.Title{charts.KeyTitle("vowel"), charts.KeyTitle("umlaut")},
		},
		{
			"year partition with no eligible artists",
			func() charts.Partition {
				p, _ := charts.YearPartition(
					charts.FromMap(map[string][]float64{"not": {0, 1}}),
					charts.FromMap(map[string][]float64{"not": {0, 1}}),
					rsrc.ParseDay("2019-12-31"),
				)
				return p
			}(),
			[]partitionTitles{
				{charts.KeyTitle("2019"), []charts.Title{}},
				{charts.KeyTitle("2020"), []charts.Title{}},
			},
			[]charts.Title{charts.KeyTitle("2019"), charts.KeyTitle("2020")},
		},
		{
			"year partition with values",
			func() charts.Partition {
				p, _ := charts.YearPartition(
					charts.FromMap(map[string][]float64{
						"not":    {0, 0, 1, 0},
						"first":  {0, 4, 10, 0}, // higher value irrelevant since 4 is reached in 2019
						"first2": {0, 2, 1, 0},
						"last":   {0, 2, 1, 0},
						"last2":  {0, 1, 2, 0},
					}),
					charts.FromMap(map[string][]float64{
						"not":    {0, 0, 1, 1},
						"first":  {0, 4, 4, 4},
						"first2": {0, 3, 4, 4},
						"last":   {0, 1, 3, 3},
						"last2":  {0, 2, 3, 3},
					}),
					rsrc.ParseDay("2019-12-30"),
				)
				return p
			}(),
			[]partitionTitles{
				{charts.KeyTitle("2019"), []charts.Title{charts.KeyTitle("first"), charts.KeyTitle("first2")}},
				{charts.KeyTitle("2020"), []charts.Title{charts.KeyTitle("last"), charts.KeyTitle("last2")}},
			},
			[]charts.Title{charts.KeyTitle("2019"), charts.KeyTitle("2020")},
		},
		{
			"tag weight partition",
			charts.TagWeightPartition(
				[]charts.Title{
					charts.ArtistTitle("first"),
					charts.ArtistTitle("second"),
					charts.ArtistTitle("firstAgain"),
					charts.ArtistTitle("none"),
				},
				[][]info.Tag{
					{info.Tag{Name: "1st", Reach: 100000}, info.Tag{Name: "2nd", Reach: 100000}},
					{info.Tag{Name: "1stnope", Reach: 100}, info.Tag{Name: "2nd", Reach: 100000}},
					{info.Tag{Name: "1stnope", Reach: 100}, info.Tag{Name: "1st", Reach: 100000}, info.Tag{Name: "2nd", Reach: 100000}},
					{info.Tag{Name: "1stnope", Reach: 100}, info.Tag{Name: "blacklisted", Reach: 100000}},
				},
				map[string]interface{}{"blacklisted": nil}),
			[]partitionTitles{
				{charts.KeyTitle("1st"), []charts.Title{charts.ArtistTitle("first"), charts.ArtistTitle("firstAgain")}},
				{charts.KeyTitle("2nd"), []charts.Title{charts.ArtistTitle("second")}},
				{charts.KeyTitle("-"), []charts.Title{charts.ArtistTitle("none")}},
			},
			[]charts.Title{charts.KeyTitle("1st"), charts.KeyTitle("2nd"), charts.KeyTitle("-")},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			for _, pt := range c.partitionTitles {
				titles, _ := c.partition.Titles(pt.partition)
				if len(titles) != len(pt.titles) {
					t.Fatalf("for partition '%v': %v != %v",
						pt.partition, len(titles), len(pt.titles))
				}
				if !areTitlesSame(pt.titles, titles) {
					t.Errorf("for partition '%v', titles unequal: %v != %v", pt.partition, pt.titles, titles)
				}
			}

			partitions, _ := c.partition.Partitions()
			if !areTitlesSame(c.partitions, partitions) {
				t.Errorf("partitions unequal: %v != %v", c.partitions, partitions)
			}
		})
	}
}
