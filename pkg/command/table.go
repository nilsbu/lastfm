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

func (cmd tableTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	cha := charts.Compile(plays)

	out, err := cmd.printCharts.getOutCharts(
		session.User,
		cha,
		func(c charts.Charts) charts.Charts { return c.Sum() }, // TODO
		s)
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
