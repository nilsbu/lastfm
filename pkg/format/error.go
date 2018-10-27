package format

import (
	"fmt"
	"io"
)

type Error struct {
	Err error
}

func (f *Error) CSV(w io.Writer, decimal string) error {
	if f.Err == nil {
		return nil
	}

	str := fmt.Sprintf("\"%v\";\n", f.Err.Error())

	_, err := io.WriteString(w, str)
	return err
}

func (f *Error) Plain(w io.Writer) error {
	if f.Err == nil {
		return nil
	}

	str := fmt.Sprintf("%v\n", f.Err.Error())

	_, err := io.WriteString(w, str)
	return err
}
