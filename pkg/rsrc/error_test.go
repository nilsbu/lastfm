package rsrc

import (
	"errors"
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
)

func TestErrorConstructors(t *testing.T) {
	err := WrapError(fail.Suspicious, errors.New("abc"))

	if err.Sev != fail.Suspicious {
		t.Error("severity must be 'Suspicious'")
	}

	str := err.Err.Error()
	if str != "abc" {
		t.Errorf("wrong error message, was '%v', expected 'abc'", str)
	}
}
