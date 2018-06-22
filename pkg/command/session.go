package command

import (
	"encoding/json"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type sessionInfo struct{}

func (sessionInfo) Execute(s store.Store) error {
	sid, err := organize.LoadSessionID(s)
	if err != nil {
		fmt.Println("No session is runnung")
		// TODO should check the kind of error, only some mean there is no session
		return err
	}

	fmt.Printf("A session is running for user '%v'\n", sid)

	return nil
}

type sessionStart struct {
	user rsrc.Name
}

func (c sessionStart) Execute(s store.Store) error {
	sid := &unpack.SessionID{User: string(c.user)}
	data, err := json.Marshal(sid)
	if err != nil {
		return err
	}

	return s.Write(data, rsrc.SessionID())
}

type sessionStop struct{}

func (sessionStop) Execute(s store.Store) error {
	return io.FileIO{}.Remove(rsrc.SessionID())
}
