package mock

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// APIKey is the API key that is used in mocked URLs.
const APIKey = "00000000000000000000000000000000"

type Resolver func(loc rsrc.Locator) (string, error)

// Path returns loc.Path()
func Path(loc rsrc.Locator) (string, error) {
	return loc.Path()
}

// URL returns loc.URL() with the default mock API key.
func URL(loc rsrc.Locator) (string, error) {
	return loc.URL(APIKey)
}

// IO constructs a mock reader and a writer. Data what is written using the
// writer can be read by the reader. The mocked data storage is initialized with
// content. They keys are the locations and the values the data contained.
// Locations that are not among the contents during initialization cannot be
// written to or read from. Locations initialized with value nil are considered
// to be non-existing files that can however be written to. Reader and writer
// can safely be copied. They are thread-safe.
func IO(
	content map[rsrc.Locator][]byte,
	resolve Resolver,
) (rsrc.Reader, rsrc.Writer, error) {

	files := make(map[string][]byte)
	for k, v := range content {
		path, err := resolve(k)
		if err != nil {
			return nil, nil, err
		}
		files[path] = v
	}

	r := make(chanReader)
	w := make(chanWriter)
	go worker(files, r, w, resolve)
	return r, w, nil
}

func worker(
	content map[string][]byte,
	readJobs <-chan readJob,
	writeJobs <-chan writeJob,
	resolve func(loc rsrc.Locator) (string, error),
) {
	for {
		select {
		case job := <-readJobs:
			path, err := resolve(job.Locator)
			if err != nil {
				job.Back <- readResult{Data: nil, Err: err}
				continue
			}

			data, ok := content[path]
			if !ok || data == nil {
				job.Back <- readResult{
					Data: nil,
					Err: &fail.AssessedError{
						Sev: fail.Control, Err: fmt.Errorf("read at '%v' failed", path)},
				}
			} else {
				job.Back <- readResult{Data: data, Err: nil}
			}
		case job := <-writeJobs:
			path, err := resolve(job.Locator)
			if err != nil {
				job.Back <- err
				continue
			}

			if _, ok := content[path]; !ok {
				job.Back <- &fail.AssessedError{
					Sev: fail.Critical, Err: fmt.Errorf("write at '%v' failed", path)}
			} else {
				content[path] = job.Data
				job.Back <- nil
			}
		}
	}
}

// TODO merge chanReader and chanWriter to a true IO (+ remove func)

type chanReader chan readJob

func (r chanReader) Read(loc rsrc.Locator) (data []byte, err error) {
	back := make(chan readResult)

	r <- readJob{Locator: loc, Back: back}
	res := <-back
	return res.Data, res.Err
}

type chanWriter chan writeJob

func (r chanWriter) Write(data []byte, loc rsrc.Locator) error {
	back := make(chan error)
	r <- writeJob{Data: data, Locator: loc, Back: back}
	return <-back
}

type readJob struct {
	Locator rsrc.Locator
	Back    chan<- readResult
}

type writeJob struct {
	Data    []byte
	Locator rsrc.Locator
	Back    chan<- error
}

type readResult struct {
	Data []byte
	Err  error
}
