package charts

// Score is a score with a name attached,
type Score struct {
	Name  string
	Score float64 // TODO rename Value
}

// Column is a column of charts sorted descendingly.
type Column []Score

func (c Column) Len() int           { return len(c) }
func (c Column) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Column) Less(i, j int) bool { return c[i].Score > c[j].Score }

// Sum sums over all values in a column.
func (c Column) Sum() (sum float64) {
	for _, line := range c {
		sum += line.Score
	}

	return sum
}

// TODO sort file by receiver

// Top returns the top n entries of col. If n is larger than len(col) the whole
// column is returned.
func (c Column) Top(n int) (top Column) {
	if n > len(c) {
		n = len(c)
	}
	return c[:n]
}
