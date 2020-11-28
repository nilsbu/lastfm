package charts2

type normalizer struct {
	chartsNode
	totals LazyCharts
}

func NormalizeColumn(c LazyCharts) LazyCharts {
	return &normalizer{
		chartsNode: chartsNode{parent: c},
		totals:     Cache(ColumnSum(c)),
	}
}

func NormalizeGaussian(
	c LazyCharts,
	sigma float64,
	width int,
	mirrorBegin, mirrorEnd bool) LazyCharts {
	smooth := Cache(Gaussian(c, sigma, width, mirrorBegin, mirrorEnd))
	return &normalizer{
		chartsNode: chartsNode{parent: smooth},
		totals:     ColumnSum(smooth),
	}
}

func (c *normalizer) Row(title Title, begin, end int) []float64 {
	return c.Data([]Title{title}, begin, end)[title.Key()].Line
}

func (c *normalizer) Column(titles []Title, index int) TitleValueMap {
	total := c.totals.Row(StringTitle("total"), index, index+1)[0]

	col := make(TitleValueMap)

	if total == 0 {
		for _, title := range titles {
			col[title.Key()] = TitleValue{
				Title: title,
				Value: 0.0,
			}
		}
		return col
	}

	par := c.parent.Column(titles, index)

	for k, tv := range par {
		col[k] = TitleValue{
			Title: tv.Title,
			Value: tv.Value / total,
		}
	}
	return col
}

func (c *normalizer) Data(titles []Title, begin, end int) TitleLineMap {
	totals := c.totals.Row(StringTitle("total"), begin, end)

	data := make(TitleLineMap)
	back := make(chan TitleLine)

	for k := range titles {
		go func(k int) {
			in := c.parent.Row(titles[k], begin, end)
			out := make([]float64, len(in))
			for i, v := range in {
				if totals[i] > 0 {
					out[i] = v / totals[i]
				}
			}

			back <- TitleLine{
				Title: titles[k],
				Line:  out,
			}
		}(k)
	}

	for range titles {
		tl := <-back
		data[tl.Title.Key()] = tl
	}
	return data
}
