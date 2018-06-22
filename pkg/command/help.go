package command

import "github.com/nilsbu/lastfm/pkg/store"

type command interface {
	Execute(s store.Store) error
}

type help struct{}

func (help) Execute(s store.Store) error {
	// TODO fill
	return nil
}
