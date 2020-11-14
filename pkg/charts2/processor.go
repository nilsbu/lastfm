package charts2

import (
	"math"
)

type LazyCharts interface {
	Row(title Title, begin, end int) []float64
	Column(titles []Title, index int) TitleValueMap
	Data(titles []Title, begin, end int) TitleLineMap

	Titles() []Title
	Len() int
}

func (l *charts) Row(title Title, begin, end int) []float64 {
	return l.values[title.Key()][begin:end]
}

func (l *charts) Column(titles []Title, index int) TitleValueMap {
	col := make(TitleValueMap)
	for _, t := range titles {
		col[t.Key()] = TitleValue{
			Title: t,
			Value: l.values[t.Key()][index],
		}
	}
	return col
}

func (l *charts) Data(titles []Title, begin, end int) TitleLineMap {
	data := make(TitleLineMap)
	for _, t := range titles {
		data[t.Key()] = TitleLine{
			Title: t,
			Line:  l.values[t.Key()][begin:end],
		}
	}
	return data
}

func (l *charts) Titles() []Title {
	ts := make([]Title, len(l.titles))

	for i, t := range l.titles {
		ts[i] = t
	}

	return ts
}

func (l *charts) Len() int {
	for _, line := range l.values {
		return len(line)
	}
	return -1
}

type chartsNode struct {
	parent LazyCharts
}

func (l chartsNode) Titles() []Title {
	return l.parent.Titles()
}

func (l chartsNode) Len() int {
	return l.parent.Len()
}

func Sum(parent LazyCharts) LazyCharts {
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(i int, line []float64) float64 {
			if i == 0 {
				acc = 0
			}
			acc += line[i]
			return acc
		},
		foldF: func(i int, line []float64) float64 {
			acc := 0.0
			for j := 0; j <= i; j++ {
				acc += line[j]
			}
			return acc
		},
		rangeF: fromBeginRange,
	}
}

func Fade(parent LazyCharts, hl float64) LazyCharts {
	fac := math.Pow(0.5, 1.0/hl)
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(i int, line []float64) float64 {
			if i == 0 {
				acc = 0
			}
			acc *= fac
			acc += line[i]
			return acc
		},
		foldF: func(i int, line []float64) float64 {
			acc := 0.0
			for j := 0; j <= i; j++ {
				acc *= fac
				acc += line[j]
			}
			return acc
		},
		rangeF: fromBeginRange,
	}
}

// Max calculates the maximum of the parent charts.
func Max(parent LazyCharts) LazyCharts {
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(i int, line []float64) float64 {
			if i == 0 {
				acc = line[i]
			} else {
				acc = math.Max(acc, line[i])
			}
			return acc
		},
		foldF: func(i int, line []float64) float64 {
			acc := 0.0
			for j := 0; j <= i; j++ {
				acc = math.Max(acc, line[j])
			}
			return acc
		},
		rangeF: fromBeginRange,
	}
}

// Gaussian blurs the data with a Gaussian kernel.
func Gaussian(
	parent LazyCharts,
	sigma float64,
	width int,
	mirrorBegin, mirrorEnd bool) LazyCharts {

	gaussian := make([]float64, width+1)
	fac := 1 / (sigma * math.Sqrt(2*math.Pi))
	for i := 0; i <= width; i++ {
		gaussian[i] = fac * math.Exp(-.5*float64(i*i)/sigma/sigma)
	}

	f := func(i int, line []float64) float64 {
		acc := 0.0
		b := i - width
		if b < 0 {
			if mirrorBegin {
				for j := 0; j < width-i; j++ {
					acc += gaussian[i+j+1] * line[j]
				}
			}

			b = 0
		}
		e := i + width + 1
		if e > len(line) {
			if mirrorEnd {
				for j := len(line) - 1; j > 2*len(line)-i-2-width; j-- {
					acc += gaussian[2*len(line)-(i+j+1)] * line[j]
				}
			}

			e = len(line)
		}

		for j := b; j < e; j++ {
			idx := j - i
			if idx < 0 {
				idx = -idx
			}
			acc += gaussian[idx] * line[j]
		}
		return acc
	}

	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF:       f,
		foldF:      f,
		rangeF: func(size, begin, end int) (b, e int) {
			if end+width > size {
				return 0, size
			}
			return 0, end + width
		},
	}
}

type lineProc func(i int, line []float64) float64
type rangeSpec func(size, begin, end int) (b, e int)

func fromBeginRange(size, begin, end int) (b, e int) {
	return 0, end
}

type lineMapCharts struct {
	chartsNode
	mapF, foldF lineProc
	rangeF      rangeSpec
}

func (l *lineMapCharts) Row(title Title, begin, end int) []float64 {
	rb, re := l.rangeF(l.parent.Len(), begin, end)
	in := l.parent.Row(title, rb, re)
	out := make([]float64, len(in))
	for i := range in {
		out[i] = l.mapF(i, in)
	}
	return out[begin-rb : end-rb]
}

func (l *lineMapCharts) Column(titles []Title, index int) TitleValueMap {
	col := make(TitleValueMap)
	back := make(chan TitleValue)
	rb, re := l.rangeF(l.parent.Len(), index, index+1)

	for t := range titles {
		go func(t int) {
			in := l.parent.Row(titles[t], rb, re)
			res := TitleValue{
				Title: titles[t],
			}
			res.Value = l.foldF(index, in)
			back <- res
		}(t)
	}

	for range titles {
		kf := <-back
		col[kf.Title.Key()] = kf
	}
	return col
}

func (l *lineMapCharts) Data(titles []Title, begin, end int) TitleLineMap {
	data := make(TitleLineMap)
	back := make(chan TitleLine)
	rb, re := l.rangeF(l.parent.Len(), begin, end)

	for k := range titles {
		go func(k int) {
			in := l.parent.Row(titles[k], rb, re)
			out := make([]float64, len(in))
			for i := range in {
				out[i] = l.mapF(i, in)
			}

			back <- TitleLine{
				Title: titles[k],
				Line:  out[begin-rb : end-rb],
			}
		}(k)
	}

	for range titles {
		tl := <-back
		data[tl.Title.Key()] = tl
	}
	return data
}

type partitionSum struct {
	chartsNode
	partition Partition
}

func Group(
	parent LazyCharts,
	partition Partition,
) LazyCharts {
	return &partitionSum{
		chartsNode: chartsNode{parent: parent},
		partition:  partition,
	}
}

func (l *partitionSum) Row(title Title, begin, end int) []float64 {
	titles := l.partition.Titles(title)
	back := make(chan []float64)

	for _, t := range titles {
		go func(t Title) {
			back <- l.parent.Row(t, begin, end)
		}(t)
	}

	var row []float64

	for i := 0; i < len(titles); i++ {
		line := <-back
		if len(row) == 0 {
			row = make([]float64, len(line))
		}
		for i, v := range line {
			row[i] += v
		}
	}
	return row
}

type titleColumn struct {
	key Title
	col TitleValueMap
}

func (l *partitionSum) Column(titles []Title, index int) TitleValueMap {
	col := make(TitleValueMap)
	back := make(chan titleColumn)

	for _, bin := range titles {
		ts := []Title{}
		for _, r := range l.partition.Titles(bin) {
			ts = append(ts, r)
		}
		go func(titles []Title, bin Title) {
			back <- titleColumn{
				key: bin,
				col: l.parent.Column(titles, index),
			}
		}(ts, bin)
	}

	for range titles {
		kf := <-back
		for _, v := range kf.col {
			col[kf.key.Key()] = TitleValue{
				Title: kf.key,
				Value: col[kf.key.Key()].Value + v.Value,
			}
		}
	}

	return col
}

func (l *partitionSum) Data(titles []Title, begin, end int) TitleLineMap {
	data := make(TitleLineMap)
	back := make(chan TitleLine)

	n := 0
	for _, bin := range titles {
		for _, key := range l.partition.Titles(bin) {
			go func(key, bin Title) {
				back <- TitleLine{
					Title: bin,
					Line:  l.parent.Row(key, begin, end),
				}
			}(key, bin)
			n++
		}
	}

	for i := 0; i < n; i++ {
		kl := <-back
		if _, ok := data[kl.Title.Key()]; !ok {
			data[kl.Title.Key()] = TitleLine{
				Title: kl.Title,
				Line:  make([]float64, len(kl.Line)),
			}
		}
		line := data[kl.Title.Key()].Line
		for i, v := range kl.Line {
			line[i] += v
		}
		data[kl.Title.Key()] = TitleLine{
			Title: data[kl.Title.Key()].Title,
			Line:  line,
		}
	}

	return data
}

func (l *partitionSum) Titles() []Title {
	return l.partition.Partitions()
}

// ColumnSum is a LazyCharts that sums up all columns.
// TODO can be optmized by bypassing partitions
func ColumnSum(parent LazyCharts) LazyCharts {
	return Group(parent, totalPartition(parent.Titles()))
}

type cache struct {
	chartsNode
	charts charts
	ranges map[string][2]int
}

// Cache is a LazyCharts that stores data to avoid duplicating work in parent.
// The cache is filled when the data is requested. The data is stored in one
// continuous block per row. Non-requested parts in between are filled.
// E.g. if Row("A", 0, 4) and Column({"A"}, 16) are called, row "A" will store
// range [0, 17).
func Cache(parent LazyCharts) LazyCharts {
	return &cache{
		chartsNode: chartsNode{parent},
		charts: charts{
			titles: make([]Title, 0),
			values: make(map[string][]float64),
		},
		ranges: make(map[string][2]int),
	}
}

func (c *cache) Row(title Title, begin, end int) []float64 {
	if line, ok := c.charts.values[title.Key()]; ok {
		r := c.ranges[title.Key()]
		if r[0] <= begin && r[1] >= end {
			return line[begin-r[0] : end-r[0]]
		}

		if begin < r[0] {
			line = append(c.parent.Row(title, begin, r[0]), line...)
			r[0] = begin
		}
		if r[1] < end {
			line = append(line, c.parent.Row(title, r[1], end)...)
			r[1] = end
		}

		c.charts.values[title.Key()] = line
		c.ranges[title.Key()] = r
		return line[begin-r[0] : end-r[0]]
	}

	pLine := c.parent.Row(title, begin, end)
	c.charts.values[title.Key()] = pLine
	c.ranges[title.Key()] = [2]int{begin, end}
	return pLine
}

func (c *cache) Column(titles []Title, index int) TitleValueMap {
	data := c.Data(titles, index, index+1)
	tvm := make(TitleValueMap)

	for k, l := range data {
		tvm[k] = TitleValue{
			Title: l.Title,
			Value: l.Line[0],
		}
	}

	return tvm
}

func (c *cache) Data(titles []Title, begin, end int) TitleLineMap {
	data := make(TitleLineMap)
	back := make(chan TitleLine)

	for k := range titles {
		go func(k int) {
			back <- TitleLine{
				Title: titles[k],
				Line:  c.Row(titles[k], begin, end),
			}
		}(k)
	}

	for range titles {
		tl := <-back
		data[tl.Title.Key()] = tl
	}
	return data
}
