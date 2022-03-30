package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type tableTotal struct {
	printCharts
	step int
}

func (cmd tableTotal) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "sum", "cache")
	steps = append(steps, fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	ranges, _ := charts.ParseRanges(fmt.Sprintf("%vd", cmd.step), pl.Registered(), cha.Len())

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

// TODO Test table fade
func (cmd tableFade) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps,
		fmt.Sprintf("fade,%v", cmd.hl),
		"cache",
		fmt.Sprintf("top,%v", cmd.n),
		fmt.Sprintf("step,%vd", cmd.step))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}
	ranges, _ := charts.ParseRanges(fmt.Sprintf("%vd", cmd.step), pl.Registered(), cha.Len())

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

func (cmd tablePeriods) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "id")
	steps = append(steps, fmt.Sprintf("periods,%v", cmd.period), "cache", fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}
	ranges, _ := charts.ParseRanges(cmd.period, pl.Registered(), cha.Len())

	f := &format.Table{
		Charts: cha,
		Ranges: ranges,
	}

	return d.Display(f)
}
