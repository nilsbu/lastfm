package charts

// Tag contains information about a tag.
type Tag struct {
	Name   string
	Total  int64
	Reach  int64
	Weight int
}

type Partition interface {
	Partitions() []Key
	Get(key Key) (partition Key)
}

type mapPart struct {
	assoc      map[string]Key
	partitions []Key
}

func (p mapPart) Partitions() []Key {
	return p.partitions
}

func (p mapPart) Get(key Key) Key {
	if partition, ok := p.assoc[key.String()]; ok {
		return partition
	}

	return simpleKey("-")
}

// artistPartition is a Partition that takes uses the Key.ArtistName() to
// categorize.
type artistPartition struct {
	mapPart
}

func (p artistPartition) Get(key Key) Key {
	if partition, ok := p.assoc[key.ArtistName()]; ok {
		return partition
	}

	return simpleKey("-")
}

func Supertags(
	tags map[string][]Tag,
	supertags map[string]string,
	corrections map[string]string,
) Partition {
	partition := artistPartition{mapPart{
		assoc:      make(map[string]Key),
		partitions: []Key{},
	}}

	// compile partitions
	names := map[string]bool{}
	for _, supertag := range supertags {
		names[supertag] = true
	}
	names["-"] = true

	for name, _ := range names {
		partition.partitions = append(partition.partitions, tagKey(name))
	}

	// compile association
	for name, values := range tags {
		supertag := tagKey("-")
		for _, tag := range values {
			if stag, ok := supertags[tag.Name]; ok {
				supertag = tagKey(stag)
				break
			}
		}

		partition.assoc[name] = supertag
	}

	for name, correction := range corrections {
		partition.assoc[name] = tagKey(correction)
	}

	return partition
}

func (c Charts) Group(partitions Partition) (tagcharts Charts) {
	size := c.Len()

	indices := map[string]int{}
	values := [][]float64{}
	for i, name := range partitions.Partitions() {
		indices[name.String()] = i
		values = append(values, make([]float64, size))
	}

	for i, name := range c.Keys {
		lineID := indices[partitions.Get(name).String()]
		line := values[lineID]
		for j := range line {
			line[j] += c.Values[i][j]
		}
	}

	keys := []Key{}
	for _, key := range partitions.Partitions() {
		keys = append(keys, key)
	}

	return Charts{
		Headers: c.Headers,
		Keys:    keys,
		Values:  values,
	}
}

func (c Charts) Split(partitions Partition) map[string]Charts {
	buckets := map[string]Charts{}

	for _, supertag := range partitions.Partitions() {
		buckets[supertag.String()] = Charts{
			Headers: c.Headers,
		}
	}

	for i, key := range c.Keys {
		p := partitions.Get(key)
		keys := append(buckets[p.String()].Keys, key)
		values := append(buckets[p.String()].Values, c.Values[i])
		buckets[p.String()] = Charts{buckets[p.String()].Headers, keys, values}
	}

	return buckets
}
