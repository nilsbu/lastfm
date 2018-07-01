package command

import (
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/store"
)

type help struct{}

func (help) Execute(s store.Store, d display.Display) error {
	// TODO fill
	return nil
}
