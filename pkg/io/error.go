package io

import (
	"github.com/nilsbu/lastfm/pkg/fail"
)

type Error struct {
	fail.AssessedError
}

func WrapError(
	severity fail.Severity,
	err error) *Error {
	return &Error{fail.AssessedError{
		Sev: severity,
		Err: err,
	}}
}
