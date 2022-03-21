package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type tableTotal struct {
	printCharts
	step int
}

func (cmd tableTotal) Accumulate(c charts.Charts) charts.Charts {
	return charts.Sum(c)
}

func (cmd tableTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "sum")
	steps = append(steps, fmt.Sprintf("top %v", cmd.n))

	w := newWeb(session, s)
	cha, err := w.Execute(steps)
	if err != nil {
		return err
	}

	ranges, _ := charts.ParseRanges(fmt.Sprintf("%vd", cmd.step), w.Registered(), cha.Len())

	f := &format.Table{
		Charts: cha,
		Ranges: ranges,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}

type tableFade struct {
	printCharts
	step int
	hl   float64
}

func (cmd tableFade) Accumulate(c charts.Charts) charts.Charts {
	return charts.Fade(c, cmd.hl)
}

// TODO Test table fade
func (cmd tableFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps,
		fmt.Sprintf("fade %v", cmd.hl),
		fmt.Sprintf("top %v", cmd.n),
		fmt.Sprintf("step %vd", cmd.step))

	w := newWeb(session, s)
	cha, err := w.Execute(steps)
	if err != nil {
		return err
	}
	ranges, _ := charts.ParseRanges(fmt.Sprintf("%vd", cmd.step), w.Registered(), cha.Len())

	f := &format.Table{
		Charts: cha,
		Ranges: ranges,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}

type tablePeriods struct {
	printCharts
	period string
}

func (cmd tablePeriods) Accumulate(c charts.Charts) charts.Charts {
	return charts.Sum(c)
}

func (cmd tablePeriods) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "id")
	steps = append(steps, fmt.Sprintf("periods %v", cmd.period), fmt.Sprintf("top %v", cmd.n))

	w := newWeb(session, s)
	cha, err := w.Execute(steps)
	if err != nil {
		return err
	}
	ranges, _ := charts.ParseRanges(cmd.period, w.Registered(), cha.Len())

	f := &format.Table{
		Charts: cha,
		Ranges: ranges,
	}

	return d.Display(f)
}
