package format

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
)

func TestErrorPlain(t *testing.T) {
	cases := []struct {
		err error
	}{
		{
			nil,
		},
		{
			errors.New("fail"),
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			formatter := &Error{Err: c.err}
			formatter.Plain(buf)

			err := buf.String()
			if err == "" && c.err != nil {
				t.Error("no message was printed but error was not nil")
			}
			if c.err != nil && err != c.err.Error()+"\n" {
				t.Errorf("false formatting:\nhas:\n%v\nwant:\n%v", err, c.err)
			}
		})
	}
}
