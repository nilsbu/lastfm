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
	for _, line := range lines {
		fmt.Fprintf(w, "\"%v\"\n", line)
	}

	return nil
}

func (f *Message) Plain(w io.Writer) error {
	if f.Msg == "" {
		return nil
	}

	fmt.Fprintf(w, "%v", f.Msg)
	io.WriteString(w, "\n")

	return nil
}

func (f *Message) HTML(w io.Writer) error {
	if f.Msg == "" {
		return nil
	}

	var html string = strings.ReplaceAll(f.Msg, "\n", "<br/>")
	_, err := fmt.Fprintf(w, "%v", html)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "<br/>")
	return err
}

func (f *Message) JSON(w io.Writer) error {
	if f.Msg == "" {
		return nil
	}

	_, err := fmt.Fprintf(w, `{"msg": "%v}`, f.Msg)
	return err
}
