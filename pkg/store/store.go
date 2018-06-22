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
	Readers []chan io.ReadJob
	Writers []chan io.WriteJob
}

// New creates an IO worker pool with the given readers and writers.
func New(
	downloaders, fileReaders []io.Reader,
	fileWriters []io.Writer) Store {
	readers := []chan io.ReadJob{make(chan io.ReadJob), make(chan io.ReadJob)}
	writers := []chan io.WriteJob{make(chan io.WriteJob), make(chan io.WriteJob)}

	p := pool{
		readers,
		writers,
	}

	p.startWorkers(
		[][]io.Reader{downloaders, fileReaders},
		[][]io.Writer{[]io.Writer{io.FailIO{}}, fileWriters})

	return p
}

func (p pool) startWorkers(
	readers [][]io.Reader,
	writers [][]io.Writer) {

	for i := range readers {
		for _, d := range readers[i] {
			go readWorker(p.Readers[i], d)
		}
	}

	for i := range writers {
		for _, w := range writers[i] {
			go writeWorker(p.Writers[i], w)
		}
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
	data, err = io.SeqReader(p.Readers[1]).Read(loc)
	if err == nil {
		return data, nil
	}

	return p.Update(loc)
}

func (p pool) Update(loc rsrc.Locator) (data []byte, err error) {
	data, err = io.SeqReader(p.Readers[0]).Read(loc)
	if err == nil {
		// TODO what happens to the result
		p.Write(data, loc)
	}
	return data, err
}

func (p pool) Write(data []byte, loc rsrc.Locator) error {
	return io.SeqWriter(p.Writers[1]).Write(data, loc)
}
