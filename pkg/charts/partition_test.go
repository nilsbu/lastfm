package charts_test

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type titlePartition struct {
	title, partition charts.Title
}

type partitionTitles struct {
	partition charts.Title
	titles    []charts.Title
}

func TestPartiton(t *testing.T) {
	for _, c := range []struct {
		name            string
		partition       charts.Partition
		titlePartitions []titlePartition
		partitionTitles []partitionTitles
		partitions      []charts.Title
	}{
		{
			"empty key partition",
			charts.KeyPartition([][2]charts.Title{}),
			[]titlePartition{
				{charts.KeyTitle("a"), charts.KeyTitle("")},
			},
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
			[]titlePartition{
				{charts.KeyTitle("a"), charts.KeyTitle("l")},
				{charts.KeyTitle("A"), charts.KeyTitle("u")},
				{charts.KeyTitle("b"), charts.KeyTitle("l")},
				{charts.KeyTitle("C"), charts.KeyTitle("u")},
				{charts.KeyTitle("B"), charts.KeyTitle("")},
			},
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
			[]titlePartition{
				{charts.KeyTitle("a"), charts.StringTitle("total")},
				{charts.KeyTitle("b"), charts.StringTitle("total")},
				{charts.KeyTitle("B"), charts.KeyTitle("")},
			},
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
			[]titlePartition{
				{charts.KeyTitle("a"), charts.KeyTitle("")},
			},
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
			[]titlePartition{
				{charts.ArtistTitle("A"), charts.KeyTitle("vowel")},
				{charts.ArtistTitle("B"), charts.KeyTitle("consonant")},
				{charts.ArtistTitle("C"), charts.KeyTitle("consonant")},
				{charts.ArtistTitle("X"), charts.KeyTitle("")},
			},
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
			[]titlePartition{
				{charts.ArtistTitle("A"), charts.KeyTitle("vowel")},
				{charts.ArtistTitle("Y"), charts.KeyTitle("vowel")},
				{charts.ArtistTitle("Ü"), charts.KeyTitle("umlaut")},
				{charts.ArtistTitle("X"), charts.KeyTitle("")},
			},
			[]partitionTitles{
				{charts.KeyTitle("consonant"), []charts.Title{}},
				{charts.KeyTitle("vowel"), []charts.Title{charts.ArtistTitle("A"), charts.ArtistTitle("Y")}},
				{charts.KeyTitle("umlaut"), []charts.Title{charts.ArtistTitle("Ü")}},
			},
			[]charts.Title{charts.KeyTitle("vowel"), charts.KeyTitle("umlaut")},
		},
		{
			"year partition with no eligible artists",
			charts.YearPartition(
				charts.FromMap(map[string][]float64{"not": {0, 1}}),
				charts.FromMap(map[string][]float64{"not": {0, 1}}),
				rsrc.ParseDay("2019-12-31"),
			),
			[]titlePartition{
				{charts.ArtistTitle("not"), charts.KeyTitle("")},
			},
			[]partitionTitles{
				{charts.KeyTitle("2019"), []charts.Title{}},
				{charts.KeyTitle("2020"), []charts.Title{}},
			},
			[]charts.Title{charts.KeyTitle("2019"), charts.KeyTitle("2020")},
		},
		{
			"year partition with values",
			charts.YearPartition(
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
			),
			[]titlePartition{
				{charts.ArtistTitle("not"), charts.KeyTitle("")},
				{charts.ArtistTitle("first"), charts.KeyTitle("2019")},
				{charts.ArtistTitle("first2"), charts.KeyTitle("2019")},
				{charts.ArtistTitle("last"), charts.KeyTitle("2020")},
				{charts.ArtistTitle("last2"), charts.KeyTitle("2020")},
			},
			[]partitionTitles{
				{charts.KeyTitle("2019"), []charts.Title{charts.KeyTitle("first"), charts.KeyTitle("first2")}},
				{charts.KeyTitle("2020"), []charts.Title{charts.KeyTitle("last"), charts.KeyTitle("last2")}},
			},
			[]charts.Title{charts.KeyTitle("2019"), charts.KeyTitle("2020")},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			for _, tp := range c.titlePartitions {
				partition := c.partition.Partition(tp.title)
				if tp.partition.Key() != partition.Key() {
					t.Errorf("'%v': '%v' != '%v'", tp.title, tp.partition, partition)
				}
			}

			for _, pt := range c.partitionTitles {
				titles := c.partition.Titles(pt.partition)
				if len(titles) != len(pt.titles) {
					t.Fatalf("for partition '%v': %v != %v",
						pt.partition, len(titles), len(pt.titles))
				}
				if !areTitlesSame(pt.titles, titles) {
					t.Errorf("for partition '%v', titles unequal: %v != %v", pt.partition, pt.titles, titles)
				}
			}

			partitions := c.partition.Partitions()
			if !areTitlesSame(c.partitions, partitions) {
				t.Errorf("partitions unequal: %v != %v", c.partitions, partitions)
			}
		})
	}
}
