package unpack

import (
	"fmt"
	"testing"

	"github.com/nilsbu/fastest"
)

func TestGetTracksPages(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		urt   *UserRecentTracks
		pages int
	}{
		{&UserRecentTracks{RecentTracks: recentTracks{Attr: recentTracksAttr{TotalPages: 3}}}, 3},
	}

	for i, tc := range testCases {
		s := fmt.Sprintf("#%v", i)
		ft.Seq(s, func(ft fastest.T) {
			pages := GetTracksPages(tc.urt)

			ft.Equals(pages, tc.pages)
		})
	}
}

func TestCountPlays(t *testing.T) {
	ft := fastest.T{T: t}

	testCases := []struct {
		urt UserRecentTracks
		dp  DayPlays
	}{
		{
			UserRecentTracks{},
			make(DayPlays),
		},
		{
			UserRecentTracks{
				RecentTracks: recentTracks{
					Track: []track{
						track{Artist: text{Str: "BTS"}},
						track{Artist: text{Str: "XX"}},
						track{Artist: text{Str: "12"}},
						track{Artist: text{Str: "BTS"}},
					},
				},
			},
			DayPlays{"BTS": 2, "XX": 1, "12": 1},
		},
	}

	for i, tc := range testCases {
		ft.Seq(fmt.Sprintf("#%v", i), func(ft fastest.T) {
			dp := CountPlays(&tc.urt)

			ft.DeepEquals(dp, tc.dp)
		})
	}
}
