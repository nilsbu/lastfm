package charts

import (
	"time"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Interval is a time span. Begin is the first moment, Before is the moment
// immediately after the interval ends.
type Interval struct {
	Begin  rsrc.Day
	Before rsrc.Day
}

type Intervals interface {
	At(index int) (interval Interval)
	Index(day rsrc.Day) (index int)
	Len() (len int)
}

type intervalsBase struct {
	begin rsrc.Day
	n     int
	step  int
}

func (i intervalsBase) Len() (len int) {
	return i.n
}

type dayIntervals struct {
	intervalsBase
}

func (i dayIntervals) At(index int) (interval Interval) {
	return Interval{
		Begin:  rsrc.DayFromTime(i.begin.Time().AddDate(0, 0, index*i.step)),
		Before: rsrc.DayFromTime(i.begin.Time().AddDate(0, 0, (index+1)*i.step)),
	}
}

func (i dayIntervals) Index(day rsrc.Day) (index int) {
	duration := day.Time().Sub(i.begin.Time())
	return int(duration.Hours()) / 24 / i.step
}

func Days(begin, end rsrc.Day) Intervals {
	return MultiDays(begin, end, 1)
}

func MultiDays(begin, end rsrc.Day, step int) Intervals {
	duration := end.Time().Sub(begin.Time())
	n := int(duration.Hours()) / 24 / step

	return dayIntervals{intervalsBase{
		begin: begin,
		n:     n,
		step:  step,
	}}
}

type monthIntervals struct {
	intervalsBase
}

func (i monthIntervals) At(index int) (interval Interval) {
	return Interval{
		Begin:  rsrc.DayFromTime(i.begin.Time().AddDate(0, index*i.step, 0)),
		Before: rsrc.DayFromTime(i.begin.Time().AddDate(0, (index+1)*i.step, 0)),
	}
}

func (i monthIntervals) Index(day rsrc.Day) (index int) {
	bt := i.begin.Time()
	by := bt.Year()
	bm := int(bt.Month())

	qt := day.Time()
	qy := qt.Year()
	qm := int(qt.Month())

	return ((qm - bm) + 12*(qy-by)) / i.step
}

func Months(begin, end rsrc.Day, step int) Intervals {
	t := begin.Time()
	y := t.Year()
	m := int(int(t.Month())-1)/step*step + 1
	begin = rsrc.DayFromTime(time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC))

	bt := begin.Time()
	by := bt.Year()
	bm := int(bt.Month())

	qt := end.Time().AddDate(0, 0, -1)
	qy := qt.Year()
	qm := int(qt.Month())

	n := ((qm - bm) + 12*(qy-by)) / step

	return monthIntervals{intervalsBase{
		begin: begin,
		n:     n + 1,
		step:  step,
	}}
}

type yearIntervals struct {
	intervalsBase
}

func (i yearIntervals) At(index int) (interval Interval) {
	return Interval{
		Begin:  rsrc.DayFromTime(i.begin.Time().AddDate(index*i.step, 0, 0)),
		Before: rsrc.DayFromTime(i.begin.Time().AddDate((index+1)*i.step, 0, 0)),
	}
}

func (i yearIntervals) Index(day rsrc.Day) (index int) {
	by := i.begin.Time().Year()
	qy := day.Time().Year()

	return (qy - by) / i.step
}

func Years(begin, end rsrc.Day, step int) Intervals {
	y := begin.Time().Year() / step * step
	begin = rsrc.DayFromTime(time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC))

	n := (end.Time().AddDate(0, 0, -1).Year() - y) / step

	return yearIntervals{intervalsBase{
		begin: begin,
		n:     n + 1,
		step:  step,
	}}
}
