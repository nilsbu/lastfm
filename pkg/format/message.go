package format

import "io"

type Message struct {
	Msg string
}

func (f Message) Plain(w io.Writer) error {
	if f.Msg == "" {
		return nil
	}

	if _, err := io.WriteString(w, f.Msg); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "\n"); err != nil {
		return err
	}
	return nil
}
