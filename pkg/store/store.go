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
		data, err := r.Read(j.Resource)
		j.Back <- io.ReadResult{data, err}
	}
}

func writeWorker(jobs <-chan io.WriteJob, r io.Writer) {
	for j := range jobs {
		err := r.Write(j.Data, j.Resource)
		j.Back <- err
	}
}

func (p pool) Read(rs rsrc.Resource) (data []byte, err error) {
	data, err = io.SeqReader(p.ReadFile).Read(rs)
	if err == nil {
		return data, nil
	}

	return p.Update(rs)
}

func (p pool) Update(rs rsrc.Resource) (data []byte, err error) {
	data, err = io.SeqReader(p.Download).Read(rs)
	if err == nil {
		// TODO what happens to the result
		p.Write(data, rs)
	}
	return data, err
}

func (p pool) Write(data []byte, rs rsrc.Resource) error {
	return io.SeqWriter(p.WriteFile).Write(data, rs)
}
