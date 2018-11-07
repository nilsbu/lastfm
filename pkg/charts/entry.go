package charts

import (
	"time"

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
				date, _ := registered.Midnight()
				date += int64(86400 * i)
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
			date, _ := registered.Midnight()
			date += int64(86400 * idx)
			entryDates = append(entryDates, EntryDate{name, rsrc.ToDay(date)})
		}
	}

	return
}

func FilterEntryDates(entryDates []EntryDate, cutoff rsrc.Day,
) (filtered []EntryDate) {
	cutoffM, _ := cutoff.Midnight()
	for _, entryDate := range entryDates {
		date, _ := entryDate.Date.Midnight()
		if date >= cutoffM {
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
		m, _ := entryDate.Date.Midnight()
		p.assoc[entryDate.Name] = time.Unix(m, 0).UTC().Format("2006")
	}

	reg, _ := registered.Midnight()
	ii := newIntervalIterator(
		Year,
		time.Unix(reg, 0).UTC(),
		reg+int64(86400*c.Len()))

	for ii.HasNext() {
		p.partitions = append(p.partitions, ii.Next().Begin.Format("2006"))
	}
	p.partitions = append(p.partitions, "")
	return p
}
