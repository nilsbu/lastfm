package testutils

import (
	"errors"

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
	return nil, errors.New("mock reader fails")
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
// When a resource is written more than all instances are stored.
// Certain resources can throw errors. FailSequences stores on which write
// requests an error is thrown. If an error is thrown, the data is not stored.
type Writer struct {
	Data          map[io.Resource][][]byte
	FailSequences map[io.Resource][]bool
}

// AsyncWriter is a mock writer analogous to Writer that works concurrently.
type AsyncWriter Writer

// NewWriter constructs a writer.
func NewWriter(failSequences map[io.Resource][]bool) *Writer {
	w := &Writer{
		Data:          make(map[io.Resource][][]byte),
		FailSequences: make(map[io.Resource][]bool),
	}

	// Copy to ensure that the input map is not altered
	for k, v := range failSequences {
		w.FailSequences[k] = v
	}
	return w
}

func (w *Writer) Write(data []byte, rsrc *io.Resource) (err error) {
	if seq, ok := w.FailSequences[*rsrc]; ok {
		if len(seq) > 0 {
			if !seq[0] {
				err = errors.New("mock writer fails")
			}

			w.FailSequences[*rsrc] = seq[1:]
		}
	}

	w.Data[*rsrc] = append(w.Data[*rsrc], data)
	return
}

// NewAsyncWriter constructs a concurrent writer.
func NewAsyncWriter(failSequences map[io.Resource][]bool) *AsyncWriter {
	return (*AsyncWriter)(NewWriter(failSequences))
}

func (w *AsyncWriter) Write(
	data []byte, rsrc *io.Resource) <-chan error {
	out := make(chan error)

	go func(data []byte, rsrc *io.Resource, out chan<- error) {
		out <- (*Writer)(w).Write(data, rsrc)
		close(out)
	}(data, rsrc, out)

	return out
}
