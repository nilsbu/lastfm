package io

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// FailReader is a Reader that always fails.
type FailReader struct{}

// FailWriter is a Writer that always fails.
type FailWriter struct{}

// FailRemover is a Remover that always fails.
type FailRemover struct{}

// FailIO is an IO that always fails non-critically.
type FailIO struct {
	FailReader
	FailWriter
	FailRemover
}

func (FailReader) Read(loc rsrc.Locator) (data []byte, err error) {
	return nil, fmt.Errorf("cannot read on FailIO")
}

func (FailWriter) Write(data []byte, loc rsrc.Locator) (err error) {
	return fmt.Errorf("cannot write on FailIO")
}

func (FailRemover) Remove(loc rsrc.Locator) (err error) {
	return fmt.Errorf("cannot remove on FailIO")
}
