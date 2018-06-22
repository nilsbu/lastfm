package io

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

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
