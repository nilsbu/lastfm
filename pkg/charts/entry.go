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
				date := c.Headers.At(j)
				entryDates = append(entryDates, EntryDate{name.String(), date})
				break
			}
		}
	}

	return
}

func (c Charts) FindEntryDatesDynamic(threshold float64,
) (entryDates []EntryDate) {

	nm := GaussianNormalizer{
		Sigma:       30,
		MirrorFront: true,
		MirrorBack:  false}
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
			entryDates = append(entryDates, EntryDate{name.String(), c.Headers.At(idx)})
		}
	}

	return
}

func (c Charts) GetYearPartition(threshold float64) Partition {
	entryDates := c.FindEntryDatesDynamic(threshold)

	p := mapPart{
		assoc:      make(map[string]string),
		partitions: []string{},
	}

	for _, entryDate := range entryDates {
		p.assoc[entryDate.Name] = entryDate.Date.Time().Format("2006")
	}

	ii := newYearIterator(
		1,
		c.Headers.At(0),
		rsrc.ToDay(c.Headers.At(0).Midnight()+int64(86400*c.Len()))) // TODO

	for ii.HasNext() {
		p.partitions = append(p.partitions, ii.Next().Begin.Time().Format("2006"))
	}
	p.partitions = append(p.partitions, "-")
	return p
}
