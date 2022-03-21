package io

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// TODO Hide all types in io.

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

	return ioutil.ReadFile(path)
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
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		// Will be *PathError
		return err
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
		return errors.New("file does not exist")
	}

	os.Remove(path)
	return nil
}
