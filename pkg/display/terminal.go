package display

import (
	"io"
	"os"

	"github.com/nilsbu/lastfm/pkg/format"
)

type Terminal struct {
	Writer io.Writer
}

func NewTerminal() *Terminal {
	return &Terminal{Writer: os.Stdout}
}

func (d Terminal) Display(f format.Formatter) error {
	return f.Plain(d.Writer)
}
