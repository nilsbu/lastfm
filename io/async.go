package io

// AsyncReader is an interface for reading resources asynchronously.
type AsyncReader interface {
	Read(rsrc *Resource) <-chan ReadResult
}

// AsyncWriter is an interface for writing resources asynchronously.
type AsyncWriter interface {
	Write(data []byte, rsrc *Resource) <-chan error
}

// PoolReader is a file reader that delegates work to a readWorker.
type PoolReader chan ReadJob

func (r PoolReader) Read(rsrc *Resource) <-chan ReadResult {
	out := make(chan ReadResult)
	go func(r PoolReader, rsrc *Resource, out chan<- ReadResult) {
		back := make(chan ReadResult)
		r <- ReadJob{Resource: rsrc, Back: back}

		out <- <-back
		close(back)
		close(out)
	}(r, rsrc, out)
	return out
}

// PoolWriter is a file writer that delegates work to a writeWorker.
type PoolWriter chan WriteJob

func (r PoolWriter) Write(data []byte, rsrc *Resource) <-chan error {
	out := make(chan error)
	go func(r PoolWriter, data []byte, rsrc *Resource, out chan<- error) {
		back := make(chan error)
		r <- WriteJob{Data: data, Resource: rsrc, Back: back}
		out <- <-back
		close(back)
		close(out)
	}(r, data, rsrc, out)
	return out
}

// ReadJob is a job for reading a resource.
type ReadJob struct {
	Resource *Resource
	Back     chan<- ReadResult
}

// WriteJob is a job for writing a resource.
type WriteJob struct {
	Data     []byte
	Resource *Resource
	Back     chan<- error
}

// ReadResult is contains the return values of Reader.Read().
type ReadResult struct {
	Data []byte
	Err  error
}

// Pool is a pool of IO workers. It contains workers for download, file reading
// and writing.
type Pool struct {
	Download  chan ReadJob
	ReadFile  chan ReadJob
	WriteFile chan WriteJob
}

// NewPool creates an IO worker pool with the given readers and writers.
func NewPool(downloaders, fileReaders []Reader, fileWriters []Writer) Pool {
	p := Pool{make(chan ReadJob), make(chan ReadJob), make(chan WriteJob)}

	startWorkers(downloaders, fileReaders, fileWriters, p)

	return p
}

func startWorkers(
	downloaders, fileReaders []Reader,
	fileWriters []Writer,
	p Pool) {
	for _, d := range downloaders {
		go readWorker(p.Download, d)
	}

	for _, r := range fileReaders {
		go readWorker(p.ReadFile, r)
	}

	for _, w := range fileWriters {
		go writeWorker(p.WriteFile, w)
	}
}

func readWorker(jobs <-chan ReadJob, r Reader) {
	for j := range jobs {
		data, err := r.Read(j.Resource)
		j.Back <- ReadResult{data, err}
	}
}

func writeWorker(jobs <-chan WriteJob, r Writer) {
	for j := range jobs {
		err := r.Write(j.Data, j.Resource)
		j.Back <- err
	}
}

// AsyncDownloadGetter is a download getter that delegates work to read and
// write workers.
type AsyncDownloadGetter Pool

func (dg AsyncDownloadGetter) Read(rsrc *Resource) <-chan ReadResult {
	out := make(chan ReadResult)
	go func(dg AsyncDownloadGetter, rsrc *Resource, out chan<- ReadResult) {
		res := <-PoolReader(dg.ReadFile).Read(rsrc)
		if res.Err == nil {
			out <- res
			close(out)
			return
		}

		res = <-PoolReader(dg.Download).Read(rsrc)
		if res.Err == nil {
			// TODO what happens to the result
			<-PoolWriter(dg.WriteFile).Write(res.Data, rsrc)
		}

		out <- res
		close(out)
	}(dg, rsrc, out)
	return out
}