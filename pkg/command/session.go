package command

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type sessionInfo struct {
	session *unpack.SessionInfo
}

func (cmd sessionInfo) Execute(s store.Store, d display.Display) error {
	if cmd.session == nil {
		d.Display(&format.Message{Msg: "no session is running"})
	} else {
		d.Display(&format.Message{
			Msg: fmt.Sprintf("a session is running for user '%v'", cmd.session.User)})
	}

	return nil
}

type sessionStart struct {
	session *unpack.SessionInfo
	user    string
}

func (cmd sessionStart) Execute(s store.Store, d display.Display) error {
	if cmd.session != nil {
		return fmt.Errorf("a session is already running for '%v'", cmd.session.User)
	}

	return unpack.WriteSessionInfo(&unpack.SessionInfo{User: cmd.user}, s)
}

type sessionStop struct {
	session *unpack.SessionInfo
}

func (cmd sessionStop) Execute(s store.Store, d display.Display) error {
	if cmd.session == nil {
		return errors.New("no session is running")
	}
	// TODO crate function
	return io.FileIO{}.Remove(rsrc.SessionInfo())
}
