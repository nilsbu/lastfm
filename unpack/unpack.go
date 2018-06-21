package unpack

import (
	"github.com/nilsbu/lastfm/rsrc"
)

// User contains relevant core information about a user.
type User struct {
	Name       rsrc.Name
	Registered rsrc.Day
}

// GetUser returns the name and regestering date of a user.
func GetUser(ui *UserInfo) *User {
	utc := ui.User.Registered.UTC
	return &User{
		rsrc.Name(ui.User.Name),
		rsrc.ToDay(utc)}
}

// DayPlays lists the number of plays for a set of artists in a given day.
// TODO find a place
type DayPlays map[rsrc.Name]int

// GetTracksPages returns the total number of pages declared in urt.
func GetTracksPages(urt *UserRecentTracks) (page int) {
	return urt.RecentTracks.Attr.TotalPages
}

//CountPlays counts the number of plays per artist in urt.
func CountPlays(urt *UserRecentTracks) DayPlays {
	dp := make(DayPlays)
	for _, track := range urt.RecentTracks.Track {
		dp[rsrc.Name(track.Artist.Str)]++
	}
	return dp
}
