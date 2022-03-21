package io

import (
	"errors"

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
	r  readWorker
	w  writeWorker
	rm removeWorker

	o observer
}

// newPool constructs a pool. It requires a non-empty list of workers, which are
// presumed to do identical jobs when provided with the same input.
func newPool(
	ios []rsrc.IO,
	o observer,
) (pool, error) {
	if len(ios) == 0 {
		return nil, errors.New("pool must have at least one IO")
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
					o.NotifyRead(j.loc)
					j.back <- readResult{data: data, err: err}
				case j := <-w:
					err := io.Write(j.data, j.loc)
					o.NotifyWrite(j.loc)
					j.back <- err
				case j := <-rm:
					err := io.Remove(j.loc)
					o.NotifyRemove(j.loc)
					j.back <- err
				}
			}
		}(io)
	}

	return workerPool{r, w, rm, o}, nil
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

func (wp workerPool) read(loc rsrc.Locator) <-chan readResult {
	wp.o.RequestRead(loc)
	resultChan := make(chan readResult)
	wp.r <- readJob{loc: loc, back: resultChan}
	return resultChan
}

type writeWorker chan writeJob

type writeJob struct {
	data []byte
	loc  rsrc.Locator
	back chan<- error
}

func (wp workerPool) write(data []byte, loc rsrc.Locator) <-chan error {
	wp.o.RequestWrite(loc)
	resultChan := make(chan error)
	wp.w <- writeJob{data: data, loc: loc, back: resultChan}
	return resultChan
}

type removeWorker chan removeJob

type removeJob struct {
	loc  rsrc.Locator
	back chan<- error
}

func (wp workerPool) remove(loc rsrc.Locator) <-chan error {
	wp.o.RequestRemove(loc)
	resultChan := make(chan error)
	wp.rm <- removeJob{loc: loc, back: resultChan}
	return resultChan
}
