package format

import (
	"bytes"
	"testing"
)

func TestMessageCSV(t *testing.T) {
	cases := []struct {
		msg string
		str string
	}{
		{
			"", "",
		},
		{
			"some text\nnew line",
			"\"some text\"\n\"new line\"\n",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Message{Msg: c.msg}
			formatter.CSV(buf, ".")

			msg := buf.String()
			if c.msg == "" {
				if msg != "" {
					t.Error("something was printed despite empty message")
				}
			} else if msg != c.str {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", msg, c.str)
			}
		})
	}
}

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

func TestMessageHTML(t *testing.T) {
	cases := []struct {
		msg       string
		formatted string
	}{
		{
			"",
			"",
		},
		{
			"some text\nnew line",
			"some text<br/>new line<br/>",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Message{Msg: c.msg}
			formatter.HTML(buf)

			msg := buf.String()
			if c.formatted == "" {
				if msg != "" {
					t.Error("something was printed despite empty message")
				}
			} else if msg != c.formatted {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", msg, c.formatted)
			}
		})
	}
}
