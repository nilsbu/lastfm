package io

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nilsbu/lastfm/rsrc"
)

// Reader is an interface for reading resources.
type Reader interface {
	Read(rs rsrc.Resource) (data []byte, err error)
}

// Writer is an interface for writing resources.
type Writer interface {
	Write(data []byte, rs rsrc.Resource) error
}

// Remover is an interface for removing a resources.
type Remover interface {
	Remove(rs rsrc.Resource) error
}

// FileReader is a reader for local files. It implements io.Reader.
type FileReader struct{}

func (FileReader) Read(rs rsrc.Resource) (data []byte, err error) {
	path, err := rs.Path()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(path)
}

// FileWriter is a writer for files. It implements io.Writer.
type FileWriter struct{}

func (FileWriter) Write(data []byte, rs rsrc.Resource) error {
	path, err := rs.Path()
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

type FileRemover struct{}

func (FileRemover) Remove(rs rsrc.Resource) error {
	path, err := rs.Path()
	if err != nil {
		return nil
	}
	return os.Remove(path)
}

// Downloader is a reader for Last.fm. It implements io.Reader.
type Downloader rsrc.Key

func (d Downloader) Read(rs rsrc.Resource) (data []byte, err error) {
	url, err := rs.URL(rsrc.Key(d))
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
