package charts

import (
	"fmt"
	"math"
	"runtime"
	"sort"
)

// Charts is table of daily accumulation of plays.
type Charts map[string][]float64

// Column is a column of charts sorted descendingly.
type Column []Score

// Score is a score with a name attached,
type Score struct {
	Name  string
	Score float64 // TODO rename Value
}

// Compile builds Charts from single day plays.
func Compile(days []Charts) Charts {
	size := len(days)
	charts := make(Charts)
	for i, day := range days {
		for name, plays := range day {
			if _, ok := charts[name]; !ok {
				charts[name] = make([]float64, size)
			}
			charts[name][i] = plays[0]
		}
	}

	return charts
}

// Sum computes partial sums for charts.
func (c Charts) Sum() Charts {
	return c.mapLine(func(in []float64, out []float64) {
		var sum float64
		for i, x := range in {
			sum += x
			out[i] = sum
		}
	})
}

func (c Charts) Fade(hl float64) Charts {
	fac := math.Pow(0.5, 1/hl)

	return c.mapLine(func(in []float64, out []float64) {
		sum := float64(0)
		for i, x := range in {
			sum *= fac
			sum += x
			out[i] = sum
		}
	})
}

func (c Charts) mapLine(f func(in []float64, out []float64)) Charts {
	result := make(Charts)

	lines := make(chan [2][]float64)
	workers := runtime.NumCPU()

	for i := 0; i < workers; i++ {
		go func() {
			for job := range lines {
				if job[0] == nil {
					break
				}

				f(job[0], job[1])
			}

		}()
	}

	for name, charts := range c {
		line := make([]float64, len(charts))
		result[name] = line
		lines <- [2][]float64{charts, line}
	}
	for i := 0; i < workers; i++ {
		lines <- [2][]float64{nil, nil}
	}
	close(lines)

	return result
}

// Column returns a column of charts sorted descendingly. Negative indices are
// used to index the chartes from behind.
func (c Charts) Column(i int) (column Column, err error) {
	var size int
	for _, line := range c {
		size = len(line)
		break
	}
	if i >= size {
		return Column{}, fmt.Errorf("Index %d >= %d (size)", i, size)
	}
	if i < 0 {
		i += size
	}
	if i < 0 {
		return Column{}, fmt.Errorf("Index %d < -%d (size)", i-size, size)
	}

	for name, line := range c {
		column = append(column, Score{name, line[i]})
	}
	sort.Sort(column)
	return column, nil
}

func (c Column) Len() int           { return len(c) }
func (c Column) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Column) Less(i, j int) bool { return c[i].Score > c[j].Score }

// Top returns the top n entries of col. If n is larger than len(col) the whole
// column is returned.
func (c Column) Top(n int) (top Column) {
	if n > len(c) {
		n = len(c)
	}
	return c[:n]
}
