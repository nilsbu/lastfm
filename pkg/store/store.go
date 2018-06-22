package store

import (
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Store interface {
	io.Reader
	io.Writer
	io.Updater
}

// pool is a pool of IO workers. It contains workers for download, file reading
// and writing.
type pool struct {
	Download  chan io.ReadJob
	ReadFile  chan io.ReadJob
	WriteFile chan io.WriteJob
}

// New creates an IO worker pool with the given readers and writers.
func New(
	downloaders, fileReaders []io.Reader,
	fileWriters []io.Writer) pool {
	p := pool{
		make(chan io.ReadJob),
		make(chan io.ReadJob),
		make(chan io.WriteJob)}

	startWorkers(downloaders, fileReaders, fileWriters, p)

	return p
}

func startWorkers(
	downloaders, fileReaders []io.Reader,
	fileWriters []io.Writer,
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

func readWorker(jobs <-chan io.ReadJob, r io.Reader) {
	for j := range jobs {
		data, err := r.Read(j.Locator)
		j.Back <- io.ReadResult{Data: data, Err: err}
	}
}

func writeWorker(jobs <-chan io.WriteJob, r io.Writer) {
	for j := range jobs {
		err := r.Write(j.Data, j.Locator)
		j.Back <- err
	}
}

func (p pool) Read(loc rsrc.Locator) (data []byte, err error) {
	data, err = io.SeqReader(p.ReadFile).Read(loc)
	if err == nil {
		return data, nil
	}

	return p.Update(loc)
}

func (p pool) Update(loc rsrc.Locator) (data []byte, err error) {
	data, err = io.SeqReader(p.Download).Read(loc)
	if err == nil {
		// TODO what happens to the result
		p.Write(data, loc)
	}
	return data, err
}

func (p pool) Write(data []byte, loc rsrc.Locator) error {
	return io.SeqWriter(p.WriteFile).Write(data, loc)
}
