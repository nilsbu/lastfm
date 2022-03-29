package charts

import (
	"github.com/nilsbu/async"
)

func fromBeginRange(size, begin, end int) (b, e int) {
	return 0, end
}

type partitionSum struct {
	chartsNode
	partition Partition
}

// Group is a LazyCharts that combines the subsets of the partition from the parent.
func Group(
	parent Charts,
	partition Partition,
) Charts {
	return &partitionSum{
		chartsNode: chartsNode{parent: parent},
		partition:  partition,
	}
}

type titleColumn struct {
	key int
	col []float64
}

func (l *partitionSum) Data(titles []Title, begin, end int) ([][]float64, error) {
	data := make([][]float64, len(titles))

	err := async.Pie(len(titles), func(i int) error {
		line := make([]float64, end-begin)
		for _, key := range l.partition.Titles(titles[i]) {
			res, err := l.parent.Data([]Title{key}, begin, end)
			if err != nil {
				return err
			}
			for j, v := range res[0] {
				line[j] += v
			}
		}
		data[i] = line
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (l *partitionSum) Titles() []Title {
	return l.partition.Partitions()
}

// Subset is a LazyCharts that picks a single subset of the partition from the parent.
func Subset(
	parent Charts,
	partition Partition,
	title Title,
) Charts {
	return Only(parent, partition.Titles(title))
}

// ColumnSum is a LazyCharts that sums up all columns.
// TODO can be optmized by bypassing partitions
func ColumnSum(parent Charts) Charts {
	return Group(parent, TotalPartition(parent.Titles()))
}

type cache struct {
	chartsNode
	rows map[string]*cacheRow
}

type cacheRow struct {
	channel chan cacheRowRequest
	begin   int
	data    []float64
}

type cacheRowRequest struct {
	back       chan cacheRowAnswer
	begin, end int
}

type cacheRowAnswer struct {
	data []float64
	err  error
}

// Cache is a LazyCharts that stores data to avoid duplicating work in parent.
// The cache is filled when the data is requested. The data is stored in one
// continuous block per row. Non-requested parts in between are filled.
// E.g. if Data({"A"}, 0, 4) and Column({"A"}, 16) are called, row "A" will store
// range [0, 17).
func Cache(parent Charts) Charts {
	rows := make(map[string]*cacheRow)
	for _, k := range parent.Titles() {
		row := &cacheRow{
			channel: make(chan cacheRowRequest),
			begin:   -1,
			data:    make([]float64, 0),
		}
		rows[k.Key()] = row

		go func(title Title, row *cacheRow, parent Charts) {
			for request := range row.channel {

				if row.begin > -1 {
					if row.begin <= request.begin && row.begin+len(row.data) >= request.end {
					} else {
						var res [][]float64
						var err error

						if request.begin < row.begin {
							res, err = parent.Data([]Title{title}, request.begin, row.begin)
							newData := []float64{}
							newData = append(newData, res[0]...)
							newData = append(newData, row.data...)
							row.data = newData

							row.begin = request.begin
						}
						if row.begin+len(row.data) < request.end {
							res, err = parent.Data([]Title{title}, row.begin+len(row.data), request.end)
							row.data = append(
								row.data,
								res[0]...)
						}

						request.back <- cacheRowAnswer{
							data: row.data[request.begin-row.begin : request.end-row.begin],
							err:  err,
						}
						continue
					}
				}

				data, err := parent.Data([]Title{title}, request.begin, request.end)
				if err == nil {
					row.data = data[0]
					row.begin = request.begin
				}
				request.back <- cacheRowAnswer{
					data: data[0],
					err:  err,
				}
			}
		}(k, row, parent)
	}

	return &cache{
		chartsNode: chartsNode{parent},
		rows:       rows,
	}
}

func (c *cache) row(title Title, begin, end int) ([]float64, error) {
	row := c.rows[title.Key()]

	back := make(chan cacheRowAnswer)

	row.channel <- cacheRowRequest{back, begin, end}
	answer := <-back
	close(back)
	return answer.data, answer.err
}

func (c *cache) Data(titles []Title, begin, end int) ([][]float64, error) {
	data := make([][]float64, len(titles))

	err := async.Pie(len(titles), func(i int) error {
		row, err := c.row(titles[i], begin, end)
		if err == nil {
			data[i] = row
		}
		return err
	})

	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

type only struct {
	chartsNode
	titles []Title
}

// Only keeps only a subset of titles from the parent
func Only(parent Charts, titles []Title) Charts {
	return &only{
		chartsNode: chartsNode{parent: parent},
		titles:     titles,
	}
}

func (c *only) Titles() []Title {
	return c.titles
}

func (c *only) Data(titles []Title, begin, end int) ([][]float64, error) {
	return c.parent.Data(titles, begin, end)
}

func Top(c Charts, n int) ([]Title, error) {
	fullTitles := c.Titles()
	col, err := c.Data(fullTitles, c.Len()-1, c.Len())
	if err != nil {
		return nil, err
	}
	m := n + 1
	if len(col) < n {
		m = len(col)
	}

	vs := make([]float64, m)
	ts := make([]Title, m)
	i := 0
	for k, tv := range col {
		vs[i] = tv[0]
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

	return ts, nil
}

// Id returns the parent
func Id(parent Charts) Charts {
	return parent
}
