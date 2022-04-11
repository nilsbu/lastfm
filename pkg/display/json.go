package display

import (
	"io"

	"github.com/nilsbu/lastfm/pkg/format"
)

type json struct {
	writer io.Writer
}

func NewJSON(w io.Writer) Display {
	return &json{writer: w}
}

func (d *json) Display(f format.Formatter) error {
	return f.JSON(d.writer)
}
