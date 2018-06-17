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

// TODO name
type write struct {
	rsrc *io.Resource
	// in and stored should be []byte but used string for convenience
	in     string
	stored string
	err    fastest.Code
}

// TestWriter tests both NewWriter() and Write()
func TestWriter(t *testing.T) {
	ft := fastest.T{T: t}

	var writeTestCases = []struct {
		writes  []write
		success map[io.Resource]bool
	}{
		{
			[]write{},
			map[io.Resource]bool{},
		},
		{
			[]write{
				{io.NewUserInfo("AS"), "X", "X", fastest.OK},
			},
			map[io.Resource]bool{}, // implicit OK
		},
		{
			[]write{
				{io.NewAPIKey(), "00", "00", fastest.OK},
				{io.NewArtistInfo("AS"), "--", "", fastest.Fail},
				{io.NewAPIKey(), "-", "-", fastest.OK}, // overwritd
			},
			map[io.Resource]bool{
				*io.NewAPIKey():         true,
				*io.NewArtistInfo("AS"): false,
			},
		},
	}

	for i, tc := range writeTestCases {
		ft.Seq(fmt.Sprintf("#%d", i), func(ft fastest.T) {
			w := NewWriter(tc.success)

			for _, write := range tc.writes {
				err := w.Write([]byte(write.in), write.rsrc)
				ft.Implies(err == nil, write.err == fastest.OK, err)
				ft.Implies(err != nil, write.err == fastest.Fail)
				ft.Only(err == nil)
				ft.Equals(string(w.Data[*write.rsrc]), write.stored)
			}
		})
	}
}

func TestWriterInterface(t *testing.T) {
	ft := fastest.T{T: t}

	var w interface{} = NewWriter(map[io.Resource]bool{})
	_, ok := w.(io.Writer)
	ft.True(ok, "Writer does not implement io.Writer")
}
