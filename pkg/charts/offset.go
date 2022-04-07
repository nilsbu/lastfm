package charts

import (
	"github.com/nilsbu/async"
)

type offset struct {
	chartsNode
	offsets map[string]int
	titles  []Title
}

// Offset shifts each line by a certain offset. offsets maps the Title key to the offset.
func Offset(parent Charts, offsets map[string]int) Charts {
	return &offset{chartsNode{parent: parent}, offsets, nil}
}

func (c *offset) Titles() []Title {
	if c.titles == nil {
		pTitles := c.parent.Titles()
		c.titles = make([]Title, len(c.offsets))
		i := 0
		for _, title := range pTitles {
			if _, ok := c.offsets[title.Key()]; ok {
				c.titles[i] = title
				i++
			}
		}
	}

	return c.titles
}

func (c *offset) Data(titles []Title, begin, end int) ([][]float64, error) {
	data := make([][]float64, len(titles))
	length := c.parent.Len()

	async.Pie(len(titles), func(i int) error {
		offset := c.offsets[titles[i].Key()]
		line := make([]float64, end-begin)

		e := end + offset
		if e > length {
			e = length
		}

		o := begin + offset
		if o >= length {
			o = length - 1
		}

		d, err := c.parent.Data([]Title{titles[i]}, o, e)
		if err != nil {
			return err
		}
		copy(line, d[0])

		if len(d[0]) > 0 {
			for j := len(d[0]); j < len(line); j++ {
				line[j] = line[len(d[0])-1]
			}
		}

		data[i] = line
		return nil
	})
	return data, nil
}

func EntryDates(gaussian, sums Charts) (map[string]int, error) {
	titles := sums.Titles()
	s, err := sums.Data(titles, 0, sums.Len())
	if err != nil {
		return nil, err
	}
	g, err := gaussian.Data(titles, 0, gaussian.Len())
	if err != nil {
		return nil, err
	}

	entries := make([]int, len(titles))
	async.Pi(len(titles), func(i int) {
		gLine := g[i]
		sLine := s[i]

		if sLine[len(sLine)-1] < 2 {
			entries[i] = -1
			return
		}

		begin, end := -1, -1
		for j, v := range sLine {
			if begin < 0 {
				if v >= 2 {
					begin = j
				}
			} else if v > 4 {
				end = j
				break
			}
		}

		m, mj := 0.0, 0
		for j := begin; j < end; j++ {
			if m < gLine[j] {
				m = gLine[j]
				mj = j
			}
		}

		entries[i] = mj
	})

	eMap := make(map[string]int)
	for i, entry := range entries {
		if entry > 0 {
			eMap[titles[i].Key()] = entry
		}
	}

	return eMap, nil
}
