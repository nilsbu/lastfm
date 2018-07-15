package store

import (
	"errors"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Store provides IO access to multi-layered data storages. Layers are ordered
// from distant to close. Stores are intended to be used in a way that the most
// distant layer is the most permantent, i.e. a file is most likely to be found
// there but also the most expensive to interact with. The closer layers reverse
// that relation.
//
// A store performs a "lazy" file retrieval when Read() is invoked. That means
// it reads the file from the closest available storage. Possible file changes
// in more distant layers can go unnoticed that way. When a resource was read in
// a layer other than the most proximate, it is written to all closer layers to
// allow for faster retrieval in subsequent requests. To ensure that the most
// recent version of a resource is loaded, use Update() or see Fresh().
//
// Update() searches for a resource starting with the most distant layer. Once
// it finds the resource it overwrites potentially outdated versions in all
// closer layers.
//
// Write() writes the resource to all layers.
//
// Remove() removes the resource from all layers.
//
// TODO What role does fail.Threat play? Thread-safety?
type Store interface {
	rsrc.IO
	Update(loc rsrc.Locator) (data []byte, err error)
}

type cache struct {
	layers []pool
}

// New creates a store. The layers are described by ios. They are ordered from
// distant to close. Each layer must have at least one IO.
//
// TODO is it thread-safe?
func New(
	ios [][]rsrc.IO,
) (Store, error) {
	if len(ios) == 0 {
		return nil, fail.WrapError(fail.Critical,
			errors.New("store must have at least one layer"))
	}

	pools := make([]pool, len(ios))
	for i := range ios {
		pool, err := newPool(ios[i])
		if err != nil {
			return nil, err
		}
		pools[i] = pool
	}

	return &cache{pools}, nil
}

func (s *cache) Read(loc rsrc.Locator) (data []byte, err error) {
	return s.read(loc, len(s.layers)-1, -1)
}

func (s *cache) Update(loc rsrc.Locator) (data []byte, err error) {
	return s.read(loc, 0, 1)
}

func (s *cache) Write(data []byte, loc rsrc.Locator) error {
	return s.write(data, loc, len(s.layers)-1, -1)
}

func (s *cache) Remove(loc rsrc.Locator) error {
	return s.remove(loc)
}

func (s *cache) read(loc rsrc.Locator, start int, di int,
) (data []byte, err error) {

	idx, err := s.cascade(start, di, func(i int) (bool, error) {
		result := <-s.layers[i].read(loc)
		var tmpErr error
		data, tmpErr = result.data, result.err
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

func (s *cache) write(data []byte, loc rsrc.Locator, start int, di int) error {
	_, err := s.cascade(start, di, func(i int) (bool, error) {
		if err := <-s.layers[i].write(data, loc); err != nil {
			if f, ok := err.(fail.Threat); ok && f.Severity() == fail.Control {
				return true, nil
			}
			return false, err
		}
		return true, nil
	})
	return err
}

func (s *cache) remove(loc rsrc.Locator) error {
	_, err := s.cascade(len(s.layers)-1, -1, func(i int) (bool, error) {
		if err := <-s.layers[i].remove(loc); err != nil {
			if f, ok := err.(fail.Threat); ok && f.Severity() == fail.Control {
				return true, nil
			}
			return false, err
		}
		return true, nil
	})

	return err
}

func (s *cache) cascade(start int, di int, f func(i int) (bool, error)) (int, error) {
	for i := start; i >= 0 && i < len(s.layers); i += di {
		if cont, err := f(i); !cont {
			return i, err
		}
	}

	return -1, nil
}
