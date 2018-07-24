package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/config"
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type printCharts struct {
	by         string
	name       string
	percentage bool
	normalized bool
	n          int
}

func (cmd printCharts) getOutCharts(
	user string,
	cha charts.Charts,
	f func(charts.Charts) charts.Charts,
	r rsrc.Reader) (charts.Charts, error) {
	replace, err := unpack.LoadArtistCorrections(user, r)
	if err != nil {
		return nil, err
	}
	cha = cha.Correct(replace)

	if cmd.normalized {
		nm := charts.GaussianNormalizer{
			Sigma:       30,
			MirrorFront: true,
			MirrorBack:  true}
		cha = nm.Normalize(cha)
	}

	cha = f(cha)

	if cmd.name == "" {
		switch cmd.by {
		case "all":
			return cha, nil
		case "super":
			supertags, err := getSupertags(cha, user, r)
			if err != nil {
				return nil, err
			}

			return cha.Group(supertags), nil
		default:
			return nil, fmt.Errorf("chart type '%v' not supported", cmd.by)
		}
	} else {
		var container map[string]charts.Charts
		switch cmd.by {
		case "all":
			return nil, errors.New("name must be empty for chart type 'all'")
		case "super":
			supertags, err := getSupertags(cha, user, r)
			if err != nil {
				return nil, err
			}

			container = cha.Split(supertags)
		default:
			return nil, fmt.Errorf("chart type '%v' not supported", cmd.by)
		}

		out, ok := container[cmd.name]
		if !ok {
			return nil, fmt.Errorf("name '%v' not found", cmd.name)
		}

		return out, nil
	}
}

type printTotal struct {
	printCharts
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	cha := charts.Compile(plays)

	out, err := cmd.printCharts.getOutCharts(
		session.User,
		cha,
		func(c charts.Charts) charts.Charts { return c.Sum() },
		s)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     out,
		Column:     -1,
		Count:      cmd.n,
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}

type printFade struct {
	printCharts
	hl float64
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	cha := charts.Compile(plays)

	out, err := cmd.printCharts.getOutCharts(
		session.User,
		cha,
		func(c charts.Charts) charts.Charts { return c.Fade(cmd.hl) },
		s)
	if err != nil {
		return err
	}

	f := &format.Charts{
		Charts:     out,
		Column:     -1,
		Count:      cmd.n,
		Numbered:   true,
		Precision:  2,
		Percentage: cmd.percentage,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}

func getSupertags(
	c charts.Charts,
	user string,
	r rsrc.Reader,
) (charts.Partition, error) {

	tags, err := organize.LoadArtistTags(c.Keys(), r)
	if err != nil {
		return nil, err
	}

	corrections, _ := unpack.LoadSupertagCorrections(user, r)

	return charts.Supertags(tags, config.Supertags, corrections), nil
}

type printPeriod struct {
	printCharts
	period string
}

func (cmd printPeriod) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	cha := charts.Compile(plays)

	out, err := cmd.printCharts.getOutCharts(
		session.User,
		cha,
		func(c charts.Charts) charts.Charts { return c.Sum() },
		s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, s)
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	period, err := charts.Period(cmd.period)
	if err != nil {
		return err
	}

	col := out.Interval(period, user.Registered)
	sumTotal := col.Sum()
	col = col.Top(cmd.n)

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Column{
		Column:     col,
		Numbered:   true,
		Percentage: cmd.percentage,
		Precision:  prec,
		SumTotal:   sumTotal,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}

type printTags struct {
	artist string
}

func (cmd printTags) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {

	tags, err := unpack.LoadArtistTags(cmd.artist, s)
	if err != nil {
		return err
	}

	col := make(charts.Column, len(tags))
	for i, tag := range tags {
		col[i] = charts.Score{Name: tag.Name, Score: float64(tag.Count)}
	}

	f := &format.Column{
		Column:   col,
		Numbered: true}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}
