package io

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
	"github.com/nilsbu/lastfm/test/mock"
)

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

			io := NewWebIO("0")

			var url string
			if c.hasURL {
				url = server.URL
			} else {
				url = ""
			}

			data, err := io.Read(stubURL(url))
			if str, ok := mock.IsThreatCorrect(err, c.ok, c.sev); !ok {
				t.Error(str)
			}

			if err == nil && string(data) != "response" {
				t.Errorf("wrong data read, has '%v', expected 'response'", string(data))
			}
		})
	}
}

func TestWebIOWrite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "response")
	}))
	defer server.Close()

	io := NewWebIO(mock.APIKey)

	err := io.Write([]byte("x"), stubURL(server.URL))
	if str, ok := mock.IsThreatCorrect(err, false, fail.Control); !ok {
		t.Error(str)
	}
}

func TestWebIORemove(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "response")
	}))
	defer server.Close()

	io := NewWebIO(mock.APIKey)

	err := io.Remove(stubURL(server.URL))
	if str, ok := mock.IsThreatCorrect(err, false, fail.Control); !ok {
		t.Error(str)
	}
}
