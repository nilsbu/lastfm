package io

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/rsrc"
)

func TestPoolReaderRead(t *testing.T) {
	ft := fastest.T{T: t}
	userInfo, _ := rsrc.UserInfo("SOX")

	testCases := []struct {
		rs   rsrc.Resource
		data string
		err  fastest.Code
	}{
		{rsrc.APIKey(), "XX", fastest.OK},
		{userInfo, "", fastest.OK},
		{userInfo, "lol", fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			r := make(PoolReader)
			c := r.Read(tc.rs)
			go func() {
				for job := range r {
					path1, err1 := job.Resource.Path()
					path2, err2 := tc.rs.Path()
					ft.Nil(err1)
					ft.Nil(err2)

					if path1 == path2 && tc.err == fastest.OK {
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
	userInfo, _ := rsrc.UserInfo("SOX")

	testCases := []struct {
		rs   rsrc.Resource
		data string
		err  fastest.Code
	}{
		{rsrc.APIKey(), "XX", fastest.OK},
		{userInfo, "", fastest.OK},
		{userInfo, "lol", fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			w := make(PoolWriter)
			c := w.Write([]byte(tc.data), tc.rs)
			var data []byte
			var rs rsrc.Resource = nil
			go func() {
				for job := range w {
					data = job.Data
					rs = job.Resource
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
			ft.Only(err == nil)

			path1, err1 := rs.Path()
			path2, err2 := tc.rs.Path()
			ft.Nil(err1)
			ft.Nil(err2)
			ft.Equals(path1, path2)
			ft.Only(err == nil)
			ft.Equals(string(data), tc.data)
		})
	}
}

func TestSeqReaderRead(t *testing.T) {
	ft := fastest.T{T: t}
	userInfo, _ := rsrc.UserInfo("SOX")

	testCases := []struct {
		rs   rsrc.Resource
		data string
		err  fastest.Code
	}{
		{rsrc.APIKey(), "XX", fastest.OK},
		{userInfo, "", fastest.OK},
		{userInfo, "lol", fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			r := make(SeqReader)

			go func() {
				for job := range r {
					path1, err1 := job.Resource.Path()
					path2, err2 := tc.rs.Path()
					ft.Nil(err1)
					ft.Nil(err2)

					if path1 == path2 && tc.err == fastest.OK {
						job.Back <- ReadResult{[]byte(tc.data), nil}
					} else {
						job.Back <- ReadResult{nil, errors.New("read failed")}
					}
				}
			}()

			data, err := r.Read(tc.rs)

			ft.Implies(err != nil, tc.err == fastest.Fail)
			ft.Implies(err == nil, tc.err == fastest.OK, err)
			ft.Only(err == nil)
			ft.Equals(string(data), tc.data)
		})
	}
}

func TestSeqWriterWrite(t *testing.T) {
	ft := fastest.T{T: t}
	userInfo, _ := rsrc.UserInfo("SOX")

	testCases := []struct {
		rs   rsrc.Resource
		data string
		err  fastest.Code
	}{
		{rsrc.APIKey(), "XX", fastest.OK},
		{userInfo, "", fastest.OK},
		{userInfo, "lol", fastest.Fail},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			w := make(SeqWriter)

			var data []byte
			var rs rsrc.Resource
			go func() {
				for job := range w {
					data = job.Data
					rs = job.Resource
					if tc.err == fastest.OK {
						job.Back <- nil
					} else {
						job.Back <- errors.New("read failed")
					}
				}
			}()

			err := w.Write([]byte(tc.data), tc.rs)
			ft.Implies(err != nil, tc.err == fastest.Fail)
			ft.Implies(err == nil, tc.err == fastest.OK, err)
			ft.Only(err == nil)

			path1, err1 := rs.Path()
			path2, err2 := tc.rs.Path()
			ft.Nil(err1)
			ft.Nil(err2)
			ft.Equals(path1, path2)
			ft.Only(err == nil)
			ft.Equals(string(data), tc.data)
		})
	}
}

type MockReader []byte

func (r MockReader) Read(rs rsrc.Resource) ([]byte, error) {
	if r != nil {
		return []byte(r), nil
	}
	return nil, errors.New("read failed")
}

type MockWriter struct {
	data []byte
	ok   bool
}

func (w *MockWriter) Write(data []byte, rs rsrc.Resource) error {
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

	res := <-PoolReader(p.Download).Read(rsrc.APIKey())
	ft.Nil(res.Err, res.Err)
	ft.Equals(string(res.Data), string(d))

	res = <-PoolReader(p.ReadFile).Read(rsrc.APIKey())
	ft.Nil(res.Err, res.Err)
	ft.Equals(string(res.Data), string(r))

	err := <-PoolWriter(p.WriteFile).Write(wStr, rsrc.APIKey())
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

			res := <-dg.Read(rsrc.APIKey())
			ft.Implies(res.Err != nil, tc.err == fastest.Fail)
			ft.Implies(res.Err == nil, tc.err == fastest.OK, res.Err)
			ft.Only(res.Err == nil)

			ft.Equals(string(res.Data), string(tc.data))

			ft.Only(tc.w)
			ft.Equals(string(w.data), string(tc.data))
		})
	}
}
