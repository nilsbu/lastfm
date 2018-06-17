package io

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Reader is an interface for reading resources.
type Reader interface {
	Read(rsrc *Resource) (data []byte, err error)
}

// Writer is an interface for writing resources.
type Writer interface {
	Write(data []byte, rsrc *Resource) error
}

// FileReader is a reader for local files. It implements io.Reader.
type FileReader struct{}

func (FileReader) Read(rsrc *Resource) (data []byte, err error) {
	path := fmtPath(rsrc)
	return ioutil.ReadFile(path)
}

// FileWriter is a writer for files. It implements io.Writer.
type FileWriter struct{}

func (FileWriter) Write(data []byte, rsrc *Resource) error {
	path := fmtPath(rsrc)

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

// Downloader is a reader for Last.fm. It implements io.Reader.
type Downloader APIKey

func (d Downloader) Read(rsrc *Resource) (data []byte, err error) {
	url := fmtURL(rsrc, APIKey(d))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, err
}
