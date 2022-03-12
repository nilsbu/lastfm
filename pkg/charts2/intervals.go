package charts2

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Range describes a range of dates. Begin is the first date of the range, End
// is the first date after the range and Redistered refers to the date the
// charts the Range will be used in was registered.
type Range struct {
	Begin, End, Registered rsrc.Day
}

func ParseRange(str string, registered rsrc.Day, l int) (Range, error) {
	if r, err := parseRange(str, registered); err != nil {
		return r, err
	} else {
		return cropRange(r, l)
	}
}

func parseRange(str string, registered rsrc.Day) (Range, error) {
	if matched, _ := regexp.Match(`^\d{4}(-\d{2}){0,2}$`, []byte(str)); !matched {
		return Range{}, errors.New("pattern is invalid")
	}
	switch len(str) {
	case 4:
		begin := rsrc.ParseDay(str + "-01-01")
		return Range{
			Begin:      begin,
			End:        begin.AddDate(1, 0, 0),
			Registered: registered,
		}, nil
	case 7:
		begin := rsrc.ParseDay(str + "-01")
		return Range{
			Begin:      begin,
			End:        begin.AddDate(0, 1, 0),
			Registered: registered,
		}, nil
	case 10:
		begin := rsrc.ParseDay(str)
		return Range{
			Begin:      begin,
			End:        begin.AddDate(0, 0, 1),
			Registered: registered,
		}, nil
	default:
		return Range{}, errors.New("pattern is invalid")
	}
}

func cropRange(r Range, l int) (Range, error) {
	c := r

	if c.Begin.Midnight() > c.Registered.AddDate(0, 0, l).Midnight() {
		return Range{}, fmt.Errorf("begin (%v) is after end of data (%v)",
			c.Begin, c.Registered.AddDate(0, 0, l))
	}
	if c.End.Midnight() <= c.Registered.Midnight() {
		return Range{}, fmt.Errorf("end (%v) is before or equal to registered (%v)",
			c.End, c.Registered)
	}

	if c.Begin.Midnight() < c.Registered.Midnight() {
		c.Begin = c.Registered
	}
	if c.End.Midnight() > c.Registered.AddDate(0, 0, l).Midnight() {
		c.End = c.Registered.AddDate(0, 0, l)
	}

	return c, nil
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
