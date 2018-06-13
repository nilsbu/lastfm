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

// AsyncReader is an interface for reading resources asynchronously.
type AsyncReader interface {
	Read(rsrc *Resource) <-chan ReadResult
}

// AsyncWriter is an interface for writing resources asynchronously.
type AsyncWriter interface {
	Write(data []byte, rsrc *Resource) <-chan error
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

// AsyncFileWriter is a file reader that delegates work to a readWorker.
type AsyncFileReader struct {
	Job chan ReadJob
}

func (r AsyncFileReader) Read(rsrc *Resource) <-chan ReadResult {
	out := make(chan ReadResult)
	go func(r AsyncFileReader, rsrc *Resource, out chan<- ReadResult) {
		back := make(chan ReadResult)
		r.Job <- ReadJob{Resource: rsrc, Back: back}

		out <- <-back
		close(back)
		close(out)
	}(r, rsrc, out)
	return out
}

// AsyncFileWriter is a file writer that delegates work to a writeWorker.
type AsyncFileWriter struct {
	Job chan WriteJob
}

func (r AsyncFileWriter) Write(data []byte, rsrc *Resource) <-chan error {
	out := make(chan error)
	go func(r AsyncFileWriter, data []byte, rsrc *Resource, out chan<- error) {
		back := make(chan error)
		r.Job <- WriteJob{Data: data, Resource: rsrc, Back: back}
		out <- <-back
		close(back)
		close(out)
	}(r, data, rsrc, out)
	return out
}

// AsyncDownloadGetter is a download getter that delegates work to read and
// write workers.
type AsyncDownloadGetter struct {
	downloader AsyncFileReader
	fileReader AsyncFileReader
	fileWriter AsyncFileWriter
}

func NewAsyncDownloadGetter(pool *Pool) *AsyncDownloadGetter {
	return &AsyncDownloadGetter{
		AsyncFileReader{pool.Download},
		AsyncFileReader{pool.ReadFile},
		AsyncFileWriter{pool.WriteFile},
	}
}

func (dg *AsyncDownloadGetter) Read(rsrc *Resource) <-chan ReadResult {
	out := make(chan ReadResult)
	go func(dg *AsyncDownloadGetter, rsrc *Resource, out chan<- ReadResult) {
		back := dg.fileReader.Read(rsrc)
		res := <-back
		if res.Err == nil {
			out <- res
			close(out)
			return
		}

		back = dg.downloader.Read(rsrc)
		res = <-back
		if res.Err == nil {
			back2 := dg.fileWriter.Write(res.Data, rsrc)
			<-back2
		}

		out <- res
		close(out)
	}(dg, rsrc, out)
	return out
}
