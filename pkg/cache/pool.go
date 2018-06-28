package cache

import (
	"errors"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Pool is a pool of readers and writers.
type Pool interface {
	Read(loc rsrc.Locator) <-chan ReadResult
	Write(data []byte, loc rsrc.Locator) <-chan error
}

// ReadResult is contains the return values of Reader.Read().
type ReadResult struct {
	Data []byte
	Err  error
}

type workerPool struct {
	readChan  chan readJob
	writeChan chan writeJob
}

type readJob struct {
	Locator rsrc.Locator
	Back    chan<- ReadResult
}

type writeJob struct {
	Data    []byte
	Locator rsrc.Locator
	Back    chan<- error
}

// NewPool constructs a Pool. It requires a non-epty list of workers, which are
// presumed to do identical jobs when provided with the same input.
func NewPool(
	readers []rsrc.Reader,
	writers []rsrc.Writer,
) (Pool, error) {
	if len(readers) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("pool must have at least one reader"))
	}
	if len(writers) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("pool must have at least one writer"))
	}

	pool := workerPool{
		make(chan readJob),
		make(chan writeJob),
	}

	for _, reader := range readers {
		go readWorker(pool.readChan, reader)
	}

	for _, writer := range writers {
		go writeWorker(pool.writeChan, writer)
	}

	return pool, nil
}

func readWorker(jobs <-chan readJob, r rsrc.Reader) {
	for j := range jobs {
		data, err := r.Read(j.Locator)
		j.Back <- ReadResult{Data: data, Err: err}
	}
}

func writeWorker(jobs <-chan writeJob, r rsrc.Writer) {
	for j := range jobs {
		err := r.Write(j.Data, j.Locator)
		j.Back <- err
	}
}

func (p workerPool) Read(loc rsrc.Locator) <-chan ReadResult {
	resultChan := make(chan ReadResult)
	p.readChan <- readJob{Locator: loc, Back: resultChan}
	return resultChan
}

func (p workerPool) Write(data []byte, loc rsrc.Locator) <-chan error {
	resultChan := make(chan error)
	p.writeChan <- writeJob{Data: data, Locator: loc, Back: resultChan}
	return resultChan
}
