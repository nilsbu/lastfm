package io

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestPoolReaderRead(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		rsrc *Resource
		data string
		err  fastest.Code
	}{
		{NewAPIKey(), "XX", fastest.OK},
		{NewUserInfo("SOX"), "", fastest.OK},
		{NewUserInfo("A"), "lol", fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			r := make(PoolReader)
			c := r.Read(tc.rsrc)
			go func() {
				for job := range r {
					if *job.Resource == *tc.rsrc && tc.err == fastest.OK {
						job.Back <- ReadResult{[]byte(tc.data), nil}
					} else {
						job.Back <- ReadResult{nil, errors.New("read failed")}
					}
				}
			}()

			res := <-c
			ft.Implies(res.Err != nil, tc.err == fastest.Fail)
			ft.Implies(res.Err == nil, tc.err == fastest.OK, res.Err)
			ft.Only(res.Err == nil)
			ft.Equals(string(res.Data), tc.data)
		})
	}
}

func TestPoolWriterWrite(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		rsrc *Resource
		data string
		err  fastest.Code
	}{
		{NewAPIKey(), "XX", fastest.OK},
		{NewUserInfo("SOX"), "", fastest.OK},
		{NewUserInfo("A"), "lol", fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			w := make(PoolWriter)
			c := w.Write([]byte(tc.data), tc.rsrc)
			var data []byte
			var rsrc *Resource
			go func() {
				for job := range w {
					data = job.Data
					rsrc = job.Resource
					if tc.err == fastest.OK {
						job.Back <- nil
					} else {
						job.Back <- errors.New("read failed")
					}
				}
			}()

			err := <-c
			ft.Implies(err != nil, tc.err == fastest.Fail)
			ft.Implies(err == nil, tc.err == fastest.OK, err)
			ft.Equals(*rsrc, *tc.rsrc)
			ft.Only(err == nil)
			ft.Equals(string(data), tc.data)
		})
	}
}

type MockReader []byte

func (r MockReader) Read(rsrc *Resource) ([]byte, error) {
	if r != nil {
		return []byte(r), nil
	}
	return nil, errors.New("read failed")
}

type MockWriter struct {
	data []byte
	ok   bool
}

func (w *MockWriter) Write(data []byte, rsrc *Resource) error {
	if w.ok {
		w.data = data
		return nil
	}
	return errors.New("write failed")
}

func TestPool(t *testing.T) {
	ft := fastest.T{T: t}

	d := MockReader("XYZ")
	r := MockReader("089i")
	w := &MockWriter{ok: true}

	wStr := []byte("uiokl.")

	p := NewPool(
		[]Reader{d},
		[]Reader{r},
		[]Writer{w})

	res := <-PoolReader(p.Download).Read(NewAPIKey())
	ft.Nil(res.Err, res.Err)
	ft.Equals(string(res.Data), string(d))

	res = <-PoolReader(p.ReadFile).Read(NewAPIKey())
	ft.Nil(res.Err, res.Err)
	ft.Equals(string(res.Data), string(r))

	err := <-PoolWriter(p.WriteFile).Write(wStr, NewAPIKey())
	ft.Nil(err, err)
	ft.Equals(string(w.data), string(wStr))
}

func TestAsyncDownloadGetterRead(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		data    []byte
		r, d, w bool
		err     fastest.Code
	}{
		// Read from disk (availability of download doesn't matter)
		{[]byte("A"), true, true, false, fastest.OK},
		{[]byte("B"), true, false, false, fastest.OK},
		// Downloaded and written
		{[]byte("C"), false, true, true, fastest.OK},
		// Read and download fails
		{[]byte("D"), false, false, false, fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			var r, d MockReader
			if tc.r {
				r = MockReader(tc.data)
			} else {
				r = MockReader(nil)
			}
			if tc.d {
				d = MockReader(tc.data)
			} else {
				d = MockReader(nil)
			}

			w := &MockWriter{ok: tc.w}

			dg := AsyncDownloadGetter(NewPool(
				[]Reader{d},
				[]Reader{r},
				[]Writer{w}))

			res := <-dg.Read(NewAPIKey())
			ft.Implies(res.Err != nil, tc.err == fastest.Fail)
			ft.Implies(res.Err == nil, tc.err == fastest.OK, res.Err)
			ft.Only(res.Err == nil)

			ft.Equals(string(res.Data), string(tc.data))

			ft.Only(tc.w)
			ft.Equals(string(w.data), string(tc.data))
		})
	}
}
