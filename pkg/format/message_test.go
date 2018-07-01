package format

import (
	"bytes"
	"testing"
)

func TestMessagePlain(t *testing.T) {
	cases := []struct {
		msg string
	}{
		{
			"",
		},
		{
			"some text\nnew line",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Message{Msg: c.msg}
			formatter.Plain(buf)

			msg := buf.String()
			if c.msg == "" {
				if msg != "" {
					t.Error("something was printed despite empty message")
				}
			} else if msg != c.msg+"\n" {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", msg, c.msg)
			}
		})
	}
}
