package command

import (
	"encoding/json"
	"time"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type updateHistory struct {
	sid organize.SessionID
}

func (c updateHistory) Execute(s store.Store, d display.Display) error {
	// TODO turn this into a function
	name, err := rsrc.UserInfo(string(c.sid))
	if err != nil {
		return err
	}

	data, err := s.Read(name)
	if err != nil {
		return err
	}

	userRaw := unpack.UserInfo{}
	err = json.Unmarshal(data, &userRaw)
	if err != nil {
		return err
	}

	user := unpack.GetUser(&userRaw)

	plays, err := organize.UpdateAllDayPlays(user, rsrc.Date(time.Now()), s)
	if err != nil {
		return err
	}

	err = organize.WriteAllDayPlays(plays, user.Name, s)
	if err != nil {
		return err
	}

	return nil
}
