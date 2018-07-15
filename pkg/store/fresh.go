package store

import "github.com/nilsbu/lastfm/pkg/rsrc"

type fresh struct {
	Cache Store
}

// Fresh returns a handle to an existing store that cirumvents the "lazy"
// behavior of its Read() method. Instead, when Read() is called, Update() is
// executed which retrieves a resource from the most distant layer and
// overwrites the data in all closer layers. All other behaviors remain the
// same.
func Fresh(cache Store) rsrc.IO {
	return &fresh{Cache: cache}
}

func (s *fresh) Read(loc rsrc.Locator) ([]byte, error) {
	return s.Cache.Update(loc)
}

func (s *fresh) Write(data []byte, loc rsrc.Locator) error {
	return s.Cache.Write(data, loc)
}

func (s *fresh) Remove(loc rsrc.Locator) error {
	return s.Cache.Remove(loc)
}
