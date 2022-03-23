package command

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type sessionInfo struct{}

func (cmd sessionInfo) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	if session == nil {
		d.Display(&format.Message{Msg: "no session is running"})
	} else {
		d.Display(&format.Message{
			Msg: fmt.Sprintf("a session is running for user '%v'", session.User)})
		// TODO print params in session info
	}

	return nil
}

type sessionStart struct {
	user string
}

func (cmd sessionStart) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	if session != nil {
		return fmt.Errorf("a session is already running for '%v'", session.User)
	}

	return unpack.WriteSessionInfo(&unpack.SessionInfo{User: cmd.user}, s)
}

type sessionStop struct{}

func (cmd sessionStop) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	if session == nil {
		return errors.New("no session is running")
	}

	return s.Remove(rsrc.SessionInfo())
}

type sessionConfig struct {
	option, value string
}

func (cmd sessionConfig) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {
	if session == nil {
		return errors.New("no session is running")
	}

	found := false
	for _, opt := range storableOptions {
		if opt.name == cmd.option {
			if _, err := parseArgument(cmd.value, opt.kind); err != nil {
				return err
			}
			found = true
		}
	}
	if !found {
		return fmt.Errorf("option '%v' doesn't exist", cmd.option)
	}

	params := make(map[string]string)
	for k, v := range session.Options {
		params[k] = v
	}
	params[cmd.option] = cmd.value

	return unpack.WriteSessionInfo(&unpack.SessionInfo{User: session.User, Options: params}, s)
}
