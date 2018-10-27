package command

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type sessionInfo struct{}

func (cmd sessionInfo) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	if session == nil {
		d.Display(&format.Message{Msg: "no session is running"})
	} else {
		d.Display(&format.Message{
			Msg: fmt.Sprintf("a session is running for user '%v'", session.User)})
	}

	return nil
}

type sessionStart struct {
	user string
}

func (cmd sessionStart) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	if session != nil {
		return fmt.Errorf("a session is already running for '%v'", session.User)
	}

	return unpack.WriteSessionInfo(&unpack.SessionInfo{User: cmd.user}, s)
}

type sessionStop struct{}

func (cmd sessionStop) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	if session == nil {
		return errors.New("no session is running")
	}

	return s.Remove(rsrc.SessionInfo())
}
