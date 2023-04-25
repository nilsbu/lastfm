package display

import "github.com/nilsbu/lastfm/pkg/format"

type null struct {
}

func NewNull() Display {
	return &null{}
}

func (n *null) Display(f format.Formatter) error {
	return nil
}
