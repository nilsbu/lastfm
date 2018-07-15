package io

import (
	"errors"
	"os"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

const path = "../../test/files/temp.txt"

type stubPath string

func (stubPath) URL(string) (string, error) {
	return "", errors.New("no URL")
}

func (path stubPath) Path() (string, error) {
	if path == "" {
		return "", errors.New("no path")
	}
	return string(path), nil
}

func TestFileIORead(t *testing.T) {
	cases := []struct {
		hasPath bool
		hasFile bool
		ok      bool
	}{
		{true, true, true},
		{false, true, false},
		{true, false, false},
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
					t.Fatal("unexpected error during setup:", err)
				}
			} else {
				loc = stubPath("")
			}

			data, err := io.Read(loc)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Errorf("expected error but none occurred")
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
	}{
		{true, true},
		{false, false},
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
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Errorf("expected error but none occurred")
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
	}{
		{true, true, true},
		{false, true, false},
		{true, false, false},
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
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Errorf("expected error but none occurred")
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
