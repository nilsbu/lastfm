package unpack

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/pkg/errors"
)

// CachedLoader if a buffer that stores tag information.
type CachedLoader struct {
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

// NewCachedLoader creates a buffer which can read and store information.
func NewCachedLoader(r rsrc.Reader) *CachedLoader {
	buf := &CachedLoader{
		reader:      r,
		requestChan: make(chan cacheRequest),
	}

	go buf.worker()
	return buf
}

func (buf *CachedLoader) Load(ob obtainer) (interface{}, error) {
	back := make(chan cachedResult)

	buf.requestChan <- cacheRequest{
		ob:   ob,
		back: back,
	}

	result := <-back
	if result.err != nil {
		return nil, result.err
	}

	return result.data, nil
}

func (buf *CachedLoader) worker() {
	resultChan := make(chan cachedResult)

	requests := make(map[string][]cacheRequest)
	itemCounts := make(map[string]cachedResult)

	hasFatal := false

	for {
		select {
		case request := <-buf.requestChan:
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
					data, err := obtain(request.ob, buf.reader)
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
	switch err.(type) {
	case *LastfmError:
		return err.(*LastfmError).IsFatal()
	default:
		return true
	}
}
