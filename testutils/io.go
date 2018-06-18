package testutils

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/io"
)

// Reader is a mock reader that stores data that is returned on request.
// For each request of the same resource the same data is returned.
type Reader map[io.Resource][]byte

// AsyncReader is a mock reader analogous to Reader that works concurrently.
type AsyncReader Reader

func (r Reader) Read(rsrc *io.Resource) (data []byte, err error) {
	data, ok := r[*rsrc]
	if ok {
		return data, nil
	}
	return nil, fmt.Errorf("mock reader fails (%v)", rsrc)
}

func (r AsyncReader) Read(rsrc *io.Resource) <-chan io.ReadResult {
	out := make(chan io.ReadResult)
	go func(rsrc *io.Resource, out chan<- io.ReadResult) {
		data, err := Reader(r).Read(rsrc)
		out <- io.ReadResult{Data: data, Err: err}
		close(out)
	}(rsrc, out)
	return out
}

// Writer is a mock writer that stores data that is written.
// When a resource is written more than one, the last value is kept.
// Certain resources can throw errors. FailSequences stores on which write
// requests an error is thrown. If an error is thrown, the data is not stored.
type Writer struct {
	Data    map[io.Resource][]byte
	Success map[io.Resource]bool
}

// AsyncWriter is a mock writer analogous to Writer that works concurrently.
type AsyncWriter Writer

// NewWriter constructs a writer.
func NewWriter(success map[io.Resource]bool) *Writer {
	w := &Writer{
		Data:    make(map[io.Resource][]byte),
		Success: make(map[io.Resource]bool),
	}

	w.Success = success
	return w
}

func (w *Writer) Write(data []byte, rsrc *io.Resource) (err error) {
	if success, ok := w.Success[*rsrc]; ok && !success {
		return errors.New("mock writer fails")
	}

	w.Data[*rsrc] = data
	return
}

// TODO no AsyncWriter right now, since Writer is not thread-safe
