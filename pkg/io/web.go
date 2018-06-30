package io

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type WebIO struct {
	Downloader
	FailWriter
	FailRemover
}

func NewWebIO(apiKey string) WebIO {
	return WebIO{
		Downloader: Downloader(apiKey),
	}
}

// Downloader is a reader for Last.fm. It implements io.Reader.
type Downloader string

// TODO test with net/http/httptest
func (d Downloader) Read(loc rsrc.Locator) (data []byte, err error) {
	url, err := loc.URL(string(d))
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fail.WrapError(fail.Critical, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusForbidden:
			err = fail.WrapError(fail.Critical,
				errors.New("forbidden (403), wrong API key?"))
		case http.StatusNotFound:
			err = fail.WrapError(fail.Suspicious,
				errors.New("resouce not found (404)"))
		default:
			err = fail.WrapError(fail.Suspicious,
				fmt.Errorf("unexpected HTTP status: %v", resp.Status))
		}
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	return data, err
}
