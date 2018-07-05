package command

import (
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type command interface {
	Execute(s store.Store, d display.Display) error
}

// Execute executes the command described in the arguments.
func Execute(
	args []string,
	session *unpack.SessionInfo,
	s store.Store,
	d display.Display) error {
	cmd, err := resolve(args, session)
	if err != nil {
		return err
	}

	return cmd.Execute(s, d)
}
