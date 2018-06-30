package store

import (
	"errors"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// ReadPool is a pool of readers.
type ReadPool interface {
	Read(loc rsrc.Locator) <-chan ReadResult
}

// WritePool is a pool of writers.
type WritePool interface {
	Write(data []byte, loc rsrc.Locator) <-chan error
}

// RemovePool is a pool removers.
type RemovePool interface {
	Remove(loc rsrc.Locator) <-chan error
}

// Pool is a pool of readers, writers and removers.
type Pool interface {
	ReadPool
	WritePool
	RemovePool
}

type workerPool struct {
	ReadPool
	WritePool
	RemovePool
}

// NewPool constructs a Pool. It requires a non-empty list of workers, which are
// presumed to do identical jobs when provided with the same input.
func NewPool(
	ios []rsrc.IO,
) (Pool, error) {
	if len(ios) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("pool must have at least one IO"))
	}

	r := make(readWorker)
	w := make(writeWorker)
	rm := make(removeWorker)

	for _, io := range ios {
		go func(io rsrc.IO) {
			for {
				select {
				case j := <-r:
					data, err := io.Read(j.Locator)
					j.Back <- ReadResult{Data: data, Err: err}
				case j := <-w:
					err := io.Write(j.Data, j.Locator)
					j.Back <- err
				case j := <-rm:
					err := io.Remove(j.Locator)
					j.Back <- err
				}
			}
		}(io)
	}

	return workerPool{r, w, rm}, nil
}

type readWorker chan readJob

type readJob struct {
	Locator rsrc.Locator
	Back    chan<- ReadResult
}

// ReadResult is contains the return values of Reader.Read().
type ReadResult struct {
	Data []byte
	Err  error
}

func (r readWorker) Read(loc rsrc.Locator) <-chan ReadResult {
	resultChan := make(chan ReadResult)
	r <- readJob{Locator: loc, Back: resultChan}
	return resultChan
}

type writeWorker chan writeJob

type writeJob struct {
	Data    []byte
	Locator rsrc.Locator
	Back    chan<- error
}

func (w writeWorker) Write(data []byte, loc rsrc.Locator) <-chan error {
	resultChan := make(chan error)
	w <- writeJob{Data: data, Locator: loc, Back: resultChan}
	return resultChan
}

type removeWorker chan removeJob

type removeJob struct {
	Locator rsrc.Locator
	Back    chan<- error
}

func (rm removeWorker) Remove(loc rsrc.Locator) <-chan error {
	resultChan := make(chan error)
	rm <- removeJob{Locator: loc, Back: resultChan}
	return resultChan
}
