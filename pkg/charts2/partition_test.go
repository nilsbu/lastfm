package charts2

import (
	"testing"
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
	} {
		t.Run(c.name, func(t *testing.T) {
			for i, tp := range c.titlePartitions {
				partition := c.partition.Partition(tp.title)
				if tp.partition.Key() != partition.Key() {
					t.Errorf("%v: '%v' != '%v'", i, tp.partition, partition)
				}
			}

			for i, pt := range c.partitionTitles {
				titles := c.partition.Titles(pt.partition)
				if len(titles) != len(pt.titles) {
					t.Fatal(len(titles), "!=", len(pt.titles))
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
