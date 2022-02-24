package display

import (
	"io"

	"github.com/nilsbu/lastfm/pkg/format"
)

type Web struct {
	writer io.Writer
}

func NewWeb(w io.Writer) *Web {
	return &Web{writer: w}
}

func (d *Web) Display(f format.Formatter) error {
	io.WriteString(d.writer, "<html><body>")
	defer io.WriteString(d.writer, "</body></html>")
	return f.HTML(d.writer)
}
