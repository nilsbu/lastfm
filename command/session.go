package command

import (
	"encoding/json"
	"fmt"

	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/unpack"
)

type sessionInfo struct{}

func (sessionInfo) Execute(ioPool io.Pool) error {
	sid, err := organize.LoadSessionID(io.SeqReader(ioPool.ReadFile))
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

func (c sessionStart) Execute(ioPool io.Pool) error {
	sid := &unpack.SessionID{User: string(c.user)}
	data, err := json.Marshal(sid)
	if err != nil {
		return err
	}

	return io.SeqWriter(ioPool.WriteFile).Write(data, rsrc.SessionID())
}

type sessionStop struct{}

func (sessionStop) Execute(ioPool io.Pool) error {
	return io.FileRemover{}.Remove(rsrc.SessionID())
}
