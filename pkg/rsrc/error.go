package rsrc

import "github.com/nilsbu/lastfm/pkg/fail"

type LocatorError struct {
	fail.AssessedError
}

func WrapError(
	severity fail.Severity,
	err error) *LocatorError {
	return &LocatorError{fail.AssessedError{
		Sev: severity,
		Err: err,
	}}
}
