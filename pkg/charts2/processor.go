package charts2

import (
	"math"
)

type LazyCharts interface {
	Row(key Title, begin, end int) []float64
	Column(keys []Title, index int) map[string]float64
	Data(keys []Title, begin, end int) map[string][]float64

	Titles() []Title
	Len() int
}

func (l *Charts) Row(key Title, begin, end int) []float64 {
	return l.Values[key.Key()][begin:end]
}

func (l *Charts) Column(keys []Title, index int) map[string]float64 {
	col := make(map[string]float64)
	for _, k := range keys {
		col[k.Key()] = l.Values[k.Key()][index]
	}
	return col
}

func (l *Charts) Data(keys []Title, begin, end int) map[string][]float64 {
	data := make(map[string][]float64)
	for _, k := range keys {
		data[k.Key()] = l.Values[k.Key()][begin:end]
	}
	return data
}

func (l *Charts) Titles() []Title {
	ts := make([]Title, len(l.titles))

	for i, t := range l.titles {
		ts[i] = t
	}

	return ts
}

func (l *Charts) Len() int {
	for _, line := range l.Values {
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

type lineMap func(i int, v float64) float64

func Sum(parent LazyCharts) LazyCharts {
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		f: func(i int, v float64) float64 {
			if i == 0 {
				acc = 0
			}
			acc += v
			return acc
		},
	}
}

func Fade(parent LazyCharts, hl float64) LazyCharts {
	fac := math.Pow(0.5, 1.0/hl)
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		f: func(i int, v float64) float64 {
			if i == 0 {
				acc = 0
			}
			acc *= fac
			acc += v
			return acc
		},
	}
}

// Max calculates the maximum of the parent charts.
func Max(parent LazyCharts) LazyCharts {
	acc := 0.0
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		f: func(i int, v float64) float64 {
			if i == 0 {
				acc = v
			} else {
				acc = math.Max(acc, v)
			}
			return acc
		},
	}
}

type lineMapCharts struct {
	chartsNode
	f lineMap
}

func (l *lineMapCharts) Row(key Title, begin, end int) []float64 {

	in := l.parent.Row(key, 0, end)
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = l.f(i, v)
	}
	return out[begin:]
}

type keyFloat struct {
	key   Title
	value float64
}

type keyLine struct {
	key  Title
	line []float64
}

func (l *lineMapCharts) Column(keys []Title, index int) map[string]float64 {
	col := make(map[string]float64)
	back := make(chan keyFloat)

	for k := range keys {
		go func(k int) {
			in := l.parent.Row(keys[k], 0, index+1)
			res := keyFloat{
				key: keys[k],
			}
			for i, v := range in {
				res.value = l.f(i, v)
			}
			back <- res
		}(k)
	}

	for range keys {
		kf := <-back
		col[kf.key.Key()] = kf.value
	}
	return col
}

func (l *lineMapCharts) Data(keys []Title, begin, end int) map[string][]float64 {
	data := make(map[string][]float64)

	titles := keys
	back := make(chan keyLine)

	for k := range titles {
		go func(k int) {
			in := l.parent.Row(titles[k], 0, end)
			out := make([]float64, len(in))
			for i, v := range in {
				out[i] = l.f(i, v)
			}

			back <- keyLine{
				key:  titles[k],
				line: out[begin:],
			}
		}(k)
	}

	for range titles {
		kl := <-back
		data[kl.key.Key()] = kl.line
	}
	return data
}

type partitionSum struct {
	chartsNode
	partition map[string]string
	key       func(Title) string
}

func (l *partitionSum) Row(key Title, begin, end int) []float64 {
	titles := inverseMap(l.partition)[key.Key()]
	back := make(chan keyLine)

	for _, t := range titles {
		go func(t Title) {
			back <- keyLine{
				line: l.parent.Row(t, begin, end),
			}
		}(KeyTitle(t))
	}

	var row []float64

	for i := 0; i < len(titles); i++ {
		kl := <-back
		if len(row) == 0 {
			row = make([]float64, len(kl.line))
		}
		for i, v := range kl.line {
			row[i] += v
		}
	}
	return row
}

func inverseMap(keys map[string]string) map[string][]string {
	rev := map[string][]string{}
	for k, p := range keys {
		if _, ok := rev[p]; !ok {
			rev[p] = []string{k}
		} else {
			rev[p] = append(rev[p], k)
		}
	}
	return rev
}

type keyColumn struct {
	key string
	col map[string]float64
}

func (l *partitionSum) Column(keys []Title, index int) map[string]float64 {
	col := make(map[string]float64)
	rev := inverseMap(l.partition)
	back := make(chan keyColumn)

	for _, bin := range keys {
		ts := []Title{}
		for _, r := range rev[bin.Key()] {
			ts = append(ts, KeyTitle(r))
		}
		go func(keys []Title, bin Title) {
			back <- keyColumn{
				key: bin.Key(),
				col: l.parent.Column(keys, index),
			}

		}(ts, bin)
	}

	for range keys {
		kf := <-back
		for _, v := range kf.col {
			col[kf.key] += v
		}
	}

	return col
}

func (l *partitionSum) Data(keys []Title, begin, end int) map[string][]float64 {
	data := make(map[string][]float64)
	rev := inverseMap(l.partition)
	back := make(chan keyLine)

	n := 0
	for _, bin := range keys {
		for _, key := range rev[bin.Key()] {
			go func(key, bin Title) {
				back <- keyLine{
					key:  bin,
					line: l.parent.Row(key, begin, end),
				}
			}(KeyTitle(key), bin)
			n++
		}
	}

	for i := 0; i < n; i++ {
		kl := <-back
		if _, ok := data[kl.key.Key()]; !ok {
			data[kl.key.Key()] = make([]float64, len(kl.line))
		}
		line := data[kl.key.Key()]
		for i, v := range kl.line {
			line[i] += v
		}
		data[kl.key.Key()] = line
	}

	return data
}

func (l *partitionSum) Titles() []Title {
	// OPTIMIZE: doesn't need full lookup
	set := map[string]Title{}
	for _, v := range l.partition {
		set[v] = KeyTitle(v)
	}
	keys := make([]Title, 0)
	for _, t := range set {
		keys = append(keys, t)
	}
	return keys
}
