package charts

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type EntryDate struct {
	Name string
	Date rsrc.Day
}

func (c Charts) FindEntryDates(registered rsrc.Day, threshold float64,
) (entryDates []EntryDate) {

	for name, values := range c {
		if values[len(values)-1] < threshold {
			continue
		}

		for i, value := range values {
			if value >= threshold {
				date := registered.Midnight() + int64(86400*i)
				entryDates = append(entryDates, EntryDate{name, rsrc.ToDay(date)})
				break
			}
		}
	}

	return
}

func (c Charts) FindEntryDatesDynamic(registered rsrc.Day, threshold float64,
) (entryDates []EntryDate) {

	nm := GaussianNormalizer{
		Sigma:       30,
		MirrorFront: true,
		MirrorBack:  false}
	nc := nm.Normalize(c)

	nsum := nc.Sum()

	for name, values := range nsum {
		if values[len(values)-1] < threshold {
			continue
		}

		idx := -1
		var maxv float64
		for i, value := range values {
			if value > 2*threshold {
				break
			}

			if maxv < nc[name][i] {
				maxv = nc[name][i]
				idx = i
			}
		}

		if idx != -1 {
			date := registered.Midnight() + int64(86400*idx)
			entryDates = append(entryDates, EntryDate{name, rsrc.ToDay(date)})
		}
	}

	return
}

func FilterEntryDates(entryDates []EntryDate, cutoff rsrc.Day,
) (filtered []EntryDate) {
	for _, entryDate := range entryDates {
		if entryDate.Date.Midnight() >= cutoff.Midnight() {
			filtered = append(filtered, entryDate)
		}
	}

	return
}

func (c Charts) GetYearPartition(registered rsrc.Day, threshold float64,
) Partition {
	entryDates := c.FindEntryDatesDynamic(registered, threshold)

	p := mapPart{
		assoc:      make(map[string]string),
		partitions: []string{},
	}

	for _, entryDate := range entryDates {
		p.assoc[entryDate.Name] = entryDate.Date.Time().Format("2006")
	}

	ii := newYearIterator(
		1,
		registered,
		rsrc.ToDay(registered.Midnight()+int64(86400*c.Len())))

	for ii.HasNext() {
		p.partitions = append(p.partitions, ii.Next().Begin.Time().Format("2006"))
	}
	p.partitions = append(p.partitions, "-")
	return p
}
