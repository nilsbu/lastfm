package charts

import (
	"github.com/nilsbu/async"
)

type offset struct {
	chartsNode
	offsets map[string]int
}

// Offset shifts each line by a certain offset. offsets maps the Title key to the offset.
func Offset(parent Charts, offsets map[string]int) Charts {
	return &offset{chartsNode{parent: parent}, offsets}
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
