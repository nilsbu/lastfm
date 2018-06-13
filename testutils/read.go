package testutils

import "github.com/nilsbu/lastfm/io"

// strerr is a simple error that consists of a string.
type strerr string

func (s strerr) Error() string {
	return string(s)
}

// Reader is a mock reader that stores data that is returned on request.
// For each request of the same resource the same data is returned.
type Reader struct {
	data map[io.Resource][]byte
}

// NewReader constructs a reader.
func NewReader(data map[io.Resource][]byte) *Reader {
	return &Reader{data}
}

func (r *Reader) Read(rsrc *io.Resource) (data []byte, err error) {
	data, ok := r.data[*rsrc]
	if ok {
		return data, nil
	}
	return nil, strerr("resource not found")
}
