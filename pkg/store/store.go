package store

import (
	"github.com/nilsbu/lastfm/pkg/cache"
	"github.com/nilsbu/lastfm/pkg/fail"
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
	Pools []cache.Pool
}

// TODO ...
func New(
	readers [][]io.Reader,
	writers [][]io.Writer) (Store, error) {
	// TODO check lenghts

	pools := make([]cache.Pool, len(readers))
	for i := range readers {
		pool, err := cache.NewPool(readers[i], writers[i])
		if err != nil {
			return nil, err
		}
		pools[i] = pool
	}

	return pool{pools}, nil
}

func (p pool) Read(loc rsrc.Locator) (data []byte, err error) {
	result := <-p.Pools[1].Read(loc)
	data, err = result.Data, result.Err
	if err == nil {
		return data, nil
	}

	return p.Update(loc)
}

func (p pool) Update(loc rsrc.Locator) (data []byte, err error) {
	result := <-p.Pools[0].Read(loc)
	data, err = result.Data, result.Err
	if err == nil {
		// TODO what happens to the result
		p.Write(data, loc)
	}
	return data, err
}

func (p pool) Write(data []byte, loc rsrc.Locator) error {
	var ferr fail.Threat
	for i := len(p.Pools) - 1; i >= 0; i-- {
		if wErr := <-p.Pools[i].Write(data, loc); wErr != nil {
			f, ok := wErr.(fail.Threat)
			if !ok {
				f = io.WrapError(fail.Critical, wErr)
			}
			ferr = f
			break
		}
	}

	if ferr != nil && ferr.Severity() == fail.Control {
		return nil
	}

	err, _ := ferr.(error)
	return err
}
