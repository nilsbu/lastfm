package command

import (
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type help struct{}

func (help) Execute(session *unpack.SessionInfo, s store.Store, d display.Display) error {
	// TODO fill
	return nil
}
