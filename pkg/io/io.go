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
			err = fail.WrapError(fail.Control, err)
		}
	}

	return data, err
}

func (FileIO) Write(data []byte, loc rsrc.Locator) error {
	path, err := loc.Path()
	if err != nil {
		return err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err = os.MkdirAll(dir, 0040755); err != nil {
			// Will be *PathError (?)
			return fail.WrapError(fail.Critical, err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		// Will be *PathError
		return fail.WrapError(fail.Critical, err)
	}

	_, err = f.Write(data)
	return err
}

func (FileIO) Remove(loc rsrc.Locator) error {
	path, err := loc.Path()
	if err != nil {
		return err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return fail.WrapError(fail.Control, errors.New("file does not exist"))
	}

	return os.Remove(path)
}

// Downloader is a reader for Last.fm. It implements io.Reader.
type Downloader rsrc.Key

// TODO test with net/http/httptest
func (d Downloader) Read(loc rsrc.Locator) (data []byte, err error) {
	url, err := loc.URL(rsrc.Key(d))
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fail.WrapError(fail.Critical, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusForbidden:
			err = fail.WrapError(fail.Critical,
				errors.New("forbidden (403), wrong API key?"))
		case http.StatusNotFound:
			err = fail.WrapError(fail.Suspicious,
				errors.New("resouce not found (404)"))
		default:
			err = fail.WrapError(fail.Suspicious,
				fmt.Errorf("unexpected HTTP status: %v", resp.Status))
		}
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}

type updateRedirect struct {
	updater rsrc.Updater
}

func RedirectUpdate(updater rsrc.Updater) *updateRedirect {
	return &updateRedirect{updater: updater}
}

func (ur updateRedirect) Read(loc rsrc.Locator) (data []byte, err error) {
	return ur.updater.Update(loc)
}

// FailIO is a reader and writer that always fails non-critically.
type FailIO struct{}

func (FailIO) Read(loc rsrc.Locator) (data []byte, err error) {
	return nil, fail.WrapError(fail.Control,
		fmt.Errorf("cannot read on FailIO"))
}

func (FailIO) Write(data []byte, loc rsrc.Locator) (err error) {
	return fail.WrapError(fail.Control,
		fmt.Errorf("cannot write on FailIO"))
}
