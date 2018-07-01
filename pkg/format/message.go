package format

import "io"

type Message struct {
	Msg string
}

func (f Message) Plain(w io.Writer) error {
	if f.Msg == "" {
		return nil
	}

	io.WriteString(w, f.Msg)
	io.WriteString(w, "\n")
	return nil
}
