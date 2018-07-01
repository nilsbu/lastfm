package format

import "io"

type Error struct {
	Err error
}

func (f *Error) Plain(w io.Writer) error {
	if f.Err == nil {
		return nil
	}

	io.WriteString(w, f.Err.Error())
	io.WriteString(w, "\n")
	return nil
}
