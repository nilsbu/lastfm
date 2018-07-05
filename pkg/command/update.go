package command

import (
	"time"

	"github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type updateHistory struct{}

func (cmd updateHistory) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	user, err := unpack.LoadUserInfo(session.User, s)
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	plays, err := organize.UpdateAllDayPlays(user, rsrc.Date(time.Now()), s)
	if err != nil {
		return errors.Wrap(err, "failed to update user history")
	}

	err = organize.WriteAllDayPlays(plays, user.Name, s)
	if err != nil {
		return errors.Wrap(err, "failed to write alldayplays")
	}

	return nil
}
