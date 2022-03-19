package charts2

func fromBeginRange(size, begin, end int) (b, e int) {
	return 0, end
}

type partitionSum struct {
	chartsNode
	partition Partition
}

// Group is a LazyCharts that combines the subsets of the partition from the parent.
func Group(
	parent LazyCharts,
	partition Partition,
) LazyCharts {
	return &partitionSum{
		chartsNode: chartsNode{parent: parent},
		partition:  partition,
	}
}

type titleColumn struct {
	key int
	col []float64
}

func (l *partitionSum) Column(titles []Title, index int) []float64 {
	col := make([]float64, len(titles))
	back := make(chan titleColumn)

	for i, bin := range titles {
		ts := l.partition.Titles(bin)
		go func(titles []Title, i int) {
			back <- titleColumn{
				key: i,
				col: l.parent.Column(titles, index),
			}
		}(ts, i)
	}

	for range titles {
		kf := <-back
		for _, v := range kf.col {
			col[kf.key] += v
		}
	}

	return col
}

func (l *partitionSum) Data(titles []Title, begin, end int) [][]float64 {
	back := make(chan indexLine)

	for i, bin := range titles {
		go func(i int, bin Title) {
			line := make([]float64, end-begin)
			for _, key := range l.partition.Titles(bin) {
				for j, v := range l.parent.Data([]Title{key}, begin, end)[0] {
					line[j] += v
				}
			}
			back <- indexLine{
				i:  i,
				vs: line,
			}
		}(i, bin)
	}

	data := make([][]float64, len(titles))
	for i := 0; i < len(titles); i++ {
		b := <-back
		data[b.i] = b.vs
	}

	return data
}

func (l *partitionSum) Titles() []Title {
	return l.partition.Partitions()
}

// Subset is a LazyCharts that picks a single subset of the partition from the parent.
func Subset(
	parent LazyCharts,
	partition Partition,
	title Title,
) LazyCharts {
	return Only(parent, partition.Titles(title))
}

// ColumnSum is a LazyCharts that sums up all columns.
// TODO can be optmized by bypassing partitions
func ColumnSum(parent LazyCharts) LazyCharts {
	return Group(parent, TotalPartition(parent.Titles()))
}

type cache struct {
	chartsNode
	rows map[string]*cacheRow
}

type cacheRow struct {
	channel    chan cacheRowRequest
	begin, end int
	data       []float64
}

type cacheRowRequest struct {
	back       chan cacheRowAnswer
	begin, end int
}

type cacheRowAnswer []float64

// Cache is a LazyCharts that stores data to avoid duplicating work in parent.
// The cache is filled when the data is requested. The data is stored in one
// continuous block per row. Non-requested parts in between are filled.
// E.g. if Data({"A"}, 0, 4) and Column({"A"}, 16) are called, row "A" will store
// range [0, 17).
func Cache(parent LazyCharts) LazyCharts {
	rows := make(map[string]*cacheRow)
	for _, k := range parent.Titles() {
		row := &cacheRow{
			channel: make(chan cacheRowRequest),
			begin:   -1, end: -1,
			data: make([]float64, 0),
		}
		rows[k.Key()] = row

		go func(title Title, row *cacheRow, parent LazyCharts) {
			for request := range row.channel {

				if row.begin > -1 {
					if row.begin <= request.begin && row.end >= request.end {

					} else {

						if request.begin < row.begin {
							row.data = append(
								parent.Data([]Title{title}, request.begin, row.begin)[0],
								row.data...)

							row.begin = request.begin
						}
						if row.end < request.end {
							row.data = append(
								row.data,
								parent.Data([]Title{title}, row.end, request.end)[0]...)

							row.end = request.end
						}

						answer := row.data[request.begin-row.begin : request.end-row.begin]
						request.back <- answer
						continue
					}
				}

				row.data = parent.Data([]Title{title}, request.begin, request.end)[0]
				row.begin, row.end = request.begin, request.end

				request.back <- row.data
			}
		}(k, row, parent)
	}

	return &cache{
		chartsNode: chartsNode{parent},
		rows:       rows,
	}
}

func (c *cache) row(title Title, begin, end int) []float64 {
	row := c.rows[title.Key()]

	back := make(chan cacheRowAnswer)

	row.channel <- cacheRowRequest{back, begin, end}
	answer := <-back
	close(back)
	return answer
}

func (c *cache) Column(titles []Title, index int) []float64 {
	data := c.Data(titles, index, index+1)
	tvm := make([]float64, len(titles))
	for i := range titles {
		tvm[i] = data[i][0]
	}
	return tvm
}

func (c *cache) Data(titles []Title, begin, end int) [][]float64 {
	data := make([][]float64, len(titles))
	back := make(chan indexLine)

	for k := range titles {
		go func(k int) {
			back <- indexLine{
				i:  k,
				vs: c.row(titles[k], begin, end),
			}
		}(k)
	}

	for range titles {
		tl := <-back
		data[tl.i] = tl.vs
	}
	return data
}

type only struct {
	chartsNode
	titles []Title
}

// Only keeps only a subset of titles from the parent
func Only(parent LazyCharts, titles []Title) LazyCharts {
	return &only{
		chartsNode: chartsNode{parent: parent},
		titles:     titles,
	}
}

func (c *only) Titles() []Title {
	return c.titles
}

func (c *only) Column(titles []Title, index int) []float64 {
	return c.parent.Column(titles, index)
}

func (c *only) Data(titles []Title, begin, end int) [][]float64 {
	return c.parent.Data(titles, begin, end)
}

func Top(c LazyCharts, n int) []Title {
	fullTitles := c.Titles()
	col := c.Column(fullTitles, c.Len()-1)
	m := n + 1
	if len(col) < n {
		m = len(col)
	}

	vs := make([]float64, m)
	ts := make([]Title, m)
	i := 0
	for k, tv := range col {
		vs[i] = tv
		ts[i] = fullTitles[k]
		for j := i; j > 0; j-- {
			if vs[j-1] < vs[j] {
				vs[j-1], vs[j] = vs[j], vs[j-1]
				ts[j-1], ts[j] = ts[j], ts[j-1]
			} else {
				break
			}
		}
		if i+1 < m {
			i++
		}
	}
	if len(ts) > n {
		ts = ts[:n]
	}

	return ts
}

func Id(parent LazyCharts) LazyCharts {
	return parent
}
