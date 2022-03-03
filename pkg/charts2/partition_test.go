package charts2

import (
	"testing"

	legacy "github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type titlePartition struct {
	title, partition Title
}

type partitionTitles struct {
	partition Title
	titles    []Title
}

func TestPartiton(t *testing.T) {
	for _, c := range []struct {
		name            string
		partition       Partition
		titlePartitions []titlePartition
		partitionTitles []partitionTitles
		partitions      []Title
	}{
		{
			"empty key partition",
			KeyPartition([][2]Title{}),
			[]titlePartition{
				{KeyTitle("a"), KeyTitle("")},
			},
			[]partitionTitles{
				{KeyTitle("x"), []Title{}},
			},
			[]Title{},
		},
		{
			"key partition",
			KeyPartition([][2]Title{
				{KeyTitle("a"), KeyTitle("l")},
				{KeyTitle("A"), KeyTitle("u")},
				{KeyTitle("b"), KeyTitle("l")},
				{KeyTitle("C"), KeyTitle("u")},
			}),
			[]titlePartition{
				{KeyTitle("a"), KeyTitle("l")},
				{KeyTitle("A"), KeyTitle("u")},
				{KeyTitle("b"), KeyTitle("l")},
				{KeyTitle("C"), KeyTitle("u")},
				{KeyTitle("B"), KeyTitle("")},
			},
			[]partitionTitles{
				{KeyTitle("l"), []Title{KeyTitle("a"), KeyTitle("b")}},
				{KeyTitle("u"), []Title{KeyTitle("A"), KeyTitle("C")}},
				{KeyTitle("x"), []Title{}},
			},
			[]Title{KeyTitle("l"), KeyTitle("u")},
		},
		{
			"totalPartition",
			TotalPartition([]Title{KeyTitle("a"), KeyTitle("b")}),
			[]titlePartition{
				{KeyTitle("a"), StringTitle("total")},
				{KeyTitle("b"), StringTitle("total")},
				{KeyTitle("B"), KeyTitle("")},
			},
			[]partitionTitles{
				{StringTitle("total"), []Title{KeyTitle("a"), KeyTitle("b")}},
				{KeyTitle("n"), []Title{}},
			},
			[]Title{StringTitle("total")},
		},
		{
			"empty first key partition",
			FirstTagPartition(
				map[string][]legacy.Tag{},
				map[string]string{},
				nil,
			),
			[]titlePartition{
				{KeyTitle("a"), KeyTitle("")},
			},
			[]partitionTitles{
				{KeyTitle("x"), []Title{}},
			},
			[]Title{},
		},
		{
			"first key partition without correction",
			FirstTagPartition(
				map[string][]legacy.Tag{
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
				{ArtistTitle("A"), KeyTitle("vowel")},
				{ArtistTitle("B"), KeyTitle("consonant")},
				{ArtistTitle("C"), KeyTitle("consonant")},
				{ArtistTitle("X"), KeyTitle("")},
			},
			[]partitionTitles{
				{KeyTitle("consonant"), []Title{ArtistTitle("B"), ArtistTitle("C")}},
				{KeyTitle("vowel"), []Title{ArtistTitle("A")}},
			},
			[]Title{KeyTitle("vowel"), KeyTitle("consonant")},
		},
		{
			"first key partition with correction",
			FirstTagPartition(
				map[string][]legacy.Tag{
					"A": {{Name: "a", Weight: 100}, {Name: "c", Weight: 25}},
					"Y": {{Name: "b", Weight: 25}, {Name: "y", Weight: 100}}, // Ignore Weight
					"Ü": {{Name: "-", Weight: 100}, {Name: "ü", Weight: 50}},
				},
				map[string]string{
					"a": "vowel", "y": "consonant", "ü": "vowel",
				},
				map[string]string{
					"y": "vowel", "ü": "umlaut",
				},
			),
			[]titlePartition{
				{ArtistTitle("A"), KeyTitle("vowel")},
				{ArtistTitle("Y"), KeyTitle("vowel")},
				{ArtistTitle("Ü"), KeyTitle("umlaut")},
				{ArtistTitle("X"), KeyTitle("")},
			},
			[]partitionTitles{
				{KeyTitle("consonant"), []Title{}},
				{KeyTitle("vowel"), []Title{ArtistTitle("A"), ArtistTitle("Y")}},
				{KeyTitle("umlaut"), []Title{ArtistTitle("Ü")}},
			},
			[]Title{KeyTitle("vowel"), KeyTitle("umlaut")},
		},
		{
			"year partition with no eligible artists",
			YearPartition(
				chartsFromMap(map[string][]float64{"not": {0, 1}}),
				chartsFromMap(map[string][]float64{"not": {0, 1}}),
				rsrc.ParseDay("2019-12-31"),
			),
			[]titlePartition{
				{ArtistTitle("not"), KeyTitle("")},
			},
			[]partitionTitles{
				{KeyTitle("2019"), []Title{}},
				{KeyTitle("2020"), []Title{}},
			},
			[]Title{KeyTitle("2019"), KeyTitle("2020")},
		},
		{
			"year partition with values",
			YearPartition(
				chartsFromMap(map[string][]float64{
					"not":    {0, 0, 1, 0},
					"first":  {0, 4, 10, 0}, // higher value irrelevant since 4 is reached in 2019
					"first2": {0, 2, 1, 0},
					"last":   {0, 2, 1, 0},
					"last2":  {0, 1, 2, 0},
				}),
				chartsFromMap(map[string][]float64{
					"not":    {0, 0, 1, 1},
					"first":  {0, 4, 4, 4},
					"first2": {0, 3, 4, 4},
					"last":   {0, 1, 3, 3},
					"last2":  {0, 2, 3, 3},
				}),
				rsrc.ParseDay("2019-12-30"),
			),
			[]titlePartition{
				{ArtistTitle("not"), KeyTitle("")},
				{ArtistTitle("first"), KeyTitle("2019")},
				{ArtistTitle("first2"), KeyTitle("2019")},
				{ArtistTitle("last"), KeyTitle("2020")},
				{ArtistTitle("last2"), KeyTitle("2020")},
			},
			[]partitionTitles{
				{KeyTitle("2019"), []Title{KeyTitle("first"), KeyTitle("first2")}},
				{KeyTitle("2020"), []Title{KeyTitle("last"), KeyTitle("last2")}},
			},
			[]Title{KeyTitle("2019"), KeyTitle("2020")},
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
