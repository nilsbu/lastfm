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
	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	today := rsrc.DayFromTime(time.Now())
	plays, err := organize.UpdateHistory(user, today, s)
	if err != nil {
		return errors.Wrap(err, "failed to update user history")
	}

	err = unpack.WriteSongHistory(plays, user.Name, s)
	if err != nil {
		return errors.Wrap(err, "failed to write song history")
	}

	err = unpack.WriteBookmark(today, user.Name, s)
	if err != nil {
		return errors.Wrap(err, "failed to write bookmark")
	}

	return nil
}
