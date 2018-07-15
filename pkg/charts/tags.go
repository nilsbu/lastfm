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

	return ""
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
	names[""] = true

	for name, _ := range names {
		partition.partitions = append(partition.partitions, name)
	}

	// compile association
	for name, values := range tags {
		var supertag string
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
	var size int
	for _, line := range c {
		size = len(line)
		break
	}

	tagcharts = make(Charts)
	for _, supertag := range partitions.Partitions() {
		tagcharts[supertag] = make([]float64, size)
	}

	for name, values := range c {
		line := tagcharts[partitions.Get(name)]
		for i := range line {
			line[i] += values[i]
		}
	}

	return tagcharts
}

func (c Charts) Split(partitions Partition) map[string]Charts {
	buckets := map[string]Charts{}

	for _, supertag := range partitions.Partitions() {
		buckets[supertag] = Charts{}
	}

	for name, values := range c {
		buckets[partitions.Get(name)][name] = values
	}

	return buckets
}
