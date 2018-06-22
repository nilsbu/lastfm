package mock

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type resolver func(rsrc.Resource) (string, error)

// APIKey is the API key that is used in mocked URLs.
const APIKey = "00000000000000000000000000000000"

// FileIO constructs a mock reader and a writer. Data what is written using the
// writer can be read by the reader. The mocked file system is initialized with
// content. They keys are the file paths and the values the data contained.
// Locations that are not amon the contents during initialization cannot be
// written to or read from. Locations initialized with value nil are considered
// to be non-existing files that can however be written to.
func FileIO(content map[string][]byte) (io.SeqReader, io.SeqWriter) {
	r := make(io.SeqReader)
	w := make(io.SeqWriter)
	go worker(content, r, w, func(rs rsrc.Resource) (string, error) {
		return rs.Path()
	})
	return r, w
}

// Downloader constructs a mock reader for data download. Content is
// initialized analagous to FileIO.
func Downloader(content map[string][]byte) io.SeqReader {
	r := make(io.SeqReader)
	go worker(
		content, r, make(chan io.WriteJob),
		func(rs rsrc.Resource) (string, error) {
			return rs.URL(APIKey)
		})
	return r
}

// AsyncFileIO constructs an asynchronous versions of the reader and writer from
// FileIO. The data access is thread-safe.
func AsyncFileIO(content map[string][]byte) (io.PoolReader, io.PoolWriter) {
	r := make(io.PoolReader)
	w := make(io.PoolWriter)
	go worker(content, r, w, func(rs rsrc.Resource) (string, error) {
		return rs.Path()
	})
	return r, w
}

// AsyncDownloader constructs an asynchronous versions of the reader from
// Downloader.
func AsyncDownloader(content map[string][]byte) io.PoolReader {
	r := make(io.PoolReader)
	go worker(
		content, r, make(chan io.WriteJob),
		func(rs rsrc.Resource) (string, error) {
			return rs.URL(APIKey)
		})
	return r
}

func worker(
	content map[string][]byte,
	readJobs <-chan io.ReadJob,
	writeJobs <-chan io.WriteJob,
	resolve resolver,
) {
	for {
		select {
		case job, ok := <-readJobs:
			if !ok {
				break
			}

			path, err := resolve(job.Resource)
			if err != nil {
				job.Back <- io.ReadResult{Data: nil, Err: err}
				continue
			}

			data, ok := content[path]
			if !ok || data == nil {
				job.Back <- io.ReadResult{
					Data: nil,
					Err:  fmt.Errorf("read at '%v' failed", path),
				}
			} else {
				job.Back <- io.ReadResult{Data: data, Err: nil}
			}
		case job, ok := <-writeJobs:
			if !ok {
				break
			}

			path, err := resolve(job.Resource)
			if err != nil {
				// cannot happen, include for safety
				job.Back <- err
				continue
			}

			if _, ok := content[path]; !ok {
				job.Back <- fmt.Errorf("write at '%v' failed", path)
			} else {
				content[path] = job.Data
				job.Back <- nil
			}
		}
	}
}
