package charts

// Tag contains information about a tag.
type Tag struct {
	Name   string
	Total  int64
	Reach  int64
	Weight int
}

func (c Charts) Supertags(
	tags map[string][]Tag,
	supertags map[string]string,
) (tagcharts Charts) {
	var size int
	for _, line := range c {
		size = len(line)
		break
	}
	tagcharts = initSupertagCharts(supertags, size)

	for name, values := range c {
		var supertag string
		for _, tag := range tags[name] {
			if stag, ok := supertags[tag.Name]; ok {
				supertag = stag
				break
			}
		}

		line := tagcharts[supertag]
		for i := range line {
			line[i] += values[i]
		}
	}

	return tagcharts
}

func initSupertagCharts(supertags map[string]string, len int) Charts {
	charts := Charts{}

	for _, supertag := range supertags {
		charts[supertag] = make([]float64, len)
	}

	charts[""] = make([]float64, len)

	return charts
}

func (c Charts) SplitBySupertag(
	tags map[string][]Tag,
	supertags map[string]string,
) map[string]Charts {

	buckets := map[string]Charts{}

	for _, supertag := range supertags {
		buckets[supertag] = Charts{}
	}

	buckets[""] = Charts{}

	for name, values := range c {
		var supertag string
		for _, tag := range tags[name] {
			if stag, ok := supertags[tag.Name]; ok {
				supertag = stag
				break
			}
		}

		buckets[supertag][name] = values
	}

	return buckets
}
