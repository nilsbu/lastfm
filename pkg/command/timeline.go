package command

import (
	"fmt"
	"time"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/timeline"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type printTimeline struct {
	from   time.Time
	before time.Time
}

func (cmd printTimeline) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {

	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	cha := charts.Compile(plays)

	replace, err := unpack.LoadArtistCorrections(session.User, s)
	if err == nil {
		cha = cha.Correct(replace)
	}

	user, err := unpack.LoadUserInfo(session.User, s)
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	events := timeline.CompileEvents(
		cha,
		user.Registered, rsrc.Date(cmd.from), rsrc.Date(cmd.before))

	for _, event := range events {
		t := event.Date.Time()
		d.Display(&format.Message{
			Msg: fmt.Sprintf(
				"%v: %v",
				fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()),
				event.Message,
			)})
	}

	return nil
}
