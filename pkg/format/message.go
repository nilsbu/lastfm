package format

import (
	"fmt"
	"io"
	"strings"
)

type Message struct {
	Msg string
}

func (f *Message) CSV(w io.Writer, decimal string) error {
	if f.Msg == "" {
		return nil
	}

	lines := strings.Split(f.Msg, "\n")
	var str string
	for _, line := range lines {
		str += fmt.Sprintf("\"%v\";\n", line)
	}

	_, err := io.WriteString(w, str)
	return err
}

func (f *Message) Plain(w io.Writer) error {
	if f.Msg == "" {
		return nil
	}

	str := fmt.Sprintf("%v\n", f.Msg)

	_, err := io.WriteString(w, str)
	return err
}
