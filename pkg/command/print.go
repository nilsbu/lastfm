package command

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/store"
)

type printTotal struct {
	sid organize.SessionID
	n   int
}

func (cmd printTotal) Execute(s store.Store, d display.Display) error {
	plays, err := organize.ReadAllDayPlays(string(cmd.sid), s)
	if err != nil {
		return err
	}

	sums := charts.Compile(plays).Sum()
	f := &format.Charts{
		Charts:    charts.Charts(sums),
		Column:    -1,
		Count:     cmd.n,
		Numbered:  true,
		Precision: 0,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}
