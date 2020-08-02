package charts

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"
)

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
	result := make([][]float64, len(c.Keys))

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

	for li, lineIn := range c.Values {
		lineOut := make([]float64, len(lineIn))
		result[li] = lineOut
		lines <- [2][]float64{lineIn, lineOut}
	}
	for i := 0; i < workers; i++ {
		lines <- [2][]float64{nil, nil}
	}
	close(lines)

	return Charts{
		Headers: c.Headers,
		Keys:    c.Keys,
		Values:  result,
	}
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

	for li, line := range c.Values {
		column = append(column, Score{c.Keys[li].String(), line[i]})
	}
	sort.Sort(column)
	return column, nil
}

func (c Charts) Correct(replace map[string]string) Charts {
	corrected := map[string][]float64{}

	for li, line := range c.Values {
		corrected[c.Keys[li].String()] = line
	}

	for _, key := range c.Keys {
		if with, ok := replace[key.String()]; ok {
			dest := corrected[with]
			src := corrected[key.String()]
			sum := make([]float64, len(dest))

			for i := range dest {
				sum[i] = src[i] + dest[i]
			}

			delete(corrected, key.String())
			corrected[with] = sum
		}
	}

	keys := []Key{}
	values := [][]float64{}

	for name, plays := range corrected {
		// TODO key may not be simple
		keys = append(keys, simpleKey(name))
		values = append(values, plays)
	}

	return Charts{
		Headers: c.Headers,
		Keys:    keys,
		Values:  values,
	}
}

// Rank the charts in each column.
func (c Charts) Rank() (ranks Charts) {
	ranks.Headers = c.Headers
	ranks.Keys = c.Keys
	ranks.Values = make([][]float64, len(c.Keys))
	for i := 0; i < len(ranks.Values); i++ {
		ranks.Values[i] = make([]float64, c.Len())
	}

	back := make(chan bool)
	for i := 0; i < c.Len(); i++ {
		go func(i int) {
			col, _ := c.Column(i)
			var last float64
			idx := 1
			for j, score := range col {
				if last != score.Score {
					idx = j + 1
					last = score.Score
				}

				for k, key := range ranks.Keys {
					if key.String() == score.Name {
						ranks.Values[k][i] = float64(idx)
						break
					}
				}
			}

			back <- true
		}(i)
	}

	for i := 0; i < c.Len(); i++ {
		<-back
	}

	return
}

// TotalPartition is a partition in which all artists belong to the same
// partition.
type TotalPartition struct{}

func (TotalPartition) Partitions() []Key {
	return []Key{simpleKey("total")}
}

func (TotalPartition) Get(key Key) Key {
	return simpleKey("")
}

func (c Charts) Total() []float64 {
	total := c.Group(TotalPartition{})
	if len(total.Values) > 0 {
		return total.Values[0]
	} else {
		return []float64{}
	}
}

// Max returns a Column where the score for each key is equal to the maximum of
// all scores in that key's line.
func (c Charts) Max() (max Column) {
	max = Column{}

	for i, key := range c.Keys {
		m := 0.0
		for _, v := range c.Values[i] {
			m = math.Max(m, v)
		}
		max = append(max, Score{Name: key.String(), Score: m})
	}

	sort.Sort(max)

	return
}

// Equal compares two charts in their headers, keys and values. Key order does
// not matter.
func (c Charts) Equal(other Charts) bool {
	return c.AssertEqual(other) == nil
}

// AssertEqual compares two charts in their headers, keys and values. Key order
// does not matter. If they differ an error is returned.
func (c Charts) AssertEqual(other Charts) error {
	if c.Len() != other.Len() {
		return fmt.Errorf("this len is '%v' but other is '%v'", c.Len(), other.Len())
	}

	// Compare headers
	if c.Headers == nil {
		return errors.New("this header is nil")
	}
	if other.Headers == nil {
		return errors.New("this header is nil")
	}

	if c.Headers.Len() != other.Headers.Len() {
		return fmt.Errorf("this header's len is '%v' but other is '%v'",
			c.Headers.Len(), other.Headers.Len())
	}

	for i := 0; i < c.Len(); i++ {
		thisInterval := c.Headers.At(i)
		otherInterval := other.Headers.At(i)
		if thisInterval.Begin.Midnight() != otherInterval.Begin.Midnight() {
			return fmt.Errorf("begin at position %v is %v but other is %v",
				i, thisInterval.Begin, otherInterval.Begin)
		}
		if thisInterval.Before.Midnight() != otherInterval.Before.Midnight() {
			return fmt.Errorf("before at position %v is %v but other is %v",
				i, thisInterval.Before, otherInterval.Before)
		}

		thisI := c.Headers.Index(thisInterval.Begin)
		otherI := other.Headers.Index(thisInterval.Begin)
		if thisI != otherI {
			return fmt.Errorf("index of '%v' is '%v' but other is '%v'",
				c.Headers.At(i).Begin, i, otherI)
		}
	}

	// Compare charts
	otherMap := map[string][]float64{}
	for i, key := range other.Keys {
		otherMap[key.String()] = other.Values[i]
	}

	thisMap := map[string][]float64{}
	for i, key := range c.Keys {
		thisMap[key.String()] = c.Values[i]
	}

	if !reflect.DeepEqual(thisMap, otherMap) {
		return fmt.Errorf("lines differ: '%v' != '%v'", thisMap, otherMap)
	}

	// Test Keys in detail
	for _, key := range c.Keys {
		for _, otherKey := range other.Keys {
			if key.String() == otherKey.String() {
				if key.ArtistName() != otherKey.ArtistName() {
					return fmt.Errorf("artist of key '%v' differs: '%v' vs. '%v'",
						key.String(), key.ArtistName(), otherKey.ArtistName())
				}
				if key.FullTitle() != otherKey.FullTitle() {
					return fmt.Errorf("full title of key '%v' differs: '%v' vs. '%v'",
						key.String(), key.FullTitle(), otherKey.FullTitle())
				}
			}
		}
	}

	return nil
}

type artistNamePartition struct {
	artistPartition
}

// NewArtistNamePartition creates a Partition in which each artist in charts is
// its own partition.
func NewArtistNamePartition(charts Charts) Partition {
	p := artistNamePartition{}
	p.assoc = make(map[string]Key)
	for _, key := range charts.Keys {
		if _, ok := p.assoc[key.ArtistName()]; !ok {
			p.assoc[key.ArtistName()] = simpleKey(key.ArtistName())
		}
	}

	for _, v := range p.assoc {
		p.partitions = append(p.partitions, v)
	}
	return p
}
