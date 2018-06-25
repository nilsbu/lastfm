package mock

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/io"
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
) (io.SeqReader, io.SeqWriter, error) {

	files := make(map[string][]byte)
	for k, v := range content {
		path, err := resolve(k)
		if err != nil {
			return nil, nil, err
		}
		files[path] = v
	}

	r := make(io.SeqReader)
	w := make(io.SeqWriter)
	go worker(files, r, w, resolve)
	return r, w, nil
}

func worker(
	content map[string][]byte,
	readJobs <-chan io.ReadJob,
	writeJobs <-chan io.WriteJob,
	resolve func(loc rsrc.Locator) (string, error),
) {
	for {
		select {
		case job, ok := <-readJobs:
			if !ok {
				break
			}

			path, err := resolve(job.Locator)
			if err != nil {
				job.Back <- io.ReadResult{Data: nil, Err: err}
				continue
			}

			data, ok := content[path]
			if !ok || data == nil {
				job.Back <- io.ReadResult{
					Data: nil,
					Err: &fail.AssessedError{
						Sev: fail.Control, Err: fmt.Errorf("read at '%v' failed", path)},
				}
			} else {
				job.Back <- io.ReadResult{Data: data, Err: nil}
			}
		case job, ok := <-writeJobs:
			if !ok {
				break
			}

			path, err := resolve(job.Locator)
			if err != nil {
				job.Back <- err
				continue
			}

			if _, ok := content[path]; !ok {
				job.Back <- &fail.AssessedError{
					Sev: fail.Control, Err: fmt.Errorf("write at '%v' failed", path)}
			} else {
				content[path] = job.Data
				job.Back <- nil
			}
		}
	}
}
