package io

import (
	"errors"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestSeqReaderRead(t *testing.T) {
	userInfo, _ := rsrc.UserInfo("SOX")

	cases := []struct {
		rs   rsrc.Resource
		data string
		ok   bool
	}{
		{rsrc.APIKey(), "XX", true},
		{userInfo, "", true},
		{userInfo, "lol", false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			r := make(SeqReader)

			go func() {
				for job := range r {
					path1, _ := job.Resource.Path()
					path2, _ := c.rs.Path()

					if path1 == path2 && c.ok == true {
						job.Back <- ReadResult{[]byte(c.data), nil}
					} else {
						job.Back <- ReadResult{nil, errors.New("read failed")}
					}
				}
			}()

			data, err := r.Read(c.rs)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but non occurred")
			}
			if err == nil {
				if string(data) != c.data {
					t.Errorf("read data is wrong:\nread:     '%v'\nexpected: '%v'",
						string(data), c.data)
				}
			}
		})
	}
}

func TestSeqWriterWrite(t *testing.T) {
	userInfo, _ := rsrc.UserInfo("SOX")

	cases := []struct {
		rs   rsrc.Resource
		data string
		ok   bool
	}{
		{rsrc.APIKey(), "XX", true},
		{userInfo, "", true},
		{userInfo, "lol", false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			w := make(SeqWriter)

			var data []byte
			var rs rsrc.Resource
			go func() {
				for job := range w {
					data = job.Data
					rs = job.Resource
					if c.ok {
						job.Back <- nil
					} else {
						job.Back <- errors.New("read failed")
					}
				}
			}()

			err := w.Write([]byte(c.data), c.rs)
			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Error("expected error but non occurred")
			}
			if err == nil {
				path1, _ := rs.Path()
				path2, _ := c.rs.Path()

				if path1 != path2 {
					t.Errorf("written to wrong path:\nwritten:  '%v'\nexpected: '%v'",
						path1, path2)
				}
				if string(data) != c.data {
					t.Errorf("written data is wrong:\nread:     '%v'\nexpected: '%v'",
						string(data), c.data)
				}
			}
		})
	}
}
