package command

import (
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type cache struct {
	port int
}

func (cmd cache) Execute(
	session *unpack.SessionInfo,
	s store.Store,
	d display.Display) error {

	io.RunCacheServer(cmd.port)
	return nil
}
