package unpack

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type obtainer interface {
	locator() rsrc.Locator
	deserializer() interface{}
	interpret(raw interface{}) (interface{}, error)
}

func obtain(o obtainer, r rsrc.Reader) (interface{}, error) {
	data, err := r.Read(o.locator())
	if err != nil {
		return nil, err
	}

	raw := o.deserializer()
	err = json.Unmarshal(data, raw)
	if err != nil {
		return nil, errors.Wrap(err, "could not deserialize")
	}

	return o.interpret(raw)
}

// User contains relevant core information about a user.
type User struct {
	Name       string
	Registered rsrc.Day
}

// DayPlays lists the number of plays for a set of artists in a given day.
// TODO find a place
type DayPlays map[string]int

// GetTracksPages returns the total number of pages declared in urt.
func GetTracksPages(urt *UserRecentTracks) (page int) {
	return urt.RecentTracks.Attr.TotalPages
}

//CountPlays counts the number of plays per artist in urt.
func CountPlays(urt *UserRecentTracks) DayPlays {
	dp := make(DayPlays)
	for _, track := range urt.RecentTracks.Track {
		if !track.Attr.NowPlaying {
			dp[track.Artist.Str]++
		}
	}
	return dp
}
