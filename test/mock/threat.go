package mock

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/fail"
)

// IsThreatCorrect checks a returned error that is supposed to implement Threat
// is returned when an error was expected and has the correct severity. In case
// of a mistake an error with an explanation is returned.
func IsThreatCorrect(
	err error,
	ok bool,
	sev fail.Severity,
) (message string, correct bool) {
	if err == nil {
		if ok {
			return "", true
		}
		return "error should have been returned but was not", false
	}

	if ok {
		return fmt.Sprintf("unexpected error: %v", err), false
	}

	f, tok := err.(fail.Threat)
	if !tok {
		return fmt.Sprintf("error must implement fail.Threat but does not: %v", err), false
	}

	hasSev := f.Severity()
	if hasSev != sev {
		str := fmt.Sprintf("severity must be '%v' but was '%v': %v",
			fail.GetSeverityString(sev), fail.GetSeverityString(f.Severity()), err)
		return str, false
	}

	return "", true
}
