package store

import "github.com/nilsbu/lastfm/pkg/rsrc"

type UpToDate struct {
	Cache Store
}

func NewUpToDate(cache Store) *UpToDate {
	return &UpToDate{Cache: cache}
}

func (s *UpToDate) Read(loc rsrc.Locator) ([]byte, error) {
	return s.Cache.Update(loc)
}

func (s *UpToDate) Write(data []byte, loc rsrc.Locator) error {
	return s.Cache.Write(data, loc)
}

func (s *UpToDate) Remove(loc rsrc.Locator) error {
	return s.Cache.Remove(loc)
}
