package cache

import "errors"

// Pool is a pool of workers.
type Pool interface {
	Work(job Job) <-chan Result
}

// Job is a generic job for a pool.
type Job interface{}

// Result is the result of a done job by a pool.
type Result struct {
	Data interface{}
	Err  error
}

// Worker has the function Do that does a job and returns the result.
type Worker interface {
	Do(job Job) (data interface{}, err error)
}

// NewPool constructs a Pool. It requires a non-epty list of workers, which are
// presumed to do identical jobs when provided with the same input.
func NewPool(workers []Worker) (Pool, error) {
	if len(workers) == 0 {
		return nil, errors.New("pool needs at least one worker")
	}

	pool := make(workerPool)

	for _, worker := range workers {
		go func(worker Worker) {
			for in := range pool {
				data, err := worker.Do(in.job)
				in.back <- Result{Data: data, Err: err}
			}
		}(worker)
	}

	return pool, nil
}

type workerPool chan jobWithBack

type jobWithBack struct {
	job  Job
	back chan Result
}

func (p workerPool) Work(job Job) <-chan Result {
	resultChan := make(chan Result)
	p <- jobWithBack{job: job, back: resultChan}
	return resultChan
}
