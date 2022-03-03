package charts2

import (
	"testing"

	legacy "github.com/nilsbu/lastfm/pkg/charts"
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
					"A": []legacy.Tag{{Name: "a", Weight: 100}, {Name: "c", Weight: 25}},
					"B": []legacy.Tag{{Name: "b", Weight: 25}, {Name: "c", Weight: 100}}, // Ignore Weight
					"C": []legacy.Tag{{Name: "-", Weight: 100}, {Name: "c", Weight: 50}},
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
					"A": []legacy.Tag{{Name: "a", Weight: 100}, {Name: "c", Weight: 25}},
					"Y": []legacy.Tag{{Name: "b", Weight: 25}, {Name: "y", Weight: 100}}, // Ignore Weight
					"Ü": []legacy.Tag{{Name: "-", Weight: 100}, {Name: "ü", Weight: 50}},
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
	} {
		t.Run(c.name, func(t *testing.T) {
			for _, tp := range c.titlePartitions {
				partition := c.partition.Partition(tp.title)
				if tp.partition.Key() != partition.Key() {
					t.Errorf("'%v': '%v' != '%v'", tp.title, tp.partition, partition)
				}
			}

			for i, pt := range c.partitionTitles {
				titles := c.partition.Titles(pt.partition)
				if len(titles) != len(pt.titles) {
					t.Fatalf("for partition '%v': %v != %v",
						pt.partition, len(titles), len(pt.titles))
				}
				for j := range titles {
					if pt.titles[j].Key() != titles[j].Key() {
						t.Errorf("%v, %v: '%v' != '%v'", i, j, pt.titles[j], titles[j])
					}
				}
			}

			partitions := c.partition.Partitions()
			if !areTitlesSame(c.partitions, partitions) {
				t.Errorf("partitions unequal: %v != %v", c.partitions, partitions)
			}
		})
	}
}
