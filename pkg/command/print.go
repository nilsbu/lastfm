package command

import (
	"fmt"
	"time"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type printCharts struct {
	keys       string // defaults to "artist"
	by         string
	name       string
	percentage bool
	normalized bool
	duration   bool
	entry      float64
	n          int
}

func (cmd printCharts) getSteps() ([]string, error) {
	if cmd.keys == "" {
		cmd.keys = "artist"
	}

	steps := []string{cmd.keys + "s"}
	if cmd.duration {
		steps[0] += "duration"
	}
	if cmd.normalized {
		steps = append(steps, "gaussian", "normalize")
	}

	steps = append(steps, "*")

	if cmd.percentage {
		steps = append(steps, "normalize")
	}

	if cmd.by != "all" {
		var s1 string
		if cmd.name == "" {
			s1 = "group," + cmd.by
		} else {
			s1 = "split," + cmd.by + "," + cmd.name
		}
		steps = append(steps, s1)
	} else {
		if cmd.name != "" {
			return nil, fmt.Errorf("cannot use name='%v' with by='all'", cmd.name)
		}
	}

	return steps, nil
}

func setStep(steps []string, sub ...string) []string {
	for i, step := range steps {
		if step == "*" {
			var filled []string
			filled = append(filled, steps[:i]...)
			filled = append(filled, sub...)
			filled = append(filled, steps[i+1:]...)
			return filled
		}
	}
	return steps
}

type printTotal struct {
	printCharts
	date time.Time
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {

	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "sum", "cache")

	null := time.Time{}
	if cmd.date != null {
		steps = append(steps, fmt.Sprintf("day,%v", rsrc.DayFromTime(cmd.date)))
	}
	steps = append(steps, fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     []charts.Charts{cha},
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printFade struct {
	printCharts
	hl   float64
	date time.Time
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("fade,%v", cmd.hl), "cache")

	null := time.Time{}
	if cmd.date != null {
		steps = append(steps, fmt.Sprintf("day,%v", rsrc.DayFromTime(cmd.date)))
	}
	steps = append(steps, fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 2

	f := &format.Charts{
		Charts:     []charts.Charts{cha},
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printPeriod struct {
	printCharts
	period string
}

func (cmd printPeriod) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("period,%v", cmd.period), "sum", "cache")
	steps = append(steps, fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     []charts.Charts{cha},
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printInterval struct {
	printCharts
	begin  time.Time
	before time.Time
}

func (cmd printInterval) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("interval,%v,%v", rsrc.DayFromTime(cmd.begin), rsrc.DayFromTime(cmd.before)), "sum", "cache")
	steps = append(steps, fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     []charts.Charts{cha},
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printFadeMax struct {
	printCharts
	hl float64
}

func (cmd printFadeMax) Accumulate(c charts.Charts) charts.Charts {
	return charts.Fade(c, cmd.hl)
}

func (cmd printFadeMax) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("fade,%v", cmd.hl), "cache")
	steps = append(steps, "max", fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 2

	f := &format.Charts{
		Charts:     []charts.Charts{cha},
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printTags struct {
	artist string
}

func (cmd printTags) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {

	tags, err := unpack.LoadArtistTags(cmd.artist, unpack.NewCacheless(s))
	if err != nil {
		return err
	}

	col := make(map[string][]float64, len(tags))
	for _, tag := range tags {
		col[tag.Name] = []float64{float64(tag.Count)}
	}

	f := &format.Charts{
		Charts:     []charts.Charts{charts.FromMap(col)},
		Numbered:   true,
		Precision:  0,
		Percentage: false,
	}

	return d.Display(f)
}
