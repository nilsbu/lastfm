package io

// ReadResult is contains the return values of Reader.Read().
type ReadResult struct {
	Data []byte
	Err  error
}

// ReadJob is a job for reading a resource.
type ReadJob struct {
	Resource *Resource
	Back     chan<- ReadResult
}

// WriteJob is a job for writgin a resource.
type WriteJob struct {
	Data     []byte
	Resource *Resource
	Back     chan<- error
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
