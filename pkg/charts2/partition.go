package charts2

import (
	legacy "github.com/nilsbu/lastfm/pkg/charts"
)

// Partition divides titles into separate groups called partitions.
// Each title belongs to no more than one partition which can be obtained by
// Partition(title). Conversely Titles(partition) will include all titles
// belonging to a given partition.
// Partitions() returns the set of partitions.
// TODO Should key errors be handled explicitely?
type Partition interface {
	Partition(title Title) (partition Title)
	Titles(partition Title) []Title
	Partitions() []Title
}

type biMapPartition struct {
	titlePartition  map[string]Title
	partitionTitles map[string][]Title
	partitions      []Title
	key             func(title Title) string
}

func (p biMapPartition) Partition(title Title) Title {
	if partition, ok := p.titlePartition[p.key(title)]; ok {
		return partition
	}
	return KeyTitle("")
}
func (p biMapPartition) Titles(partition Title) []Title {
	if titles, ok := p.partitionTitles[p.key(partition)]; ok {
		return titles
	}
	return []Title{}
}

func (p biMapPartition) Partitions() []Title {
	return p.partitions
}

// KeyPartition returns a Partition where Title.Key() determines a Title's
// membership in a partition. Membership is constructed by tpPairs. The first
// element in each array denotes the title, the second the partition.
func KeyPartition(tpPairs [][2]Title) Partition {
	f := func(t Title) string { return t.Key() }
	return biMapPartition{
		titlePartition:  buildTPMap(tpPairs, f),
		partitionTitles: buildPTMap(tpPairs, f),
		partitions:      buildPartitions(tpPairs, f),
		key:             f,
	}
}

func buildTPMap(tpPairs [][2]Title, key func(Title) string) map[string]Title {
	tpMap := make(map[string]Title)
	for _, tp := range tpPairs {
		tpMap[key(tp[0])] = tp[1]
	}
	return tpMap
}

func buildPTMap(tpPairs [][2]Title, key func(Title) string) map[string][]Title {
	ptMap := map[string][]Title{}
	for _, tp := range tpPairs {
		k := key(tp[1])
		if _, ok := ptMap[k]; !ok {
			ptMap[k] = []Title{tp[0]}
		} else {
			ptMap[k] = append(ptMap[k], tp[0])
		}
	}
	return ptMap
}

func buildPartitions(tpPairs [][2]Title, key func(Title) string) []Title {
	set := make(map[string]Title)
	for _, tp := range tpPairs {
		set[key(tp[1])] = tp[1]
	}

	partitions := make([]Title, len(set))
	i := 0
	for _, p := range set {
		partitions[i] = p
		i++
	}
	return partitions
}

// TotalPartition is a Partition where all titles are mapped to the same
// partition.
func TotalPartition(titles []Title) Partition {
	p := StringTitle("total")
	tpPairs := make([][2]Title, len(titles))
	for i, title := range titles {
		tpPairs[i] = [2]Title{title, p}
	}
	return KeyPartition(tpPairs)
}

// FirstTagPartition creates a partition where a select group of tags point to
// the partitions. Each key is assigned to its partition by first tag included
// in tagToPartition. Corrections can override this.
func FirstTagPartition(
	tags map[string][]legacy.Tag, // TODO: use different type or move struct here
	tagToPartition map[string]string,
	corrections map[string]string,
) Partition {
	tpPairs := [][2]Title{}

	for title, list := range tags {
		for _, tag := range list {
			if partition, ok := corrections[tag.Name]; ok {
				tpPairs = append(tpPairs, [2]Title{ArtistTitle(title), KeyTitle(partition)})
				break
			}
			if partition, ok := tagToPartition[tag.Name]; ok {
				tpPairs = append(tpPairs, [2]Title{ArtistTitle(title), KeyTitle(partition)})
				break
			}
		}
	}

	return KeyPartition(tpPairs)
}
