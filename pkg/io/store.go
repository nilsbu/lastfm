package io

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/format"
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
// Write() writes the resource to all layers. The error is always nil.
//
// Remove() removes the resource from all layers. The error is always nil.
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
func new(
	ios [][]rsrc.IO,
	obChans []chan<- format.Formatter,
) (Store, error) {
	if len(ios) == 0 {
		return nil, errors.New("store must have at least one layer")
	}
	if len(ios) != len(obChans) {
		return nil, fmt.Errorf("has %v IOs but %v observer channels",
			len(ios), len(obChans))
	}

	pools := make([]pool, len(ios))
	for i := range ios {
		ob := newObserver(obChans[i])
		pool, err := newPool(ios[i], ob)
		if err != nil {
			return nil, err
		}
		pools[i] = pool
	}

	return &cache{pools}, nil
}

func dumpChan() chan<- format.Formatter {
	obChan := make(chan format.Formatter)
	go func() {
		for range obChan {
		}
	}()

	return obChan
}

func dumpChans(n int) []chan<- format.Formatter {
	chans := make([]chan<- format.Formatter, n)
	for i := 0; i < n; i++ {
		chans[i] = dumpChan()
	}
	return chans
}

// NewStore creates a store. The layers are described by ios. They are ordered from
// distant to close. Each layer must have at least one IO.
//
// TODO is it thread-safe?
func NewStore(
	ios [][]rsrc.IO,
) (Store, error) {
	return new(ios, dumpChans(len(ios)))
}

// NewObservedStore creates a store. The layers are described by ios. They are
// ordered from distant to close. Each layer must have at least one IO.
//
// Progress updates for each layer are sent to the obChans.
//
// TODO is it thread-safe?
func NewObservedStore(
	ios [][]rsrc.IO,
	obChans []chan<- format.Formatter,
) (Store, error) {
	return new(ios, obChans)
}

func (s *cache) Read(loc rsrc.Locator) (data []byte, err error) {
	return s.read(loc, len(s.layers)-1, -1)
}

func (s *cache) Update(loc rsrc.Locator) (data []byte, err error) {
	return s.read(loc, 0, 1)
}

func (s *cache) Write(data []byte, loc rsrc.Locator) error {
	s.write(data, loc, len(s.layers)-1, -1)
	return nil
}

func (s *cache) Remove(loc rsrc.Locator) error {
	s.cascade(len(s.layers)-1, -1, func(i int) bool {
		<-s.layers[i].remove(loc)
		return false
	})

	return nil
}

func (s *cache) read(loc rsrc.Locator, start int, di int,
) (data []byte, err error) {

	idx, found := s.cascade(start, di, func(i int) bool {
		result := <-s.layers[i].read(loc)
		data = result.data
		return result.err == nil
	})

	if !found {
		s, _ := loc.Path()
		return nil, fmt.Errorf("resource '%v' not found", s)
	}

	s.write(data, loc, idx+1, 1)
	return data, nil
}

func (s *cache) write(data []byte, loc rsrc.Locator, start int, di int) {
	s.cascade(start, di, func(i int) bool {
		<-s.layers[i].write(data, loc)
		return false
	})
}

func (s *cache) cascade(
	start, di int,
	f func(i int) (found bool),
) (idx int, found bool) {

	for i := start; i >= 0 && i < len(s.layers); i += di {
		if f(i) {
			return i, true
		}
	}

	return -1, false
}
