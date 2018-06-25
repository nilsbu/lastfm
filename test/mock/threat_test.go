package mock

import (
	"errors"
	"strings"
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
)

func TestIsThreatCorrect(t *testing.T) {
	generic := errors.New("A")
	control := &fail.AssessedError{Sev: fail.Control, Err: errors.New("B")}
	suspicious := &fail.AssessedError{Sev: fail.Suspicious, Err: errors.New("C")}
	critical := &fail.AssessedError{Sev: fail.Critical, Err: errors.New("D")}

	cases := []struct {
		err     error
		ok      bool
		sev     fail.Severity
		message string
		correct bool
	}{
		{generic, false, fail.Control, "error must implement fail.Threat but does not", false},
		{control, false, fail.Control, "", true},
		{suspicious, true, fail.Suspicious, "unexpected error", false},
		{control, false, fail.Suspicious, "severity must be 'suspicious' but was 'control'", false},
		{critical, false, fail.Control, "severity must be 'control' but was 'critical'", false},
		{nil, true, fail.Control, "", true},
		{nil, false, fail.Control, "error should have been returned but was not", false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			message, correct := IsThreatCorrect(c.err, c.ok, c.sev)

			if correct && !c.correct {
				t.Error("error was accepted but should not have")
			} else if !correct && c.correct {
				t.Error("error should not have been accepted")
			}

			if !c.correct {
				idx := strings.LastIndexByte(message, ':')
				if idx == -1 {
					if c.err == nil {
						if message != c.message {
							t.Errorf("message does not fit, has \"%v\", expected \"%v\"",
								message, c.message)
						}
					} else {
						t.Error("message must contain ':'")
					}
				} else {
					stripped := message[:idx]
					if stripped != c.message {
						t.Errorf("message does not fit, has \"%v\", expected \"%v\"",
							stripped, c.message)
					}
				}
			} else {
				if message != "" {
					t.Errorf("message must be \"\" but was \"%v\"", message)
				}
			}
		})
	}
}
