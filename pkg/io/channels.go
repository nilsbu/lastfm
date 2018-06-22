package io

import "github.com/nilsbu/lastfm/pkg/rsrc"

// SeqReader provides sequential read access to an io.pool.
type SeqReader chan ReadJob

func (r SeqReader) Read(rs rsrc.Resource) (data []byte, err error) {
	back := make(chan ReadResult)

	r <- ReadJob{Resource: rs, Back: back}
	res := <-back
	return res.Data, res.Err
}

// SeqWriter provides sequential write access to an io.pool.
type SeqWriter chan WriteJob

func (r SeqWriter) Write(data []byte, rs rsrc.Resource) error {
	back := make(chan error)
	r <- WriteJob{Data: data, Resource: rs, Back: back}
	return <-back
}

// ReadJob is a job for reading a resource.
type ReadJob struct {
	Resource rsrc.Resource
	Back     chan<- ReadResult
}

// WriteJob is a job for writing a resource.
type WriteJob struct {
	Data     []byte
	Resource rsrc.Resource
	Back     chan<- error
}

// ReadResult is contains the return values of Reader.Read().
type ReadResult struct {
	Data []byte
	Err  error
}
