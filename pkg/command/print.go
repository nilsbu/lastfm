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

type printTotal struct {
	by         string
	name       string
	percentage bool
	n          int
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	sums := charts.Compile(plays).Sum()

	replace, err := unpack.LoadArtistCorrections(session.User, s)
	if err != nil {
		return err
	}
	sums = sums.Correct(replace)

	out, err := getOutCharts(session.User, cmd.by, cmd.name, sums, s)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage {
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
	by         string
	name       string
	n          int
	percentage bool
	hl         float64
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	fade := charts.Compile(plays).Fade(cmd.hl)

	replace, err := unpack.LoadArtistCorrections(session.User, s)
	if err != nil {
		return err
	}
	fade = fade.Correct(replace)

	out, err := getOutCharts(session.User, cmd.by, cmd.name, fade, s)
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

func getOutCharts(
	user, by, name string,
	cha charts.Charts,
	r rsrc.Reader) (charts.Charts, error) {
	if name == "" {
		switch by {
		case "all":
			return cha, nil
		case "super":
			supertags, err := getSupertags(cha, user, r)
			if err != nil {
				return nil, err
			}

			return cha.Group(supertags), nil
		default:
			return nil, fmt.Errorf("chart type '%v' not supported", by)
		}
	} else {
		var container map[string]charts.Charts
		switch by {
		case "all":
			return nil, errors.New("name must be empty for chart type 'all'")
		case "super":
			supertags, err := getSupertags(cha, user, r)
			if err != nil {
				return nil, err
			}

			container = cha.Split(supertags)
		default:
			return nil, fmt.Errorf("chart type '%v' not supported", by)
		}

		out, ok := container[name]
		if !ok {
			return nil, fmt.Errorf("name '%v' not found", name)
		}

		return out, nil
	}
}

type printPeriod struct {
	period     string
	by         string
	name       string
	n          int
	percentage bool
}

func (cmd printPeriod) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	sum := charts.Compile(plays).Sum()

	replace, err := unpack.LoadArtistCorrections(session.User, s)
	if err != nil {
		return err
	}
	sum = sum.Correct(replace)

	out, err := getOutCharts(session.User, cmd.by, cmd.name, sum, s)
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
	if cmd.percentage {
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
