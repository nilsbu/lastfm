package cache

import (
	"errors"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Pool is a pool of readers and writers.
type Pool interface {
	Read(loc rsrc.Locator) <-chan io.ReadResult
	Write(data []byte, loc rsrc.Locator) <-chan error
}

type workerPool struct {
	readChan  chan io.ReadJob
	writeChan chan io.WriteJob
}

// NewPool constructs a Pool. It requires a non-epty list of workers, which are
// presumed to do identical jobs when provided with the same input.
func NewPool(
	readers []rsrc.Reader,
	writers []rsrc.Writer,
) (Pool, error) {
	if len(readers) == 0 {
		return nil, io.WrapError(fail.Critical,
			errors.New("pool must have at least one reader"))
	}
	if len(writers) == 0 {
		return nil, io.WrapError(fail.Critical,
			errors.New("pool must have at least one writer"))
	}

	pool := workerPool{
		make(chan io.ReadJob),
		make(chan io.WriteJob),
	}

	for _, reader := range readers {
		go readWorker(pool.readChan, reader)
	}

	for _, writer := range writers {
		go writeWorker(pool.writeChan, writer)
	}

	return pool, nil
}

func readWorker(jobs <-chan io.ReadJob, r rsrc.Reader) {
	for j := range jobs {
		data, err := r.Read(j.Locator)
		j.Back <- io.ReadResult{Data: data, Err: err}
	}
}

func writeWorker(jobs <-chan io.WriteJob, r rsrc.Writer) {
	for j := range jobs {
		err := r.Write(j.Data, j.Locator)
		j.Back <- err
	}
}

func (p workerPool) Read(loc rsrc.Locator) <-chan io.ReadResult {
	resultChan := make(chan io.ReadResult)
	p.readChan <- io.ReadJob{Locator: loc, Back: resultChan}
	return resultChan
}

func (p workerPool) Write(data []byte, loc rsrc.Locator) <-chan error {
	resultChan := make(chan error)
	p.writeChan <- io.WriteJob{Data: data, Locator: loc, Back: resultChan}
	return resultChan
}
