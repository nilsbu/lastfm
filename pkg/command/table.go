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

	w := newWeb(session, s)
	cha, err := w.Execute(steps)
	if err != nil {
		return err
	}

	f := &format.Table{
		Charts: cha,
		First:  w.Registered(),
		Step:   cmd.step,
		Count:  cmd.n,
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

func (cmd tableFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("fade %v", cmd.hl))

	w := newWeb(session, s)
	cha, err := w.Execute(steps)
	if err != nil {
		return err
	}

	f := &format.Table{
		Charts: cha,
		First:  w.Registered(),
		Step:   cmd.step,
		Count:  cmd.n,
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
	steps = append(steps, fmt.Sprintf("periods %v", cmd.period))

	w := newWeb(session, s)
	cha, err := w.Execute(steps)
	if err != nil {
		return err
	}

	f := &format.Table{
		Charts: cha,
		// First:  ranges.Delims[0], // TODO
		Step:  1,
		Count: cmd.n,
	}

	return d.Display(f)
}
