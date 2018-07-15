package format

import "io"

type Error struct {
	Err error
}

func (f *Error) Plain(w io.Writer) {
	if f.Err == nil {
		return
	}

	io.WriteString(w, f.Err.Error())
	io.WriteString(w, "\n")
}
