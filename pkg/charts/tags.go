package charts

// Tag contains information about a tag.
type Tag struct {
	Name   string
	Total  int64
	Reach  int64
	Weight int
}

type Partition interface {
	Partitions() []string
	Get(key string) (partition string)
}

type mapPart struct {
	assoc      map[string]string
	partitions []string
}

func (p mapPart) Partitions() []string {
	return p.partitions
}

func (p mapPart) Get(key string) string {
	if partition, ok := p.assoc[key]; ok {
		return partition
	}

	return "-"
}

func Supertags(
	tags map[string][]Tag,
	supertags map[string]string,
	corrections map[string]string,
) Partition {
	partition := mapPart{
		assoc:      make(map[string]string),
		partitions: []string{},
	}

	// compile partitions
	names := map[string]bool{}
	for _, supertag := range supertags {
		names[supertag] = true
	}
	names["-"] = true

	for name, _ := range names {
		partition.partitions = append(partition.partitions, name)
	}

	// compile association
	for name, values := range tags {
		supertag := "-"
		for _, tag := range values {
			if stag, ok := supertags[tag.Name]; ok {
				supertag = stag
				break
			}
		}

		partition.assoc[name] = supertag
	}

	for name, correction := range corrections {
		partition.assoc[name] = correction
	}

	return partition
}

func (c Charts) Group(partitions Partition) (tagcharts Charts) {
	size := c.Len()

	indices := map[string]int{}
	values := [][]float64{}
	for i, name := range partitions.Partitions() {
		indices[name] = i
		values = append(values, make([]float64, size))
	}

	for i, name := range c.Keys {
		lineID := indices[partitions.Get(name.String())]
		line := values[lineID]
		for j := range line {
			line[j] += c.Values[i][j]
		}
	}

	keys := []Key{}
	for _, key := range partitions.Partitions() {
		keys = append(keys, simpleKey(key))
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
		buckets[supertag] = Charts{
			Headers: c.Headers,
		}
	}

	for i, key := range c.Keys {
		p := partitions.Get(key.String())
		keys := append(buckets[p].Keys, key)
		values := append(buckets[p].Values, c.Values[i])
		buckets[p] = Charts{buckets[p].Headers, keys, values}
	}

	return buckets
}
