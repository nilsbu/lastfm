package testutils

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/io"
)

var readTestCases = []struct {
	rsrc *io.Resource
	data []byte
	err  fastest.Code
}{
	{io.NewUserInfo("AS"), []byte("xx"), fastest.OK},
	// Twice to ensure that requests are repeatable
	{io.NewUserInfo("AS"), []byte("xx"), fastest.OK},
	{io.NewArtistInfo("AS"), nil, fastest.Fail},
	{io.NewUserRecentTracks("D", 1, 86400), []byte("a"), fastest.OK},
	{io.NewUserRecentTracks("D", 2, 0), []byte("c"), fastest.OK},
	{io.NewUserRecentTracks("D", 1, 0), []byte("b"), fastest.OK},
	{io.NewUserRecentTracks("D", 3, 0), nil, fastest.Fail},
}

var mockReaderContent = map[io.Resource][]byte{
	*io.NewUserInfo("AS"):                  []byte("xx"),
	*io.NewUserRecentTracks("D", 1, 86400): []byte("a"),
	*io.NewUserRecentTracks("D", 1, 0):     []byte("b"),
	*io.NewUserRecentTracks("D", 2, 0):     []byte("c"),
}

func TestReader(t *testing.T) {
	ft := fastest.T{T: t}

	r := Reader(mockReaderContent)

	for i, tc := range readTestCases {
		ft.Seq(fmt.Sprintf("#%d", i), func(ft fastest.T) {
			data, err := r.Read(tc.rsrc)
			ft.Equals(tc.err == fastest.Fail, err != nil)
			ft.Equals(string(tc.data), string(data))
		})
	}
}

func TestReaderInterface(t *testing.T) {
	ft := fastest.T{T: t}

	var r interface{} = Reader(map[io.Resource][]byte{})
	_, ok := r.(io.Reader)
	ft.True(ok, "Reader does not implement io.Reader")
}

func TestAyncReader(t *testing.T) {
	ft := fastest.T{T: t}

	r := AsyncReader(mockReaderContent)

	for i, tc := range readTestCases {
		ft.Seq(fmt.Sprintf("#%d", i), func(ft fastest.T) {
			res := <-r.Read(tc.rsrc)
			ft.Equals(tc.err == fastest.Fail, res.Err != nil)
			ft.Equals(string(tc.data), string(res.Data))
		})
	}
}

func TestAsyncReaderInterface(t *testing.T) {
	ft := fastest.T{T: t}

	var r interface{} = AsyncReader(map[io.Resource][]byte{})
	_, ok := r.(io.AsyncReader)
	ft.True(ok, "AsyncReader does not implement io.AsyncReader")
}

var writeTestCases = []struct {
	rsrc *io.Resource
	// in and stored should be []byte but used string for convenience
	in     string
	stored []string
	err    fastest.Code
}{
	{io.NewUserInfo("AS"), "xx", []string{"xx"}, fastest.OK},
	{io.NewUserInfo("X"), "++", []string{"++"}, fastest.Fail},
	{io.NewUserInfo("AS"), "xy", []string{"xx", "xy"}, fastest.OK},
	{io.NewUserInfo("AS"), "--", []string{"xx", "xy", "--"}, fastest.Fail},
	{io.NewUserInfo("AS"), "++", []string{"xx", "xy", "--", "++"}, fastest.OK},
	{io.NewUserRecentTracks("D", 1, 86400), "", []string{""}, fastest.OK},
}

var mockWriterContent = map[io.Resource][]bool{
	*io.NewUserInfo("AS"): []bool{true, true, false},
	*io.NewUserInfo("X"):  []bool{false},
}

// TestWriter tests both NewWriter() and Write()
func TestWriter(t *testing.T) {
	ft := fastest.T{T: t}

	w := NewWriter(mockWriterContent)

	for i, tc := range writeTestCases {
		ft.Seq(fmt.Sprintf("#%d", i), func(ft fastest.T) {
			err := w.Write([]byte(tc.in), tc.rsrc)
			data, ok := w.Data[*tc.rsrc]

			ft.Implies(tc.err == fastest.OK, ok)
			ft.Equals(err != nil, tc.err == fastest.Fail)

			ft.Equals(len(tc.stored), len(data))
			for j := range tc.stored {
				ft.Equals(tc.stored[j], string(data[j]),
					fmt.Sprintf("Failed at element %v", j))
			}
		})
	}
}

func TestWriterInterface(t *testing.T) {
	ft := fastest.T{T: t}

	var w interface{} = NewWriter(map[io.Resource][]bool{})
	_, ok := w.(io.Writer)
	ft.True(ok, "Writer does not implement io.Writer")
}

// TestAsyncWriter tests both NewAsyncWriter() and Write()
func TestAsyncWriter(t *testing.T) {
	ft := fastest.T{T: t}

	w := NewAsyncWriter(mockWriterContent)

	for i, tc := range writeTestCases {
		ft.Seq(fmt.Sprintf("#%d", i), func(ft fastest.T) {
			err := <-w.Write([]byte(tc.in), tc.rsrc)
			data, ok := w.Data[*tc.rsrc]

			ft.Implies(tc.err == fastest.OK, ok)
			ft.Equals(err != nil, tc.err == fastest.Fail)

			ft.Equals(len(tc.stored), len(data))
			for j := range tc.stored {
				ft.Equals(tc.stored[j], string(data[j]),
					fmt.Sprintf("Failed at element %v", j))
			}
		})
	}
}

func TestAsyncWriterInterface(t *testing.T) {
	ft := fastest.T{T: t}

	var w interface{} = NewAsyncWriter(map[io.Resource][]bool{})
	_, ok := w.(io.AsyncWriter)
	ft.True(ok, "AsyncWriter does not implement io.AsyncWriter")
}
