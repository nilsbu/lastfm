package io

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// Reader is an interface for reading resources.
type Reader interface {
	Read(loc rsrc.Locator) (data []byte, err error)
}

// Writer is an interface for writing resources.
type Writer interface {
	Write(data []byte, loc rsrc.Locator) error
}

// Remover is an interface for removing a resources.
type Remover interface {
	Remove(loc rsrc.Locator) error
}

type Updater interface {
	Update(loc rsrc.Locator) (data []byte, err error)
}

// FileIO privides access to the local file system. It implements Reader,
// Writer and Remover.
type FileIO struct{}

func (FileIO) Read(loc rsrc.Locator) (data []byte, err error) {
	path, err := loc.Path()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(path)
}

func (FileIO) Write(data []byte, loc rsrc.Locator) error {
	path, err := loc.Path()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0040755); err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func (FileIO) Remove(loc rsrc.Locator) error {
	path, err := loc.Path()
	if err != nil {
		return nil
	}
	return os.Remove(path)
}

// Downloader is a reader for Last.fm. It implements io.Reader.
type Downloader rsrc.Key

func (d Downloader) Read(loc rsrc.Locator) (data []byte, err error) {
	url, err := loc.URL(rsrc.Key(d))
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, err
}

type updateRedirect struct {
	updater Updater
}

func RedirectUpdate(updater Updater) *updateRedirect {
	return &updateRedirect{updater: updater}
}

func (ur updateRedirect) Read(loc rsrc.Locator) (data []byte, err error) {
	return ur.updater.Update(loc)
}
