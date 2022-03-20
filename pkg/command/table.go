package command

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type tableTotal struct {
	printCharts
	step int
}

func (cmd tableTotal) Accumulate(c charts.LazyCharts) charts.LazyCharts {
	return charts.Sum(c)
}

func (cmd tableTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	f := &format.Table{
		Charts: out,
		First:  user.Registered,
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

func (cmd tableFade) Accumulate(c charts.LazyCharts) charts.LazyCharts {
	return charts.Fade(c, cmd.hl)
}

func (cmd tableFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	f := &format.Table{
		Charts: out,
		First:  user.Registered,
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

func (cmd tablePeriods) Accumulate(c charts.LazyCharts) charts.LazyCharts {
	return charts.Sum(c)
}

func (cmd tablePeriods) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	ranges, err := charts.ParseRanges(cmd.period, user.Registered, out.Len())
	if err != nil {
		return errors.Wrap(err, "failed to parse interval")
	}
	out = charts.Intervals(out, ranges, charts.Id)

	f := &format.Table{
		Charts: out,
		First:  ranges.Delims[0],
		Step:   1,
		Count:  cmd.n,
	}

	return d.Display(f)
}
