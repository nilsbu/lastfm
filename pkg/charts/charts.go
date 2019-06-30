package charts

import (
	"fmt"
	"math"
	"runtime"
	"sort"
)

// Charts is table of daily accumulation of plays.
type Charts map[string][]float64

// Compile builds Charts from single day plays.
func Compile(days []map[string][]float64) Charts {
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

// UnravelDays takes Charts and disassembles it into single day plays. It acts
// as an inverse to Compile().
func (c Charts) UnravelDays() []map[string][]float64 {
	days := []map[string][]float64{}
	for i := 0; i < c.Len(); i++ {
		day := map[string][]float64{}

		for name, line := range c {
			if line[i] != 0 {
				day[name] = []float64{line[i]}
			}
		}

		days = append(days, day)
	}

	return days
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

func (c Charts) Len() int {
	for _, line := range c {
		return len(line)
	}

	return 0
}

// Column returns a column of charts sorted descendingly. Negative indices are
// used to index the chartes from behind.
func (c Charts) Column(i int) (column Column, err error) {
	size := c.Len()
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

// Keys returns the keys of the charts.
func (c Charts) Keys() []string {
	keys := []string{}
	for key := range c {
		keys = append(keys, key)
	}
	return keys
}

func (c Charts) Correct(replace map[string]string) Charts {
	corrected := Charts{}

	for key, values := range c {
		corrected[key] = values
	}

	for key := range c {
		if with, ok := replace[key]; ok {
			dest := corrected[with]
			src := corrected[key]
			sum := make([]float64, len(dest))

			for i := range dest {
				sum[i] = src[i] + dest[i]
			}

			delete(corrected, key)
			corrected[with] = sum
		}
	}

	return corrected
}

// append a column at the end of the charts. The keys are not required to be in
// the charts prior.
func (c Charts) append(col Column) {
	for _, score := range col {
		c[score.Name] = append(c[score.Name], score.Score)
	}
}

// Rank the charts in each column.
func (c Charts) Rank() (ranks Charts) {
	ranks = make(Charts)

	for i := 0; i < c.Len(); i++ {
		col, _ := c.Column(i)

		var last float64
		idx := 1
		for j, score := range col {
			if last != score.Score {
				idx = j + 1
				last = score.Score
			}
			col[j].Score = float64(idx)
		}

		ranks.append(col)
	}

	return
}

type totalPartition struct{}

func (totalPartition) Partitions() []string {
	return []string{""}
}

func (totalPartition) Get(key string) string {
	return ""
}

func (c Charts) Total() []float64 {
	return c.Group(totalPartition{})[""]
}

// Max returns a Column where the score for each key is equal to the maximum of
// all scores in that key's line.
func (c Charts) Max() (max Column) {
	max = Column{}

	for key, values := range c {
		m := 0.0
		for _, v := range values {
			m = math.Max(m, v)
		}
		max = append(max, Score{Name: key, Score: m})
	}

	sort.Sort(max)

	return
}
