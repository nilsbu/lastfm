package format

import "io"

type Message struct {
	Msg string
}

func (f Message) Plain(w io.Writer) {
	if f.Msg == "" {
		return
	}

	io.WriteString(w, f.Msg)
	io.WriteString(w, "\n")
}
