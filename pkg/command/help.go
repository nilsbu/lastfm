package command

import "github.com/nilsbu/lastfm/io"

type command interface {
	Execute(ioPool io.Pool) error
}

type help struct{}

func (help) Execute(ioPool io.Pool) error {
	// TODO fill
	return nil
}
