package io

import (
	"io/ioutil"
)

// Reader is an interface for reading resources.
type Reader interface {
	Read(rsrc *Resource) (data []byte, err error)
}

// Writer is an interface for writing resources.
type Writer interface {
	Write(data []byte, rsrc *Resource) error
}

// FileWriter is a writer for files. It implements io.Writer.
type FileWriter struct {
}

func (FileWriter) Write(data []byte, rsrc *Resource) error {
	return write(data, fmtPath(rsrc))
}

// FileReader is a reader for local files. It implements io.Reader.
type FileReader struct {
}

func (FileReader) Read(rsrc *Resource) (data []byte, err error) {
	path := fmtPath(rsrc)
	return ioutil.ReadFile(path)
}

// Downloader is a reader for Last.fm. It implements io.Reader.
type Downloader struct {
	apiKey APIKey
}

// NewDownloader creates a Downloader.
func NewDownloader(key APIKey) Downloader {
	return Downloader{key}
}

func (d Downloader) Read(rsrc *Resource) (data []byte, err error) {
	url := fmtURL(rsrc, d.apiKey)
	return download(url)
}

// DownloadGetter is a reader and writer that downloads a resource,
// saves it and returns it.
type DownloadGetter struct {
	Downloader
	FileReader
	FileWriter
}

// NewDownloadGetter creates a DownloadGetter
func NewDownloadGetter(key APIKey) DownloadGetter {
	return DownloadGetter{Downloader: Downloader{key}}
}

func (dg DownloadGetter) Read(rsrc *Resource) (data []byte, err error) {
	data, err = dg.FileReader.Read(rsrc)
	if err == nil {
		return
	}

	data, err = dg.Downloader.Read(rsrc)
	if err == nil {
		err = dg.FileWriter.Write(data, rsrc)
	}

	return
}
