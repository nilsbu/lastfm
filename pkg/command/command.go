package command

import (
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/store"
)

type command interface {
	Execute(s store.Store, d display.Display) error
}

// Execute executes the command described in the arguments.
func Execute(
	args []string,
	sid organize.SessionID,
	s store.Store,
	d display.Display) error {
	cmd, err := resolve(args, sid)
	if err != nil {
		return err
	}

	return cmd.Execute(s, d)
}
