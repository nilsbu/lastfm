package unpack

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// User contains relevant core information about a user.
type User struct {
	Name       string
	Registered rsrc.Day
}

type obUserInfo struct {
	name string
}

// LoadUserInfo loads a user's registration date. It is returned along with the
// name.
func LoadUserInfo(name string, r rsrc.Reader) (*User, error) {
	data, err := obtain(&obUserInfo{name}, r)
	if err != nil {
		return nil, err
	}
	user := data.(*User)
	return user, nil
}

func (o *obUserInfo) locator() rsrc.Locator {
	return rsrc.UserInfo(o.name)
}

func (o *obUserInfo) deserializer() interface{} {
	return &jsonUserInfo{}
}

func (o *obUserInfo) interpret(raw interface{}) (interface{}, error) {
	ui := raw.(*jsonUserInfo)

	utc := ui.User.Registered.UTC
	return &User{ui.User.Name, rsrc.ToDay(utc)}, nil
}

// PlayCount assigns artists a play count.
type PlayCount map[string]int

// HistoryDayPage is a single page of a day of a user's played tracks.
type HistoryDayPage struct {
	Plays PlayCount
	Pages int
}

type obHistory struct {
	user string
	page int
	day  rsrc.Day
}

// LoadHistoryDayPage loads a page of a user's played tracks.
func LoadHistoryDayPage(
	user string, page int, day rsrc.Day, r rsrc.Reader) (*HistoryDayPage, error) {
	data, err := obtain(&obHistory{user, page, day}, r)
	if err != nil {
		return nil, err
	}
	hist := data.(*HistoryDayPage)
	return hist, nil
}

func (o *obHistory) locator() rsrc.Locator {
	return rsrc.History(o.user, o.page, o.day)
}

func (o *obHistory) deserializer() interface{} {
	return &jsonUserRecentTracks{}
}

func (o *obHistory) interpret(raw interface{}) (interface{}, error) {
	data := raw.(*jsonUserRecentTracks)

	return &HistoryDayPage{
		countPlays(data),
		data.RecentTracks.Attr.TotalPages}, nil
}

func countPlays(urt *jsonUserRecentTracks) PlayCount {
	dp := make(PlayCount)
	for _, track := range urt.RecentTracks.Track {
		if !track.Attr.NowPlaying {
			dp[track.Artist.Str]++
		}
	}
	return dp
}
