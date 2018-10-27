package format

import "io"

type Error struct {
	Err error
}

func (f *Error) Plain(w io.Writer) error {
	if f.Err == nil {
		return nil
	}

	if _, err := io.WriteString(w, f.Err.Error()); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "\n"); err != nil {
		return err
	}
	return nil
}
