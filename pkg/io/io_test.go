package io

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/test/mock"
)

const path = "../../test/files/temp.txt"

type stubPath string

func (stubPath) URL(string) (string, error) {
	return "", errors.New("no URL")
}

func (path stubPath) Path() (string, error) {
	if path == "" {
		return "",
			&fail.AssessedError{Sev: fail.Control, Err: errors.New("no path")}
	}
	return string(path), nil
}

func TestFileIORead(t *testing.T) {
	cases := []struct {
		hasPath bool
		hasFile bool
		ok      bool
		sev     fail.Severity
	}{
		{true, true, true, fail.Control},
		{false, true, false, fail.Control},
		{true, false, false, fail.Control},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io := FileIO{}

			var loc rsrc.Locator
			var err error
			if c.hasPath {
				loc = stubPath(path)
				if c.hasFile {
					err = io.Write([]byte("some text"), loc)
				} else {
					err = io.Remove(loc)
				}
				if err != nil {
					if f, ok := err.(fail.Threat); !ok || f.Severity() > fail.Control {
						t.Fatal("unexpected error during setup:", err)
					}
				}
			} else {
				loc = stubPath("")
			}

			data, err := io.Read(loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}

			if err == nil && string(data) != "some text" {
				t.Errorf("wrong data read, has '%v', expected 'some text'", string(data))
			}

		})
	}
}

func TestFileIOWrite(t *testing.T) {
	cases := []struct {
		hasPath bool
		ok      bool
		sev     fail.Severity
	}{
		{true, true, fail.Control},
		{false, false, fail.Control},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io := FileIO{}

			var loc rsrc.Locator
			if c.hasPath {
				loc = stubPath(path)
			} else {
				loc = stubPath("")
			}

			err := io.Write([]byte("some text"), loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}
			if err != nil {
				return
			}

			data, err := io.Read(loc)
			if err != nil {
				t.Fatal("error during read")
			}
			if string(data) != "some text" {
				t.Errorf("wrong data read, has '%v', expected 'some text'", string(data))
			}

		})
	}
}

func TestFileIORemove(t *testing.T) {
	cases := []struct {
		hasPath bool
		hasFile bool
		ok      bool
		sev     fail.Severity
	}{
		{true, true, true, fail.Control},
		{false, true, false, fail.Control},
		{true, false, false, fail.Control},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			io := FileIO{}

			var loc rsrc.Locator
			if c.hasPath {
				loc = stubPath(path)
				if c.hasFile {
					file, err := os.Create(path)
					if err != nil {
						t.Fatal("error during setup")
					}
					file.Close()
				}
			} else {
				loc = stubPath("")
			}

			err := io.Remove(loc)
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}
			if err != nil {
				return
			}

			if _, err := os.Stat(path); os.IsExist(err) {
				t.Error("file exists")
			}
		})
	}
}

type stubURL string

func (stubURL) Path() (string, error) {
	return "", errors.New("no path")
}

func (url stubURL) URL(string) (string, error) {
	if url == "" {
		return "", &fail.AssessedError{Sev: fail.Control, Err: errors.New("no URL")}
	}
	return string(url), nil
}

func TestDownloaderRead(t *testing.T) {
	cases := []struct {
		hasURL       bool
		connectionOK bool
		httpCode     int
		ok           bool
		sev          fail.Severity
	}{
		{false, true, http.StatusOK, false, fail.Control},
		{true, true, http.StatusOK, true, fail.Control},
		{true, true, http.StatusForbidden, false, fail.Critical},
		{true, true, http.StatusNotFound, false, fail.Suspicious},
		{true, true, http.StatusTeapot, false, fail.Suspicious},
		{true, false, http.StatusOK, false, fail.Critical},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(
				w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.httpCode)
				fmt.Fprint(w, "response")
			}))
			if c.connectionOK {
				defer server.Close()
			} else {
				server.Close()
			}

			d := Downloader("0")

			var url string
			if c.hasURL {
				url = server.URL
			} else {
				url = ""
			}

			data, err := d.Read(stubURL(url))
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}

			if err == nil && string(data) != "response" {
				t.Errorf("wrong data read, has '%v', expected 'response'", string(data))
			}
		})
	}
}

func TestFailReader(t *testing.T) {
	r := FailIO{}
	data, err := r.Read(rsrc.APIKey())
	if err == nil {
		t.Error("expected error but none occurred")
	} else {
		if f, ok := err.(fail.Threat); ok {
			if f.Severity() != fail.Control {
				t.Error("severity must be 'Control':", err)
			}
		} else {
			t.Error("error must implement Threat but does not:", err)
		}
	}
	if data != nil {
		t.Errorf("data should be nil but was '%v'", string(data))
	}
}

func TestFailWriter(t *testing.T) {
	r := FailIO{}
	err := r.Write([]byte("xyz"), rsrc.APIKey())
	if err == nil {
		t.Error("expected error but none occurred")
	} else {
		if f, ok := err.(fail.Threat); ok {
			if f.Severity() != fail.Control {
				t.Error("severity must be 'Control':", err)
			}
		} else {
			t.Error("error must implement Threat but does not:", err)
		}
	}
}
