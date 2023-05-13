package charts

import "sort"

type DiffCharts interface {
	// Data(titles []Title, begin, end int) ([][]float64, error)

	// Titles() []Title
	// Len() int

	Charts

	Previous(title Title) (place int, value float64, err error)
}

type diffCharts struct {
	c Charts

	prevIdx int
	prev    []struct {
		title Title
		value float64
	}
}

func NewDiffCharts(c Charts, prevIdx int) DiffCharts {
	return &diffCharts{
		c:       c,
		prevIdx: prevIdx,
	}
}

func (l *diffCharts) Data(titles []Title, begin, end int) ([][]float64, error) {
	return l.c.Data(titles, begin, end)
}

func (l *diffCharts) Titles() []Title {
	return l.c.Titles()
}

func (l *diffCharts) Len() int {
	return l.c.Len()
}

func (l *diffCharts) Previous(title Title) (place int, value float64, err error) {
	titles := l.Titles()
	if l.prev == nil {
		l.prev = make([]struct {
			title Title
			value float64
		}, len(titles))
		if prevData, err := l.c.Data(titles, l.prevIdx, l.prevIdx+1); err != nil {
			return -1, 0, err
		} else {
			for i, prev := range prevData {
				l.prev[i].title = titles[i]
				l.prev[i].value = prev[0]
			}
			sort.Slice(l.prev, func(i, j int) bool {
				return l.prev[i].value > l.prev[j].value
			})
		}
	}

	for i, prev := range l.prev {
		if prev.title.Key() == title.Key() {
			return i, prev.value, nil
		}
	}

	return -1, 0, nil
}
