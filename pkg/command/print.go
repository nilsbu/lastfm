package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type printCharts struct {
	keys       string
	by         string
	name       string
	percentage bool
	normalized bool
	duration   bool
	entry      float64
	n          int
}

func (cmd printCharts) getSteps() ([]string, error) {
	steps := []string{cmd.keys + "s"}
	if cmd.duration {
		steps[0] += "duration"
	}
	if cmd.normalized {
		steps = append(steps, "gaussian", "cache", "normalize")
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
	date rsrc.Day
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {

	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "sum", "cache")

	if cmd.date != nil {
		steps = append(steps, fmt.Sprintf("day,%v", cmd.date))
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
	f := &format.DiffCharts{
		Charts:     []charts.DiffCharts{charts.NewDiffCharts(cha, cha.Len()-7)},
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printFade struct {
	printCharts
	hl   float64
	date rsrc.Day
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("fade,%v", cmd.hl), "cache")

	if cmd.date != nil {
		steps = append(steps, fmt.Sprintf("day,%v", cmd.date))
	}
	steps = append(steps, fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 2

	f := &format.DiffCharts{
		Charts:     []charts.DiffCharts{charts.NewDiffCharts(cha, cha.Len()-7)},
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
	f := &format.DiffCharts{
		Charts:     []charts.DiffCharts{charts.NewDiffCharts(cha, cha.Len()-7)},
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printInterval struct {
	printCharts
	begin rsrc.Day
	end   rsrc.Day
}

func (cmd printInterval) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("interval,%v,%v", cmd.begin, cmd.end), "sum", "cache")
	steps = append(steps, fmt.Sprintf("top,%v", cmd.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.DiffCharts{
		Charts:     []charts.DiffCharts{charts.NewDiffCharts(cha, cha.Len()-7)},
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

type printAfter struct {
	printCharts
	n int
}

func (cmd printAfter) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "sum", "cache", "offset")
	steps = append(steps, fmt.Sprintf("column,%d", cmd.n), fmt.Sprintf("top,%v", cmd.printCharts.n))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	prec := 2

	f := &format.DiffCharts{
		Charts:     []charts.DiffCharts{charts.NewDiffCharts(cha, cha.Len()-7)},
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

type printPeriods struct {
	printCharts
	period     string
	begin, end rsrc.Day
}

func (cmd printPeriods) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, "id")
	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}
	ll := cha.Len()

	interval, err := charts.CroppedRange(cmd.begin, cmd.end, pl.Registered(), ll)
	if err != nil {
		return err
	}
	steps = append(steps, fmt.Sprintf("interval,%v,%v", cmd.begin, cmd.end))

	ranges, _ := charts.ParseRanges(cmd.period, interval.Begin, rsrc.Between(cmd.begin, cmd.end).Days())

	steps = append(steps, fmt.Sprintf("periods,%v", cmd.period), "cache")

	cha, err = pl.Execute(steps)
	if err != nil {
		return err
	}

	chas := make([]charts.Charts, cha.Len())
	for i := range chas {
		args := []string{}
		args = append(args, steps...)
		args = append(args, fmt.Sprintf("column,%v", i), fmt.Sprintf("top,%v", cmd.n))

		chas[i], err = pl.Execute(args)
		if err != nil {
			return err
		}
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     chas,
		Ranges:     ranges,
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printFades struct {
	printCharts
	hl         float64
	period     string
	begin, end rsrc.Day
}

func (cmd printFades) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	steps, err := cmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps,
		fmt.Sprintf("fade,%v", cmd.hl),
		"cache")

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	interval, err := charts.CroppedRange(cmd.begin, cmd.end, pl.Registered(), cha.Len())
	if err != nil {
		return err
	}
	steps = append(steps, fmt.Sprintf("interval,%v,%v", cmd.begin, cmd.end))

	ranges, _ := charts.ParseRanges(cmd.period, interval.Begin, rsrc.Between(cmd.begin, cmd.end).Days())
	steps = append(steps, fmt.Sprintf("step,%v", cmd.period))
	cha, err = pl.Execute(steps)
	if err != nil {
		return err
	}

	chas := make([]charts.Charts, cha.Len())
	for i := range chas {
		args := []string{}
		args = append(args, steps...)
		args = append(args, fmt.Sprintf("column,%v", i), fmt.Sprintf("top,%v", cmd.n))

		chas[i], err = pl.Execute(args)
		if err != nil {
			return err
		}
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     chas,
		Ranges:     ranges,
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printRaw struct {
	precision int
	steps     []string
}

func (cmd printRaw) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	cha, err := pl.Execute(cmd.steps)
	if err != nil {
		return err
	}

	f := &format.Charts{
		Charts:     []charts.Charts{cha},
		Numbered:   true,
		Precision:  cmd.precision,
		Percentage: false,
	}

	return d.Display(f)
}
