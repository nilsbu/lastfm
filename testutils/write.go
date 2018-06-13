package testutils

import "github.com/nilsbu/lastfm/io"

// Writer is a mock writer that stores data that is written.
// When a resource is written more than all instances are stored.
// Certain resources can throw errors. failSequences stores on which write
// requests an error is thrown. If an error is thrown, the data is not stored.
type Writer struct {
	data          map[io.Resource][][]byte
	failSequences map[io.Resource][]bool
}

// NewWriter constructs a writer.
func NewWriter(failSequences map[io.Resource][]bool) *Writer {
	return &Writer{
		data:          make(map[io.Resource][][]byte),
		failSequences: failSequences,
	}
}

func (w *Writer) Write(data []byte, rsrc *io.Resource) (err error) {
	if seq, ok := w.failSequences[*rsrc]; ok {
		if len(seq) > 0 {
			if !seq[0] {
				err = strerr("fail")
			}
			w.failSequences[*rsrc] = seq[1:]
		}
	}

	w.data[*rsrc] = append(w.data[*rsrc], data)
	return
}
