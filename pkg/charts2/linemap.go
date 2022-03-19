package charts2

import "math"

type lineMapCharts struct {
	chartsNode
	mapF   lineProc
	foldF  valueProc
	rangeF rangeSpec
}

type valueProc func(i int, line []float64) float64
type lineProc func(line []float64) []float64
type rangeSpec func(size, begin, end int) (b, e int)

func Sum(parent LazyCharts) LazyCharts {
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			acc := 0.0
			for i := range in {
				acc += in[i]
				out[i] = acc
			}
			return out
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
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			acc := 0.0
			for i := range in {
				acc *= fac
				acc += in[i]
				out[i] = acc
			}
			return out
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
	return &lineMapCharts{
		chartsNode: chartsNode{parent: parent},
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			acc := -math.MaxFloat64
			for i := range in {
				acc = math.Max(acc, in[i])
				out[i] = acc
			}
			return out
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
				ee := width - i
				if ee > len(line) {
					ee = len(line)
				}
				for j := 0; j < ee; j++ {
					acc += gaussian[i+j+1] * line[j]
				}
			}

			b = 0
		}
		e := i + width + 1
		if e > len(line) {
			if mirrorEnd {
				bb := 2*len(line) - i - 1 - width
				if bb < 0 {
					bb = 0
				}
				for j := len(line) - 1; j >= bb; j-- {
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
		mapF: func(in []float64) []float64 {
			out := make([]float64, len(in))
			for i := range in {
				out[i] = f(i, in)
			}
			return out
		},
		foldF: f,
		rangeF: func(size, begin, end int) (b, e int) {
			if end+width > size {
				return 0, size
			}
			return 0, end + width
		},
	}
}

func (l *lineMapCharts) Column(titles []Title, index int) []float64 {
	type iv struct {
		i int
		v float64
	}

	col := make([]float64, len(titles))
	back := make(chan iv)
	rb, re := l.rangeF(l.parent.Len(), index, index+1)

	for t := range titles {
		go func(t int) {
			in := l.parent.Data([]Title{titles[t]}, rb, re)[0]
			back <- iv{
				i: t,
				v: l.foldF(index, in),
			}
		}(t)
	}

	for range titles {
		kf := <-back
		col[kf.i] = kf.v
	}
	return col
}

func (l *lineMapCharts) Data(titles []Title, begin, end int) [][]float64 {

	data := make([][]float64, len(titles))
	back := make(chan indexLine)
	rb, re := l.rangeF(l.parent.Len(), begin, end)

	for k := range titles {
		go func(k int) {
			in := l.parent.Data([]Title{titles[k]}, rb, re)[0]
			out := l.mapF(in)

			back <- indexLine{
				i:  k,
				vs: out[begin-rb : end-rb],
			}
		}(k)
	}

	for range titles {
		il := <-back
		data[il.i] = il.vs
	}
	return data
}

type indexLine struct {
	i  int
	vs []float64
}
