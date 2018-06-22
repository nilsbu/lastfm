package charts

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// Charts is table of daily accumulation of plays.
type Charts map[rsrc.Name][]int

// Sums are special charts where each row consists of the some from the
// beginning until the current row.
type Sums Charts

// Column is a column of charts sorted descendingly.
type Column []Score

// Score is a score with a name attached,
type Score struct {
	Name  rsrc.Name
	Score int
}

// Compile builds Charts from DayPlays.
func Compile(dayPlays []unpack.DayPlays) Charts {
	size := len(dayPlays)
	charts := make(Charts)
	for i, day := range dayPlays {
		for name, plays := range day {
			if _, ok := charts[name]; !ok {
				charts[name] = make([]int, size)
			}
			charts[name][i] = plays
		}
	}

	return charts
}

// Sum computes partial sums for charts.
func (c Charts) Sum() Sums {
	sums := make(Sums)

	lines := make(chan [2][]int)
	workers := runtime.NumCPU()

	for i := 0; i < workers; i++ {
		go func() {
			for job := range lines {
				if job[0] == nil {
					break
				}

				sum := 0
				for i, x := range job[1] {
					sum += x
					job[0][i] = sum
				}
			}

		}()
	}

	for name, charts := range c {
		line := make([]int, len(charts))
		sums[name] = line
		lines <- [2][]int{line, charts}
	}
	for i := 0; i < workers; i++ {
		lines <- [2][]int{nil, nil}
	}
	close(lines)

	return sums
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
