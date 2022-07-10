package rsrc

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Locator interface {
	URL(apiKey string) (string, error)
	Path() (string, error)
}

type lastFM struct {
	method   string
	nameType string
	name     string
	track    string
	page     int
	day      Day
	limit    int
}

// UserInfo returns a locator for the Last.fm API call "user.getInfo".
func UserInfo(user string) Locator {
	return &lastFM{
		method:   "user.getInfo",
		nameType: "user",
		name:     user,
		track:    "",
		page:     -1,
		day:      nil,
		limit:    -1,
	}
}

func History(user string, page int, day Day) Locator {
	return &lastFM{
		method:   "user.getRecentTracks",
		nameType: "user",
		name:     user,
		track:    "",
		page:     page,
		day:      day,
		limit:    200,
	}
}

// ArtistInfo returns a locator for the Last.fm API call "artist.getInfo".
func ArtistInfo(artist string) Locator {
	return &lastFM{
		method:   "artist.getInfo",
		nameType: "artist",
		name:     artist,
		track:    "",
		page:     -1,
		day:      nil,
		limit:    -1,
	}
}

// ArtistTags returns a locator for the Last.fm API call "artist.getTopTags".
func ArtistTags(artist string) Locator {
	return &lastFM{
		method:   "artist.getTopTags",
		nameType: "artist",
		name:     artist,
		track:    "",
		page:     -1,
		day:      nil,
		limit:    -1,
	}
}

// TagInfo returns a locator for the Last.fm API call "tag.getInfo".
func TagInfo(tag string) Locator {
	return &lastFM{
		method:   "tag.getInfo",
		nameType: "tag",
		name:     tag,
		track:    "",
		page:     -1,
		day:      nil,
		limit:    -1,
	}
}

// TrackInfo returns a locator for the Last.fm API call "track.getInfo".
func TrackInfo(artist, track string) Locator {
	return &lastFM{
		method:   "track.getInfo",
		nameType: "artist",
		name:     artist,
		track:    track,
		page:     -1,
		day:      nil,
		limit:    -1,
	}
}

func (loc *lastFM) URL(apiKey string) (string, error) {
	if err := CheckAPIKey(apiKey); err != nil {
		return "", err
	}
	base := "http://ws.audioscrobbler.com/2.0/"
	params := "?format=json&api_key=%v&method=%v&%v=%v"

	name := strings.Replace(url.PathEscape(loc.name), "&", "%26", -1)
	track := strings.Replace(url.PathEscape(loc.track), "&", "%26", -1)
	url := base + fmt.Sprintf(params, apiKey,
		loc.method, loc.nameType, name)

	if loc.track != "" {
		url += fmt.Sprintf("&track=%s", track)
	}

	if loc.page > 0 {
		url += fmt.Sprintf("&page=%d", loc.page)
	}

	if loc.day != nil {
		timestamp := loc.day.Midnight()
		url += fmt.Sprintf("&from=%d&to=%d",
			timestamp-1, timestamp+86400)
	}

	if loc.limit > -1 {
		url += fmt.Sprintf("&limit=%d", loc.limit)
	}

	return url, nil
}

// CheckAPIKey checks if an API key is a 32 digit hex string. Letters
// have to be lower case.
func CheckAPIKey(apiKey string) error {
	if len(apiKey) != 32 {
		return errors.New("API key does not have length 32")
	}

	for _, char := range apiKey[1:] {
		switch {
		case rune(char) >= rune('a') && rune(char) <= rune('z'):
		case rune(char) >= rune('0') && rune(char) <= rune('9'):
		default:
			return fmt.Errorf("user name contains invalid character '%v'",
				string(char))
		}
	}

	return nil
}

func (loc *lastFM) Path() (string, error) {
	var path string
	switch loc.method {
	case "user.getRecentTracks":
		t := loc.day.Time()
		path = fmt.Sprintf("%v/%v/%04v-%02v-%02vT%02v-%02v-%02v-%v",
			loc.name, 86400,
			t.Year(), int(t.Month()), t.Day(),
			t.Hour(), t.Minute(), t.Second(),
			loc.page)
	case "track.getInfo":
		key := fmt.Sprintf("%v\t%v", loc.name, loc.track)
		h8 := sha256.Sum256([]byte(key))
		hash := hex.EncodeToString(h8[:])
		path = fmt.Sprintf("%v/%v/%v", hash[0:2], hash[2:4], hash[4:])
	default:
		h8 := sha256.Sum256([]byte(loc.name))
		hash := hex.EncodeToString(h8[:])
		path = fmt.Sprintf("%v/%v/%v", hash[0:2], hash[2:4], hash[4:])
	}

	return fmt.Sprintf(".lastfm/raw/%v/%v.json", loc.method, path), nil
}

// TODO docu
type util struct {
	method string
	public bool
}

func APIKey() Locator {
	return &util{
		method: "apikey",
		public: false,
	}
}

func SessionInfo() Locator {
	return &util{
		method: "session",
		public: false,
	}
}

func (u util) URL(apiKey string) (string, error) {
	return "", fmt.Errorf("'%v' cannot be used as a URL", u.method)
}

func (u util) Path() (string, error) {
	return fmt.Sprintf(".lastfm/util/%v.json", u.method), nil
}

type userData struct {
	method string
	day    Day
	name   string
}

func Bookmark(user string) Locator {
	return &userData{
		method: "bookmark",
		name:   user,
	}
}

func BackupBookmark(user string) Locator {
	return &userData{
		method: "bookmark2",
		name:   user,
	}
}

func AllDayPlays(user string) Locator {
	return &userData{
		method: "alldayplays",
		name:   user,
	}
}

func SongHistory(user string) Locator {
	return &userData{
		method: "history",
		name:   user,
	}
}

func DayHistory(user string, day Day) Locator {
	return &userData{
		method: "days",
		day:    day,
		name:   user,
	}
}

func ArtistCorrections(user string) Locator {
	return &userData{
		method: "artistcorrections",
		name:   user,
	}
}

func SupertagCorrections(user string) Locator {
	return &userData{
		method: "supertagcorrections",
		name:   user,
	}
}

func CountryCorrections(user string) Locator { // TODO shouldn't be user dependent
	return &userData{
		method: "countrycorrections",
		name:   user,
	}
}

func Groups(user string) Locator {
	return &userData{
		method: "groups",
		name:   user,
	}
}

func (u userData) URL(apiKey string) (string, error) {
	return "", fmt.Errorf("'%v' cannot be used as a URL", u.method)
}

func (u userData) Path() (string, error) {
	if u.method == "days" {
		return fmt.Sprintf(".lastfm/user/%v/history/%v.json", u.name, u.day), nil
	}
	return fmt.Sprintf(".lastfm/user/%v/%v.json", u.name, u.method), nil
}
