package store

import (
	"errors"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Store interface {
	rsrc.IO
	rsrc.Updater
}

// pool is a pool of IO workers. It contains workers for download, file reading
// and writing.
type pool struct {
	Pools []Pool
}

// TODO ...
func New(
	readers [][]rsrc.Reader,
	writers [][]rsrc.Writer,
	removers [][]rsrc.Remover,
) (Store, error) {
	if len(readers) != len(writers) {
		return nil, fail.WrapError(fail.Critical,
			errors.New("readers and writers must have equal numbers of layers"))
	}
	if len(readers) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("readers and writers must have at least one layer"))
	}

	pools := make([]Pool, len(readers))
	for i := range readers {
		pool, err := NewPool(readers[i], writers[i], removers[i])
		if err != nil {
			return nil, err
		}
		pools[i] = pool
	}

	return pool{pools}, nil
}

func (p pool) Read(loc rsrc.Locator) (data []byte, err error) {
	data, err = p.read(loc, len(p.Pools)-1, -1)
	return data, err
}

func (p pool) Update(loc rsrc.Locator) (data []byte, err error) {
	data, err = p.read(loc, 0, 1)
	return data, err
}

func (p pool) Write(data []byte, loc rsrc.Locator) error {
	_, err := p.write(data, loc, len(p.Pools)-1, -1)
	return err
}

func (p pool) Remove(loc rsrc.Locator) error {
	err := p.remove(loc)
	return err
}

func (p pool) read(loc rsrc.Locator, start int, di int,
) (data []byte, err error) {

	idx, err := p.cascade(start, di, func(i int) (bool, error) {
		result := <-p.Pools[i].Read(loc)
		var tmpErr error
		data, tmpErr = result.Data, result.Err
		if tmpErr == nil {
			return false, nil
		}
		if f, ok := tmpErr.(fail.Threat); ok && f.Severity() == fail.Control {
			return true, nil
		}
		return false, tmpErr
	})

	if idx < 0 {
		return nil, fail.WrapError(fail.Control, errors.New("resource not found"))
	}

	if err != nil {
		return nil, err
	}

	_, err = p.write(data, loc, idx+1, 1)
	return data, err
}

func (p pool) write(data []byte, loc rsrc.Locator, start int, di int) (int, error) {
	return p.cascade(start, di, func(i int) (bool, error) {
		if err := <-p.Pools[i].Write(data, loc); err != nil {
			if f, ok := err.(fail.Threat); ok && f.Severity() == fail.Control {
				return true, nil
			}
			return false, err
		}
		return true, nil
	})
}

func (p pool) remove(loc rsrc.Locator) error {
	_, err := p.cascade(len(p.Pools)-1, -1, func(i int) (bool, error) {
		if err := <-p.Pools[i].Remove(loc); err != nil {
			if f, ok := err.(fail.Threat); ok && f.Severity() == fail.Control {
				return true, nil
			}
			return false, err
		}
		return true, nil
	})

	return err
}

func (p pool) cascade(start int, di int, f func(i int) (bool, error)) (int, error) {
	for i := start; i >= 0 && i < len(p.Pools); i += di {
		if cont, err := f(i); !cont {
			return i, err
		}
	}

	return -1, nil
}
