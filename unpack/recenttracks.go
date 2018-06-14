package unpack

import "github.com/nilsbu/lastfm/io"

// DayPlays lists the number of plays for a set of artists in a given day.
// TODO find a place
type DayPlays map[io.Name]int

// GetTracksPages returns the total number of pages declared in urt.
func GetTracksPages(urt *UserRecentTracks) (page int) {
	return urt.RecentTracks.Attr.TotalPages
}

//CountPlays counts the number of plays per artist in urt.
func CountPlays(urt *UserRecentTracks) DayPlays {
	dp := make(DayPlays)
	for _, track := range urt.RecentTracks.Track {
		dp[io.Name(track.Artist.Str)]++
	}
	return dp
}
