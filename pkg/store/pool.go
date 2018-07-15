package store

import (
	"errors"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type readPool interface {
	read(loc rsrc.Locator) <-chan readResult
}

type writePool interface {
	write(data []byte, loc rsrc.Locator) <-chan error
}

type removePool interface {
	remove(loc rsrc.Locator) <-chan error
}

// pool is a pool of readers, writers and removers.
type pool interface {
	readPool
	writePool
	removePool
}

type workerPool struct {
	readPool
	writePool
	removePool
}

// newPool constructs a pool. It requires a non-empty list of workers, which are
// presumed to do identical jobs when provided with the same input.
func newPool(
	ios []rsrc.IO,
) (pool, error) {
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
					data, err := io.Read(j.loc)
					j.back <- readResult{data: data, err: err}
				case j := <-w:
					err := io.Write(j.data, j.loc)
					j.back <- err
				case j := <-rm:
					err := io.Remove(j.loc)
					j.back <- err
				}
			}
		}(io)
	}

	return workerPool{r, w, rm}, nil
}

type readWorker chan readJob

type readJob struct {
	loc  rsrc.Locator
	back chan<- readResult
}

// readResult is contains the return values of Reader.Read().
type readResult struct {
	data []byte
	err  error
}

func (r readWorker) read(loc rsrc.Locator) <-chan readResult {
	resultChan := make(chan readResult)
	r <- readJob{loc: loc, back: resultChan}
	return resultChan
}

type writeWorker chan writeJob

type writeJob struct {
	data []byte
	loc  rsrc.Locator
	back chan<- error
}

func (w writeWorker) write(data []byte, loc rsrc.Locator) <-chan error {
	resultChan := make(chan error)
	w <- writeJob{data: data, loc: loc, back: resultChan}
	return resultChan
}

type removeWorker chan removeJob

type removeJob struct {
	loc  rsrc.Locator
	back chan<- error
}

func (rm removeWorker) remove(loc rsrc.Locator) <-chan error {
	resultChan := make(chan error)
	rm <- removeJob{loc: loc, back: resultChan}
	return resultChan
}
