package command

import (
	"fmt"

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

	intervals, err := parsePeriod(out, user.Registered, cmd.period)
	if len(intervals) == 0 || err != nil {
		return err
	}

	out = out.Intervals(intervals, user.Registered)

	f := &format.Table{
		Charts: out,
		First:  rsrc.ToDay(intervals[0].Begin.Unix()),
		Step:   1,
		Count:  cmd.n,
	}

	return d.Display(f)
}

func parsePeriod(
	cha charts.Charts,
	registered rsrc.Day,
	descr string) ([]charts.Interval, error) {

	if descr == "m" {
		return cha.ToIntervals(charts.Month, registered), nil
	} else if descr == "y" {
		return cha.ToIntervals(charts.Year, registered), nil
	} else {
		return nil, fmt.Errorf("period '%v' is invalid", descr)
	}
}
