package fail

import (
	"errors"
	"testing"
)

func TestSeverity(t *testing.T) {
	if Control >= Suspicious {
		t.Error("severity 'control' must be lower than 'suspicious'")
	}
	if Suspicious >= Critical {
		t.Error("severity 'suspicious' must be lower than 'critical'")
	}
}

func TestAssessedErrorInterface(t *testing.T) {
	if _, ok := error(&AssessedError{}).(Threat); !ok {
		t.Error("ee")
	}
}

func TestAssessedErrorError(t *testing.T) {
	cases := []struct {
		sev    Severity
		msg    string
		errStr string
	}{
		{Control, "I'm harmless", "[control] I'm harmless"},
		{Suspicious, "xx", "[suspicious] xx"},
		{Critical, "AAHH!", "[critical] AAHH!"},
		{-1, "this is invalid", "this is invalid"},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			err := &AssessedError{c.sev, errors.New(c.msg)}

			errStr := err.Error()
			if errStr != c.errStr {
				t.Errorf("faulty error message, '%v' ought to be '%v'",
					errStr, c.errStr)
			}
		})
	}
}

func TestAssessedErrorSeverity(t *testing.T) {
	cases := []struct {
		sev Severity
	}{
		{Control},
		{Suspicious},
		{Critical},
		{-1},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			err := &AssessedError{c.sev, errors.New("")}

			sev := err.Severity()
			if sev != c.sev {
				// see Severity consts if this fails
				t.Errorf("wrong severity level, was '%v', expected '%v'", sev, c.sev)
			}
		})
	}
}

func TestGetSeverityString(t *testing.T) {
	cases := []struct {
		sev Severity
		str string
	}{
		{Control, "control"},
		{Suspicious, "suspicious"},
		{Critical, "critical"},
		{-1, ""},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			str := GetSeverityString(c.sev)

			if str != c.str {
				t.Errorf("wrong severity string, was '%v', expected '%v'", str, c.str)
			}
		})
	}
}

func TestErrorConstructors(t *testing.T) {
	err := WrapError(Suspicious, errors.New("abc"))

	if err.Sev != Suspicious {
		t.Error("severity must be 'Suspicious'")
	}

	str := err.Err.Error()
	if str != "abc" {
		t.Errorf("wrong error message, was '%v', expected 'abc'", str)
	}
}
