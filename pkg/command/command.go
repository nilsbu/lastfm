package command

import (
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type command interface {
	Execute(session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error
}

// Execute executes the command described in the arguments.
func Execute(
	args []string,
	session *unpack.SessionInfo,
	s io.Store,
	pl pipeline.Pipeline,
	d display.Display) error {
	cmd, err := resolve(args, session)
	if err != nil {
		return err
	}

	return cmd.Execute(session, s, pl, d)
}
