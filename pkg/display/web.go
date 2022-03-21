package display

import (
	"io"

	"github.com/nilsbu/lastfm/pkg/format"
)

type web struct {
	writer io.Writer
}

func NewWeb(w io.Writer) Display {
	return &web{writer: w}
}

func (d *web) Display(f format.Formatter) error {
	io.WriteString(d.writer, "<html><body>")
	defer io.WriteString(d.writer, "</body></html>")
	return f.HTML(d.writer)
}
