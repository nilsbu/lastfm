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

type Cache struct {
	Layers []Pool
}

// TODO ...
func NewCache(
	ios [][]rsrc.IO,
) (*Cache, error) {
	if len(ios) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("store must have at least one layer"))
	}

	pools := make([]Pool, len(ios))
	for i := range ios {
		pool, err := NewPool(ios[i])
		if err != nil {
			return nil, err
		}
		pools[i] = pool
	}

	return &Cache{pools}, nil
	// return nil, nil
}

func (s *Cache) Read(loc rsrc.Locator) (data []byte, err error) {
	return s.read(loc, len(s.Layers)-1, -1)
}

func (s *Cache) Update(loc rsrc.Locator) (data []byte, err error) {
	return s.read(loc, 0, 1)
}

func (s *Cache) Write(data []byte, loc rsrc.Locator) error {
	return s.write(data, loc, len(s.Layers)-1, -1)
}

func (s *Cache) Remove(loc rsrc.Locator) error {
	return s.remove(loc)
}

func (s *Cache) read(loc rsrc.Locator, start int, di int,
) (data []byte, err error) {

	idx, err := s.cascade(start, di, func(i int) (bool, error) {
		result := <-s.Layers[i].Read(loc)
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

	return data, s.write(data, loc, idx+1, 1)
}

func (s *Cache) write(data []byte, loc rsrc.Locator, start int, di int) error {
	_, err := s.cascade(start, di, func(i int) (bool, error) {
		if err := <-s.Layers[i].Write(data, loc); err != nil {
			if f, ok := err.(fail.Threat); ok && f.Severity() == fail.Control {
				return true, nil
			}
			return false, err
		}
		return true, nil
	})
	return err
}

func (s *Cache) remove(loc rsrc.Locator) error {
	_, err := s.cascade(len(s.Layers)-1, -1, func(i int) (bool, error) {
		if err := <-s.Layers[i].Remove(loc); err != nil {
			if f, ok := err.(fail.Threat); ok && f.Severity() == fail.Control {
				return true, nil
			}
			return false, err
		}
		return true, nil
	})

	return err
}

func (s *Cache) cascade(start int, di int, f func(i int) (bool, error)) (int, error) {
	for i := start; i >= 0 && i < len(s.Layers); i += di {
		if cont, err := f(i); !cont {
			return i, err
		}
	}

	return -1, nil
}
