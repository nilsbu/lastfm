package io

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestCacheServer(t *testing.T) {
	go RunCacheServer(12322)

	steps := []struct {
		params url.Values
		status int
		result string
	}{
		{
			url.Values{"rsrc": {"r1"}, "action": {"write"}, "data": {"abc"}},
			200, "",
		},
		{
			url.Values{"rsrc": {"r1"}, "action": {"read"}},
			200, "abc",
		},
		{
			url.Values{"rsrc": {"r1"}, "action": {"remove"}},
			200, "",
		},
		{
			url.Values{"rsrc": {"r1"}, "action": {"read"}},
			404, "resource 'r1' not found",
		},
		{
			url.Values{"rsrc": {"r1"}},
			400, "no action provided",
		},
		{
			url.Values{"action": {"read"}},
			400, "no resource locator provided",
		},
		{
			url.Values{"rsrc": {"r2"}, "action": {"write"}},
			400, "no data provided",
		},
	}

	for _, s := range steps {
		resp, err := http.PostForm("http://localhost:12322", s.params)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if resp.StatusCode == 200 {
			data, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatal("unexpected error:", err)
			} else if string(data) != s.result {
				t.Errorf("unexpected result:\nhas:  %v\nwant: %v", string(data), s.result)
			}
		} else {
			if errs, ok := resp.Header["Err"]; !ok || len(errs) == 0 {
				t.Error("failing status must contain 'Err' in header")
			} else {
				if errs[0] != s.result {
					t.Errorf("unexpected error:\nhas:  %v\nwant: %v", errs[0], s.result)
				}
			}
		}
		if resp.StatusCode != s.status {
			t.Errorf("unexpected status: has %v, expected %v",
				resp.StatusCode, s.status)
		}
	}
}

type failResource struct{}

func (failResource) URL(string) (string, error) {
	return "", fail.WrapError(fail.Control, errors.New("not implemented"))
}
func (failResource) Path() (string, error) {
	return "", fail.WrapError(fail.Control, errors.New("not implemented"))
}

type job struct {
	loc  rsrc.Locator
	data []byte
	sev  fail.Severity
	ok   bool
}

func TestCache(t *testing.T) {
	cases := []struct {
		runServer bool
		write     []job
		remove    []job
		read      []job
	}{
		{
			true,
			[]job{{rsrc.APIKey(), []byte("xx"), fail.Control, true}},
			[]job{{rsrc.APIKey(), nil, fail.Control, true}},
			[]job{{rsrc.APIKey(), nil, fail.Control, false}},
		},
		{
			true,
			[]job{{rsrc.APIKey(), []byte("xx"), fail.Control, true}},
			[]job{},
			[]job{{rsrc.APIKey(), []byte("xx"), fail.Control, true}},
		},
		{
			true,
			[]job{},
			[]job{},
			[]job{{rsrc.APIKey(), nil, fail.Control, false}},
		},
		{
			true,
			[]job{{failResource{}, []byte("zz"), fail.Control, false}},
			[]job{{failResource{}, nil, fail.Control, false}},
			[]job{{failResource{}, nil, fail.Control, false}},
		},
		{
			false,
			[]job{{rsrc.APIKey(), nil, fail.Critical, false}},
			[]job{{rsrc.APIKey(), nil, fail.Critical, false}},
			[]job{{rsrc.APIKey(), nil, fail.Critical, false}},
		},
	}

	for i, c := range cases {
		t.Run("", func(t *testing.T) {
			if c.runServer {
				go RunCacheServer(12323 + i)
			}

			cacheIO := &CacheIO{Port: 12323 + i}

			for _, job := range c.write {
				err := cacheIO.Write(job.data, job.loc)
				if str, ok := mock.IsThreatCorrect(err, job.ok, job.sev); !ok {
					t.Error(str)
				}
			}

			for _, job := range c.remove {
				err := cacheIO.Remove(job.loc)
				if str, ok := mock.IsThreatCorrect(err, job.ok, job.sev); !ok {
					t.Error(str)
				}
			}

			for i, job := range c.read {
				data, err := cacheIO.Read(job.loc)
				if str, ok := mock.IsThreatCorrect(err, job.ok, job.sev); !ok {
					t.Error(str)
				}

				if err == nil {
					if string(data) != string(job.data) {
						t.Errorf("read #%v: wrong data:\nhas:  %v\nwant: %v",
							i, string(data), string(job.data))
					}
				}
			}
		})
	}
}
