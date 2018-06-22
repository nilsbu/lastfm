package io

import "github.com/nilsbu/lastfm/pkg/rsrc"

// SeqReader provides sequential read access to an io.Pool.
type SeqReader chan ReadJob

func (r SeqReader) Read(rs rsrc.Resource) (data []byte, err error) {
	back := make(chan ReadResult)

	r <- ReadJob{Resource: rs, Back: back}
	res := <-back
	return res.Data, res.Err
}

// SeqWriter provides sequential write access to an io.Pool.
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

// Pool is a pool of IO workers. It contains workers for download, file reading
// and writing.
type Pool struct {
	Download  chan ReadJob
	ReadFile  chan ReadJob
	WriteFile chan WriteJob
}

// NewPool creates an IO worker pool with the given readers and writers.
func NewPool(downloaders, fileReaders []Reader, fileWriters []Writer) Pool {
	p := Pool{make(chan ReadJob), make(chan ReadJob), make(chan WriteJob)}

	startWorkers(downloaders, fileReaders, fileWriters, p)

	return p
}

func startWorkers(
	downloaders, fileReaders []Reader,
	fileWriters []Writer,
	p Pool) {
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

// TODO docu, name & test
type ForcedDownloadGetter Pool

func (dg ForcedDownloadGetter) Read(rs rsrc.Resource) (data []byte, err error) {
	data, err = SeqReader(dg.Download).Read(rs)
	if err == nil {
		// TODO what happens to the result
		SeqWriter(dg.WriteFile).Write(data, rs)
	}
	return data, err
}

// AsyncDownloadGetter is a download getter that delegates work to read and
// write workers.
type AsyncDownloadGetter ForcedDownloadGetter

func (dg AsyncDownloadGetter) Read(rs rsrc.Resource) (data []byte, err error) {
	data, err = SeqReader(dg.ReadFile).Read(rs)
	if err == nil {
		return data, nil
	}

	return ForcedDownloadGetter(dg).Read(rs)
}
