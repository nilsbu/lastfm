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
	readers []rsrc.Reader,
	writers []rsrc.Writer,
	removers []rsrc.Remover,
) (Pool, error) {
	r, err := NewReadPool(readers)
	if err != nil {
		return nil, err
	}

	w, err := NewWritePool(writers)
	if err != nil {
		return nil, err
	}

	rm, err := NewRemovePool(removers)
	if err != nil {
		return nil, err
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

// NewReadPool constructs a ReadPool. It fails if 0 readers were given.
func NewReadPool(readers []rsrc.Reader) (ReadPool, error) {
	if len(readers) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("pool must have at least one reader"))
	}

	pool := make(readWorker)

	for _, reader := range readers {
		go func(jobs <-chan readJob, r rsrc.Reader) {
			for j := range jobs {
				data, err := r.Read(j.Locator)
				j.Back <- ReadResult{Data: data, Err: err}
			}
		}(pool, reader)
	}

	return pool, nil
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

// NewWritePool constructs a RemovePool. It fails if 0 writers were given.
func NewWritePool(writers []rsrc.Writer) (WritePool, error) {
	if len(writers) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("pool must have at least one writer"))
	}

	pool := make(writeWorker)

	for _, writer := range writers {
		go func(jobs <-chan writeJob, r rsrc.Writer) {
			for j := range jobs {
				err := r.Write(j.Data, j.Locator)
				j.Back <- err
			}
		}(pool, writer)
	}

	return pool, nil
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

// NewRemovePool constructs a RemovePool. It fails if 0 removers were given.
func NewRemovePool(removers []rsrc.Remover) (RemovePool, error) {
	if len(removers) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("pool must have at least one remover"))
	}

	pool := make(removeWorker)

	for _, remover := range removers {
		go func(jobs <-chan removeJob, r rsrc.Remover) {
			for j := range jobs {
				err := r.Remove(j.Locator)
				j.Back <- err
			}
		}(pool, remover)
	}

	return pool, nil
}

func (rm removeWorker) Remove(loc rsrc.Locator) <-chan error {
	resultChan := make(chan error)
	rm <- removeJob{Locator: loc, Back: resultChan}
	return resultChan
}
