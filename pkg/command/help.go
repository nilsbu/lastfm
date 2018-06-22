package command

import "github.com/nilsbu/lastfm/pkg/io"

type command interface {
	Execute(store io.Store) error
}

type help struct{}

func (help) Execute(store io.Store) error {
	// TODO fill
	return nil
}
