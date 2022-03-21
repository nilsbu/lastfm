package charts

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

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

func CroppedRange(begin, end, registered rsrc.Day, l int) (Range, error) {
	return cropRange(Range{
		Begin:      begin,
		End:        end,
		Registered: registered,
	}, l)
}

type interval struct {
	chartsNode
	begin, end int
}

// Interval crops the charts to a certain Range.
func Interval(parent Charts, r Range) Charts {
	begin := int(r.Begin.Time().Sub(r.Registered.Time()).Hours()) / 24
	end := int(r.End.Time().Sub(r.Registered.Time()).Hours()) / 24

	return Crop(parent, begin, end)
}

func Crop(parent Charts, begin, end int) Charts {
	return interval{
		chartsNode: chartsNode{parent: parent},
		begin:      begin,
		end:        end,
	}
}

func Column(parent Charts, col int) Charts {
	if col < 0 {
		col += parent.Len()
	}
	return interval{
		chartsNode: chartsNode{parent: parent},
		begin:      col,
		end:        col + 1,
	}
}

func (c interval) Len() int {
	return c.end - c.begin
}

func (c interval) Column(titles []Title, index int) []float64 {
	data := c.Data(titles, index, index+1)
	tvm := make([]float64, len(titles))
	for i := range titles {
		tvm[i] = data[i][0]
	}
	return tvm
}

func (c interval) Data(titles []Title, begin, end int) [][]float64 {
	data := make([][]float64, len(titles))
	back := make(chan indexLine)

	for k := range titles {
		go func(k int) {
			back <- indexLine{
				i:  k,
				vs: c.parent.Data([]Title{titles[k]}, begin+c.begin, end+c.begin)[0],
			}
		}(k)
	}

	for range titles {
		tl := <-back
		data[tl.i] = tl.vs
	}
	return data
}

type Ranges struct {
	Delims     []rsrc.Day
	Registered rsrc.Day
}

func ParseRanges(
	descr string, registered rsrc.Day, l int,
) (Ranges, error) {

	re := regexp.MustCompile(`^\d*[yMd]$`)
	if !re.MatchString(descr) {
		return Ranges{}, fmt.Errorf("ranges descriptor '%v' invalid", descr)
	}

	dates := []rsrc.Day{registered}

	t := registered.Time()
	y, M := t.Year(), t.Month()
	var date rsrc.Day
	k := descr[len(descr)-1]
	switch k {
	case 'y':
		date = rsrc.DayFromTime(time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC))
	case 'M':
		date = rsrc.DayFromTime(time.Date(y, M, 1, 0, 0, 0, 0, time.UTC))
	default:
		date = registered
	}

	n, err := strconv.Atoi(descr[:len(descr)-1])
	if err != nil {
		n = 1
	}

	end := registered.AddDate(0, 0, l)
	for {
		switch k {
		case 'y':
			date = date.AddDate(n, 0, 0)
		case 'M':
			date = date.AddDate(0, n, 0)
		default:
			date = date.AddDate(0, 0, n)
		}
		if date.Midnight() >= end.Midnight() {
			break
		}
		dates = append(dates, date)
	}
	dates = append(dates, end)

	return Ranges{
		Delims:     dates,
		Registered: registered,
	}, nil
}

type intervals struct {
	chartsNode
	delims []int
	f      func(Charts) Charts
}

// Intervals
func Intervals(parent Charts, rs Ranges, f func(Charts) Charts) Charts {
	delims := make([]int, len(rs.Delims))
	for i, r := range rs.Delims {
		delims[i] = int(r.Time().Sub(rs.Registered.Time()).Hours()) / 24
	}

	return crops(parent, delims, f)
}

func crops(parent Charts, delims []int, f func(Charts) Charts) Charts {
	return intervals{
		chartsNode: chartsNode{parent: parent},
		delims:     delims,
		f:          f,
	}
}

func (c intervals) Len() int {
	return len(c.delims) - 1
}

func (c intervals) Column(titles []Title, index int) []float64 {
	data := c.Data(titles, index, index+1)
	tvm := make([]float64, len(titles))
	for i := range titles {
		tvm[i] = data[i][0]
	}
	return tvm
}

func (c intervals) Data(titles []Title, begin, end int) [][]float64 {
	// TODO speedup
	// data := c.parent.Data(titles, c.delims[begin], c.delims[end])

	lines := make([][]float64, len(titles))
	for j := range titles {
		lines[j] = make([]float64, end-begin)
	}

	// cha := make([]LazyCharts, end-begin)
	for i := begin; i < end; i++ {
		cha := c.f(Crop(c.parent, c.delims[i], c.delims[i+1]))
		cdata := cha.Column(titles, cha.Len()-1)
		for j := range titles {
			lines[j][i-begin] = cdata[j]
		}
	}

	return lines
}
