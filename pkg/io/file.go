package io

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// FileReader is a Reader to read from the local file system.
type FileReader struct{}

// FileWriter is a Write to write to the local file system.
type FileWriter struct{}

// FileRemover is a Remover to remove from the local file system.
type FileRemover struct{}

// FileIO is an IO to access to the local file system.
type FileIO struct {
	FileReader
	FileWriter
	FileRemover
}

func (FileReader) Read(loc rsrc.Locator) ([]byte, error) {
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

func (FileWriter) Write(data []byte, loc rsrc.Locator) error {
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
	f.Close()
	return err
}

func (FileRemover) Remove(loc rsrc.Locator) error {
	path, err := loc.Path()
	if err != nil {
		return err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return fail.WrapError(fail.Control, errors.New("file does not exist"))
	}

	return os.Remove(path)
}

type updateRedirect struct {
	updater rsrc.Updater
}
