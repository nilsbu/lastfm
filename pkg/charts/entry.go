package charts

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type EntryDate struct {
	Name string
	Date rsrc.Day
}

func (c Charts) FindEntryDates(threshold float64) (entryDates []EntryDate) {

	for i, name := range c.Keys {
		if c.Values[i][c.Len()-1] < threshold {
			continue
		}

		for j, value := range c.Values[i] {
			if value >= threshold {
				date := c.Headers.At(j).Begin
				entryDates = append(entryDates, EntryDate{name.String(), date})
				break
			}
		}
	}

	return
}

func (c Charts) FindEntryDatesDynamic(threshold float64,
) (entryDates []EntryDate) {

	nm := GaussianNormalizer{Sigma: 30}
	nc := nm.Normalize(c)

	nsum := nc.Sum()

	for i, name := range nsum.Keys {
		values := nsum.Values[i]
		if values[len(values)-1] < threshold {
			continue
		}

		idx := -1
		var maxv float64
		for j, value := range values {
			if value > 2*threshold {
				break
			}

			if maxv < nc.Values[i][j] {
				maxv = nc.Values[i][j]
				idx = j
			}
		}

		if idx != -1 {
			date := c.Headers.At(idx).Begin
			entryDates = append(entryDates, EntryDate{name.String(), date})
		}
	}

	return
}

func (c Charts) GetYearPartition(threshold float64) Partition {
	entryDates := c.FindEntryDatesDynamic(threshold)

	p := mapPart{
		assoc:      make(map[string]Key),
		partitions: []Key{},
	}

	for _, entryDate := range entryDates {
		p.assoc[entryDate.Name] = simpleKey(entryDate.Date.Time().Format("2006"))
	}

	ii := Years(
		c.Headers.At(0).Begin,
		c.Headers.At(c.Len()-1).Before,
		1)

	for i := 0; i < ii.Len(); i++ {
		p.partitions = append(
			p.partitions,
			simpleKey(ii.At(i).Begin.Time().Format("2006")))
	}
	p.partitions = append(p.partitions, simpleKey("-"))
	return p
}
