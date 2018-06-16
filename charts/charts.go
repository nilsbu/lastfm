package charts

import (
	"runtime"

	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/unpack"
)

// Charts is table of daily accumulation of plays.
type Charts map[io.Name][]int

// Sums are special charts where each row consists of the some from the
// beginning until the current row.
type Sums Charts

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
