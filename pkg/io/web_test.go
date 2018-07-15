package io

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nilsbu/lastfm/test/mock"
)

type stubURL string

func (stubURL) Path() (string, error) {
	return "", errors.New("no path")
}

func (url stubURL) URL(string) (string, error) {
	if url == "" {
		return "", errors.New("no URL")
	}
	return string(url), nil
}

func TestDownloaderRead(t *testing.T) {
	cases := []struct {
		hasURL       bool
		connectionOK bool
		httpCode     int
		ok           bool
	}{
		{false, true, http.StatusOK, false},
		{true, true, http.StatusOK, true},
		{true, true, http.StatusForbidden, false},
		{true, true, http.StatusNotFound, false},
		{true, true, http.StatusTeapot, false},
		{true, false, http.StatusOK, false},
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
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Errorf("expected error but none occurred")
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

	if err := io.Write([]byte("x"), stubURL(server.URL)); err == nil {
		t.Error("expected error but none occurred")
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

	if err := io.Remove(stubURL(server.URL)); err == nil {
		t.Error("expected error but none occurred")
	}
}
