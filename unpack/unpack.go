package unpack

import "github.com/nilsbu/lastfm/io"

// User contains relevant core information about a user.
type User struct {
	Name       io.Name
	Registered io.Midnight
}

// GetUser returns the name and regestering date of a user.
func GetUser(ui *UserInfo) *User {
	utc := ui.User.Registered.UTC
	return &User{
		io.Name(ui.User.Name),
		io.Midnight(utc - utc%86400)}
}

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
