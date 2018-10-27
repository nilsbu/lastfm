package mock

import (
	"github.com/nilsbu/lastfm/pkg/format"
)

// Display is a mock display.Display that stores formatters in the order as they
// are received. It is not thread-safe against parallel calls of Display().
type Display struct {
	Msgs []format.Formatter
}

// NewDisplay constructs a *mock.Distplay.
func NewDisplay() *Display {
	return &Display{}
}

// Display stores the Formatter f and returns nil.
func (d *Display) Display(f format.Formatter) error {
	d.Msgs = append(d.Msgs, f)
	return nil
}
