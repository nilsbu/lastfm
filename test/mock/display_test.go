package mock

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/format"
)

func TestDisplay(t *testing.T) {
	cases := []struct {
		fs []format.Formatter
	}{
		{
			[]format.Formatter{format.Message{Msg: "A"}},
		},
		{
			[]format.Formatter{
				format.Message{Msg: "A"},
				format.Message{Msg: "LOL"}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			d := NewDisplay()

			for i, f := range c.fs {
				err := d.Display(f)
				if err != nil {
					t.Errorf("unexpected error at input %v: %v", i, err)
				}
			}

			if len(d.Msgs) != len(c.fs) {
				t.Fatalf("saved %v formatters but expected %v", len(d.Msgs), len(c.fs))
			}
			for i, f := range c.fs {
				if !reflect.DeepEqual(d.Msgs[i], f) {
					t.Errorf("got '%v', expected '%v'", d.Msgs[i], f)
				}
			}
		})
	}
}
