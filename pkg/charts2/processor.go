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
		mapF: func(i int, v float64) float64 {
			if i == 0 {
				acc = 0
			}
			acc += v
			return acc
		},
		foldF: func(line []float64) float64 {
			acc := 0.0
			for _, v := range line {
				acc += v
			}
			return acc
		},
	}
}

func Fade(parent LazyCharts, hl float64) LazyCharts {
	fac := math.Pow(0.5, 1.0/hl)
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(i int, v float64) float64 {
			if i == 0 {
				acc = 0
			}
			acc *= fac
			acc += v
			return acc
		},
		foldF: func(line []float64) float64 {
			acc := 0.0
			for _, v := range line {
				acc *= fac
				acc += v
			}
			return acc
		},
	}
}

// Max calculates the maximum of the parent charts.
func Max(parent LazyCharts) LazyCharts {
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(i int, v float64) float64 {
			if i == 0 {
				acc = v
			} else {
				acc = math.Max(acc, v)
			}
			return acc
		},
		foldF: func(line []float64) float64 {
			acc := 0.0
			for _, v := range line {
				acc = math.Max(acc, v)
			}
			return acc
		},
	}
}

type lineMap func(i int, v float64) float64
type lineFold func(line []float64) float64

type lineMapCharts struct {
	chartsNode
	mapF  lineMap
	foldF lineFold
}

func (l *lineMapCharts) Row(title Title, begin, end int) []float64 {

	in := l.parent.Row(title, 0, end)
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = l.mapF(i, v)
	}
	return out[begin:]
}

func (l *lineMapCharts) Column(titles []Title, index int) TitleValueMap {
	col := make(TitleValueMap)
	back := make(chan TitleValue)

	for t := range titles {
		go func(t int) {
			in := l.parent.Row(titles[t], 0, index+1)
			res := TitleValue{
				Title: titles[t],
			}
			res.Value = l.foldF(in)
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

	for k := range titles {
		go func(k int) {
			in := l.parent.Row(titles[k], 0, end)
			out := make([]float64, len(in))
			for i, v := range in {
				out[i] = l.mapF(i, v)
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
	key string
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
				key: bin.Key(),
				col: l.parent.Column(titles, index),
			}
		}(ts, bin)
	}

	for range titles {
		kf := <-back
		for _, v := range kf.col {
			col[kf.key] = TitleValue{
				Title: KeyTitle(kf.key),
				Value: col[kf.key].Value + v.Value,
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
