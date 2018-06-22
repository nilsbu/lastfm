package io

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Updater interface {
	Update(rs rsrc.Resource) (data []byte, err error)
}

type Store interface {
	Reader
	Writer
	Updater
}

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

// pool is a pool of IO workers. It contains workers for download, file reading
// and writing.
type pool struct {
	Download  chan ReadJob
	ReadFile  chan ReadJob
	WriteFile chan WriteJob
}

// NewStore creates an IO worker pool with the given readers and writers.
func NewStore(downloaders, fileReaders []Reader, fileWriters []Writer) pool {
	p := pool{make(chan ReadJob), make(chan ReadJob), make(chan WriteJob)}

	startWorkers(downloaders, fileReaders, fileWriters, p)

	return p
}

func startWorkers(
	downloaders, fileReaders []Reader,
	fileWriters []Writer,
	p pool) {
	for _, d := range downloaders {
		go readWorker(p.Download, d)
	}

	for _, r := range fileReaders {
		go readWorker(p.ReadFile, r)
	}

	for _, w := range fileWriters {
		go writeWorker(p.WriteFile, w)
	}
}

func readWorker(jobs <-chan ReadJob, r Reader) {
	for j := range jobs {
		data, err := r.Read(j.Resource)
		j.Back <- ReadResult{data, err}
	}
}

func writeWorker(jobs <-chan WriteJob, r Writer) {
	for j := range jobs {
		err := r.Write(j.Data, j.Resource)
		j.Back <- err
	}
}

func (p pool) Read(rs rsrc.Resource) (data []byte, err error) {
	data, err = SeqReader(p.ReadFile).Read(rs)
	if err == nil {
		return data, nil
	}

	return p.Update(rs)
}

func (p pool) Update(rs rsrc.Resource) (data []byte, err error) {
	data, err = SeqReader(p.Download).Read(rs)
	if err == nil {
		// TODO what happens to the result
		p.Write(data, rs)
	}
	return data, err
}

func (p pool) Write(data []byte, rs rsrc.Resource) error {
	return SeqWriter(p.WriteFile).Write(data, rs)
}

type updateRedirect struct {
	updater Updater
}

func RedirectUpdate(updater Updater) *updateRedirect {
	return &updateRedirect{updater: updater}
}

func (ur updateRedirect) Read(rs rsrc.Resource) (data []byte, err error) {
	return ur.updater.Update(rs)
}
