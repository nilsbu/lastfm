package charts2

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Range describes a range of dates. Begin is the first date of the range, End
// is the first date after the range and Redistered refers to the date the
// charts the Range will be used in was registered.
type Range struct {
	Begin, End, Registered rsrc.Day
}

type interval struct {
	chartsNode
	Range Range
}

// Interval crops the charts to a certain Range.
func Interval(parent LazyCharts, r Range) LazyCharts {
	return interval{
		chartsNode: chartsNode{parent: parent},
		Range:      r,
	}
}

func (c interval) Len() int {
	return int(c.Range.End.Time().Sub(c.Range.Begin.Time()).Hours()) / 24
}

func (c interval) Row(title Title, begin, end int) []float64 {
	return c.Data([]Title{title}, begin, end)[title.Key()].Line
}

func (c interval) Column(titles []Title, index int) TitleValueMap {
	data := c.Data(titles, index, index+1)
	tvm := make(TitleValueMap)
	for title, line := range data {
		tvm[title] = TitleValue{
			Title: line.Title,
			Value: line.Line[0],
		}
	}
	return tvm
}

func (c interval) Data(titles []Title, begin, end int) TitleLineMap {
	data := make(TitleLineMap)
	back := make(chan TitleLine)

	boffset := int(c.Range.Begin.Time().Sub(c.Range.Registered.Time()).Hours()) / 24

	for k := range titles {
		go func(k int) {
			back <- TitleLine{
				Title: titles[k],
				Line:  c.parent.Row(titles[k], begin+boffset, end+boffset),
			}
		}(k)
	}

	for range titles {
		tl := <-back
		data[tl.Title.Key()] = tl
	}
	return data
}
