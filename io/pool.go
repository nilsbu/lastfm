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

// ReadWorker is a worker that executes ReadJobs.
func ReadWorker(jobs <-chan ReadJob, r Reader) {
	for j := range jobs {
		data, err := r.Read(j.Resource)
		j.Back <- ReadResult{data, err}
	}
}

// WriteWorker is a worker that executes WriteJobs.
func WriteWorker(jobs <-chan WriteJob, r Writer) {
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
