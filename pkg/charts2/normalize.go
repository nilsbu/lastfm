package charts2

type normalizer struct {
	chartsNode
	totals []float64
	kernel []float64
}

func SingleColumnNormalizer(c LazyCharts) LazyCharts {
	return &normalizer{
		chartsNode: chartsNode{parent: c},
		kernel:     []float64{1},
	}
}

func (c *normalizer) Row(title Title, begin, end int) []float64 {
	c.calcTotals()
	return c.Data([]Title{title}, begin, end)[title.Key()].Line
}

func (c *normalizer) Column(titles []Title, index int) TitleValueMap {
	c.calcTotals()

	col := make(TitleValueMap)
	back := make(chan TitleValue)

	for t := range titles {
		go func(t int) {
			in := c.parent.Row(titles[t], 0, index+1)
			res := TitleValue{
				Title: titles[t],
			}
			for _, v := range in {
				if c.totals[index] > 0 {
					res.Value = v / c.totals[index]
				}
			}
			back <- res
		}(t)
	}

	for range titles {
		kf := <-back
		col[kf.Title.Key()] = kf
	}
	return col
}

func (c *normalizer) Data(titles []Title, begin, end int) TitleLineMap {
	c.calcTotals()

	data := make(TitleLineMap)
	back := make(chan TitleLine)

	for k := range titles {
		go func(k int) {
			in := c.parent.Row(titles[k], 0, end)
			out := make([]float64, len(in))
			for i, v := range in {
				if c.totals[i] > 0 {
					out[i] = v / c.totals[i]
				}
			}

			back <- TitleLine{
				Title: titles[k],
				Line:  out[begin:],
			}
		}(k)
	}

	for range titles {
		tl := <-back
		data[tl.Title.Key()] = tl
	}
	return data
}

func (c *normalizer) calcTotals() {
	if len(c.totals) > 0 {
		return
	}

	c.totals = make([]float64, c.Len())
	data := c.parent.Data(c.Titles(), 0, c.Len())

	for _, line := range data {
		for i, v := range line.Line {
			c.totals[i] += v
		}
	}
}
