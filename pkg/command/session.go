package command

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type sessionInfo struct {
	sid organize.SessionID
}

func (c sessionInfo) Execute(s store.Store, d display.Display) error {
	if c.sid == "" {
		d.Display(&format.Message{Msg: "no session is running"})
	} else {
		d.Display(&format.Message{
			Msg: fmt.Sprintf("a session is running for user '%v'", c.sid)})
	}

	return nil
}

type sessionStart struct {
	sid  organize.SessionID
	user string
}

func (c sessionStart) Execute(s store.Store, d display.Display) error {
	if c.sid != "" {
		return fmt.Errorf("a session is already running for '%v'", c.sid)
	}

	// TODO create function
	sid := &unpack.SessionID{User: string(c.user)}
	data, err := json.Marshal(sid)
	if err != nil {
		return err
	}

	return s.Write(data, rsrc.SessionID())
}

type sessionStop struct {
	sid organize.SessionID
}

func (c sessionStop) Execute(s store.Store, d display.Display) error {
	if c.sid == "" {
		return errors.New("no session is running")
	}
	// TODO crate function
	return io.FileIO{}.Remove(rsrc.SessionID())
}
