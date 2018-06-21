package io

import "github.com/nilsbu/lastfm/rsrc"

// AsyncReader is an interface for reading resources asynchronously.
type AsyncReader interface {
	Read(rs rsrc.Resource) <-chan ReadResult
}

// AsyncWriter is an interface for writing resources asynchronously.
type AsyncWriter interface {
	Write(data []byte, rs rsrc.Resource) <-chan error
}

// PoolReader is a file reader that delegates work to a readWorker.
type PoolReader chan ReadJob

func (r PoolReader) Read(rs rsrc.Resource) <-chan ReadResult {
	out := make(chan ReadResult)
	go func(r PoolReader, rs rsrc.Resource, out chan<- ReadResult) {
		back := make(chan ReadResult)
		r <- ReadJob{Resource: rs, Back: back}

		out <- <-back
		close(back)
		close(out)
	}(r, rs, out)
	return out
}

// PoolWriter is a file writer that delegates work to a writeWorker.
type PoolWriter chan WriteJob

func (r PoolWriter) Write(data []byte, rs rsrc.Resource) <-chan error {
	out := make(chan error)
	go func(r PoolWriter, data []byte, rs rsrc.Resource, out chan<- error) {
		back := make(chan error)
		r <- WriteJob{Data: data, Resource: rs, Back: back}
		out <- <-back
		close(back)
		close(out)
	}(r, data, rs, out)
	return out
}

// SeqReader provides sequential access to a PoolReader.
type SeqReader PoolReader

func (r SeqReader) Read(rs rsrc.Resource) (data []byte, err error) {
	res := <-PoolReader(r).Read(rs)
	return res.Data, res.Err
}

// SeqWriter provides sequential access to a PoolWriter.
type SeqWriter PoolWriter

func (r SeqWriter) Write(data []byte, rs rsrc.Resource) error {
	return <-PoolWriter(r).Write(data, rs)
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

// AsyncDownloadGetter is a download getter that delegates work to read and
// write workers.
type AsyncDownloadGetter Pool

func (dg AsyncDownloadGetter) Read(rs rsrc.Resource) <-chan ReadResult {
	out := make(chan ReadResult)
	go func(dg AsyncDownloadGetter, rs rsrc.Resource, out chan<- ReadResult) {
		res := <-PoolReader(dg.ReadFile).Read(rs)
		if res.Err == nil {
			out <- res
			close(out)
			return
		}

		res = <-PoolReader(dg.Download).Read(rs)
		if res.Err == nil {
			// TODO what happens to the result
			<-PoolWriter(dg.WriteFile).Write(res.Data, rs)
		}

		out <- res
		close(out)
	}(dg, rs, out)
	return out
}

// TODO docu, name & test
type ForcedDownloadGetter Pool

func (dg ForcedDownloadGetter) Read(rs rsrc.Resource) <-chan ReadResult {
	out := make(chan ReadResult)
	go func(dg ForcedDownloadGetter, rs rsrc.Resource, out chan<- ReadResult) {
		res := <-PoolReader(dg.Download).Read(rs)
		if res.Err == nil {
			// TODO what happens to the result
			<-PoolWriter(dg.WriteFile).Write(res.Data, rs)
		}

		out <- res
		close(out)
	}(dg, rs, out)
	return out
}
