package command

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type tableTotal struct {
	printCharts
	step int
}

func (cmd tableTotal) Accumulate(c charts.Charts) charts.Charts {
	return c.Sum()
}

func (cmd tableTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, s)
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

func (cmd tableFade) Accumulate(c charts.Charts) charts.Charts {
	return c.Fade(cmd.hl)
}

func (cmd tableFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, s)
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

func (cmd tablePeriods) Accumulate(c charts.Charts) charts.Charts {
	return c.Sum()
}

func (cmd tablePeriods) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, s)
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	end := rsrc.ToDay(user.Registered.Midnight() + int64(86400*out.Len()))
	intervals, err := charts.ToIntervals(cmd.period, user.Registered, end)
	if err != nil || intervals.Len() == 0 {
		return err
	}

	out = out.Intervals(intervals)

	f := &format.Table{
		Charts: out,
		First:  intervals.At(0).Begin,
		Step:   1,
		Count:  cmd.n,
	}

	return d.Display(f)
}
