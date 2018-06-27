package io

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nilsbu/lastfm/pkg/fail"
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

type ReadWriter interface {
	Reader
	Writer
}

type IO interface {
	ReadWriter
	Remover
}

// FileIO privides access to the local file system. It implements Reader,
// Writer and Remover.
type FileIO struct{}

func (FileIO) Read(loc rsrc.Locator) ([]byte, error) {
	path, err := loc.Path()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			err = WrapError(fail.Control, err)
		default:
			// possible cause: bytes.ErrTooLarge
			err = WrapError(fail.Critical, err)
		}
	}

	return data, err
}

func (FileIO) Write(data []byte, loc rsrc.Locator) error {
	path, err := loc.Path()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0040755); err != nil {
			// Will be *PathError (?)
			return WrapError(fail.Critical, err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		// Will be *PathError
		return WrapError(fail.Critical, err)
	}

	_, err = f.Write(data)
	if err != nil {
		return WrapError(fail.Critical, err)
	}

	return nil
}

func (FileIO) Remove(loc rsrc.Locator) error {
	path, err := loc.Path()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return WrapError(fail.Critical, err)
	}
	return nil
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
		return nil, WrapError(fail.Critical, err)
	} else {
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 403:
			err = WrapError(fail.Critical,
				errors.New("forbidden (403), wrong API key?"))
		case 404:
			err = WrapError(fail.Suspicious,
				errors.New("resouce not found (404)"))
		default:
			err = WrapError(fail.Suspicious,
				fmt.Errorf("unexpected HTTP status: %v", resp.Status))
		}
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// possible cause: bytes.ErrTooLarge
		return nil, WrapError(fail.Critical, err)
	}
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

// FailIO is a reader and writer that always fails non-critically.
type FailIO struct{}

func (FailIO) Read(loc rsrc.Locator) (data []byte, err error) {
	return nil, WrapError(fail.Control,
		fmt.Errorf("cannot read on FailIO"))
}

func (FailIO) Write(data []byte, loc rsrc.Locator) (err error) {
	return WrapError(fail.Control,
		fmt.Errorf("cannot write on FailIO"))
}
