package organize

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/fastest"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/nilsbu/lastfm/test/mock"
)

func TestLoadAllDayPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		user  unpack.User
		until rsrc.Day
		data  [][]string
		dps   []unpack.PlayCount
		err   fastest.Code
	}{
		{
			unpack.User{Name: "", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(86400),
			[][]string{[]string{}, []string{}},
			nil,
			fastest.Fail,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ToDay(86400)},
			rsrc.ToDay(2 * 86400),
			[][]string{
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`},
				[]string{`{"recenttracks":{"track":[{"artist":{"#text":"XXX"}}], "@attr":{"totalPages":"1"}}}`},
			},
			[]unpack.PlayCount{unpack.PlayCount{"ASDF": 1}, unpack.PlayCount{"XXX": 1}},
			fastest.OK,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			[][]string{
				[]string{
					`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
					`{"recenttracks":{"track":[{"artist":{"#text":"Y"}}], "@attr":{"page":"2","totalPages":"3"}}}`,
					`{"recenttracks":{"track":[{"artist":{"#text":"Z"}}, {"artist":{"#text":"X"}}], "@attr":{"page":"3","totalPages":"3"}}}`,
				},
			},
			[]unpack.PlayCount{unpack.PlayCount{"X": 2, "Y": 1, "Z": 1}},
			fastest.OK,
		},
		{
			unpack.User{Name: "ASDF", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			[][]string{
				[]string{
					`{"recenttracks":{"track":[{"artist":{"#text":"X"}}], "@attr":{"page":"1","totalPages":"3"}}}`,
					"", "",
				},
			},
			nil,
			fastest.Fail,
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			files := make(map[rsrc.Locator][]byte)
			for j, day := range tc.data {
				for k, d := range day {
					reg, _ := tc.user.Registered.Midnight()
					time := reg + int64(j*86400)
					files[rsrc.History(tc.user.Name, k+1, rsrc.ToDay(time))] = []byte(d)
				}
			}
			io, _ := mock.IO(files, mock.Path)

			dps, err := LoadHistory(tc.user, tc.until, io)
			ft.Implies(err != nil, tc.err == fastest.Fail, err)
			ft.Implies(err == nil, tc.err == fastest.OK)
			ft.DeepEquals(dps, tc.dps)
		})
	}
}

func TestUpdateHistory(t *testing.T) {
	h0 := rsrc.History("AA", 1, rsrc.ToDay(0*86400))
	h1 := rsrc.History("AA", 1, rsrc.ToDay(1*86400))
	h2 := rsrc.History("AA", 1, rsrc.ToDay(2*86400))
	h3 := rsrc.History("AA", 1, rsrc.ToDay(3*86400))

	testCases := []struct {
		user           unpack.User
		until          rsrc.Day
		saved          []unpack.PlayCount
		tracksFile     map[rsrc.Locator][]byte
		tracksDownload map[rsrc.Locator][]byte
		plays          []unpack.PlayCount
		ok             bool
	}{
		{ // No data
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.PlayCount{},
			false,
		},
		{ // Registration day invalid
			unpack.User{Name: "AA", Registered: rsrc.NoDay()},
			rsrc.ToDay(0),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.PlayCount{},
			false,
		},
		{ // Begin no valid day
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.NoDay(),
			nil,
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.PlayCount{},
			false,
		},
		{ // download one day
			unpack.User{Name: "AA", Registered: rsrc.ToDay(300)}, // registered at 0:05
			rsrc.ToDay(0),
			[]unpack.PlayCount{},
			map[rsrc.Locator][]byte{h0: nil},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[]unpack.PlayCount{
				unpack.PlayCount{"ASDF": 1},
			},
			true,
		},
		{ // download some, have some
			unpack.User{Name: "AA", Registered: rsrc.ToDay(86400)},
			rsrc.ToDay(3 * 86400),
			[]unpack.PlayCount{
				unpack.PlayCount{"XX": 4},
				unpack.PlayCount{}, // will be overwritten
			},
			map[rsrc.Locator][]byte{
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h2: []byte(`{"recenttracks":{"track":[], "@attr":{"totalPages":"1"}}}`),
				h3: nil,
			},
			map[rsrc.Locator][]byte{
				h1: nil,
				h2: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"ASDF"}}], "@attr":{"totalPages":"1"}}}`),
				h3: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"B"}}], "@attr":{"totalPages":"1"}}}`),
			},
			[]unpack.PlayCount{
				unpack.PlayCount{"XX": 4},
				unpack.PlayCount{"ASDF": 1},
				unpack.PlayCount{"B": 1},
			},
			true,
		},
		{ // have more than want
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(86400),
			[]unpack.PlayCount{
				unpack.PlayCount{"XX": 2},
				unpack.PlayCount{"A": 1},
				unpack.PlayCount{"DropMe": 1},
				unpack.PlayCount{"DropMeToo": 100},
			},
			map[rsrc.Locator][]byte{
				h0: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"XX"}},{"artist":{"#text":"XX"}}], "@attr":{"totalPages":"1"}}}`),
				h1: []byte(`{"recenttracks":{"track":[{"artist":{"#text":"A"}}], "@attr":{"totalPages":"1"}}}`),
			},
			map[rsrc.Locator][]byte{},
			[]unpack.PlayCount{
				unpack.PlayCount{"XX": 2},
				unpack.PlayCount{"A": 1},
			},
			true,
		},
		{ // download error
			unpack.User{Name: "AA", Registered: rsrc.ToDay(0)},
			rsrc.ToDay(0),
			[]unpack.PlayCount{},
			map[rsrc.Locator][]byte{},
			map[rsrc.Locator][]byte{},
			[]unpack.PlayCount{},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tc.tracksFile[rsrc.AllDayPlays(tc.user.Name)] = nil
			io1, _ := mock.IO(tc.tracksFile, mock.Path)
			if tc.saved != nil {
				if err := unpack.WriteAllDayPlays(tc.saved, tc.user.Name, io1); err != nil {
					t.Error("unexpected error during write of all day plays:", err)
				}

			}

			io0, _ := mock.IO(tc.tracksDownload, mock.URL)

			pool, _ := store.NewCache([][]rsrc.IO{[]rsrc.IO{io0}, []rsrc.IO{io1}})

			plays, err := UpdateHistory(&tc.user, tc.until, pool)
			if err != nil && tc.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !tc.ok {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				if !reflect.DeepEqual(plays, tc.plays) {
					t.Errorf("updated plays faulty:\nhas:      %v\nexpected: %v",
						plays, tc.plays)
				}
			}
		})
	}
}
