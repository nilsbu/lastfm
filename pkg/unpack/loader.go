package unpack

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/pkg/errors"
)

// Loader is able to load and unpack data
type Loader interface {
	load(ob obtainer) (interface{}, error)
}

// cacheless is a loader that doesn't use a cache
type cacheless struct {
	rsrc.Reader
}

func NewCacheless(r rsrc.Reader) Loader {
	return &cacheless{Reader: r}
}

func (l *cacheless) load(ob obtainer) (interface{}, error) {
	return obtain(ob, l)
}

// cached if a buffer that stores tag information.
type cached struct {
	reader      rsrc.Reader
	requestChan chan cacheRequest
}

type cacheRequest struct {
	ob   obtainer
	back chan cachedResult
}

type cachedResult struct {
	loc  rsrc.Locator
	data interface{}
	err  error
}

// NewCached creates a buffer which can read and store information.
func NewCached(r rsrc.Reader) Loader {
	l := &cached{
		reader:      r,
		requestChan: make(chan cacheRequest),
	}

	go l.worker()
	return l
}

func (l *cached) load(ob obtainer) (interface{}, error) {
	back := make(chan cachedResult)

	l.requestChan <- cacheRequest{
		ob:   ob,
		back: back,
	}

	result := <-back
	if result.err != nil {
		return nil, result.err
	}

	return result.data, nil
}

func (l *cached) worker() {
	resultChan := make(chan cachedResult)

	requests := make(map[string][]cacheRequest)
	itemCounts := make(map[string]cachedResult)

	hasFatal := false

	for {
		select {
		case request := <-l.requestChan:
			path, _ := request.ob.locator().Path()
			if ic, ok := itemCounts[path]; ok {
				request.back <- ic
				close(request.back)
			} else if hasFatal {
				request.back <- cachedResult{
					loc:  request.ob.locator(),
					data: nil,
					err:  errors.New("abort due to previous error"),
				}
				close(request.back)
			} else {
				path, _ := request.ob.locator().Path()
				requests[path] = append(requests[path], request)
				go func(request cacheRequest) {
					data, err := obtain(request.ob, l.reader)
					if err != nil {
						hasFatal = hasFatal || isFatal(err)
						resultChan <- cachedResult{request.ob.locator(), nil, err}
					} else {
						resultChan <- cachedResult{request.ob.locator(), data, nil}
					}
				}(request)
			}

		case item := <-resultChan:
			path, _ := item.loc.Path()
			itemCounts[path] = item

			for _, request := range requests[path] {
				request.back <- item
				close(request.back)
			}

			requests[path] = nil
		}
	}
}

func isFatal(err error) bool {
	switch err := err.(type) {
	case *LastfmError:
		return err.IsFatal()
	default:
		return true
	}
}
