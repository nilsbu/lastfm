package rsrc

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/nilsbu/lastfm/pkg/fail"
)

type Locator interface {
	URL(apiKey string) (string, error)
	Path() (string, error)
}

// TODO docu
type lastFM struct {
	method   string
	nameType string
	name     string
	page     int
	day      Day
	limit    int
}

// UserInfo returens a locator for the Last.fm API call "user.getInfo". if the
// user name is malformed, it returns a critical error.
func UserInfo(user string) (*lastFM, error) {
	if err := checkUserName(user); err != nil {
		return nil, err
	}
	return &lastFM{
		method:   "user.getInfo",
		nameType: "user",
		name:     user,
		page:     -1,
		day:      NoDay(),
		limit:    -1,
	}, nil
}

func checkUserName(user string) error {
	if len(user) < 2 {
		return fail.WrapError(fail.Critical,
			fmt.Errorf("user name '%v' too short, min length is 2", user))
	} else if len(user) > 15 {
		return fail.WrapError(fail.Critical,
			fmt.Errorf("user name '%v' too long, max length is 15", user))
	} else if !isLetter(rune(user[0])) {
		return fail.WrapError(fail.Critical,
			fmt.Errorf("user name '%v' doesn't begin with a character", user))
	}

	for _, char := range user[1:] {
		switch {
		case rune(char) == rune('-') || rune(char) == rune('_'):
		case rune(char) >= rune('0') && rune(char) <= rune('9'):
		case isLetter(char):
		default:
			return fail.WrapError(fail.Critical,
				fmt.Errorf("user name contains invalid character '%v'", string(char)))
		}
	}
	return nil
}

func isLetter(char rune) bool {
	if rune(char) >= rune('A') && rune(char) <= rune('Z') {
		return true
	} else if rune(char) >= rune('a') && rune(char) <= rune('z') {
		return true
	}
	return false
}

func History(user string, page int, day Day) (*lastFM, error) {
	if err := checkUserName(user); err != nil {
		return nil, err
	} else if page <= 0 {
		return nil, fail.WrapError(fail.Critical,
			fmt.Errorf("page number must be positive, was %v", page))
	} else if _, ok := day.Midnight(); !ok {
		return nil, fail.WrapError(fail.Critical,
			errors.New("invalid day, must have positive midnight"))
	}

	return &lastFM{
		method:   "user.getRecentTracks",
		nameType: "user",
		name:     user,
		page:     page,
		day:      day,
		limit:    200,
	}, nil
}

func (loc *lastFM) URL(apiKey string) (string, error) {
	if err := checkAPIKey(apiKey); err != nil {
		return "", err
	}
	base := "http://ws.audioscrobbler.com/2.0/"
	params := "?format=json&api_key=%v&method=%v&%v=%v"

	url := base + fmt.Sprintf(params, apiKey,
		loc.method, loc.nameType, url.PathEscape(string(loc.name)))

	if loc.page > 0 {
		url += fmt.Sprintf("&page=%d", int(loc.page))
	}

	if timestamp, ok := loc.day.Midnight(); ok {
		url += fmt.Sprintf("&from=%d&to=%d",
			timestamp-1, timestamp+86400)
	}

	if loc.limit > -1 {
		url += fmt.Sprintf("&limit=%d", loc.limit)
	}

	return url, nil
}

func checkAPIKey(apiKey string) error {
	if len(apiKey) != 32 {
		return fail.WrapError(fail.Critical,
			errors.New("API key does not have length 32"))
	}

	for _, char := range apiKey[1:] {
		switch {
		case rune(char) >= rune('a') && rune(char) <= rune('z'):
		case rune(char) >= rune('0') && rune(char) <= rune('9'):
		default:
			return fail.WrapError(fail.Critical,
				fmt.Errorf("user name contains invalid character '%v'", string(char)))
		}
	}

	return nil
}

func (loc *lastFM) Path() (string, error) {
	path := fmt.Sprintf(".lastfm/data/%v/%v",
		loc.method, parseForPath(loc.name))

	if timestamp, ok := loc.day.Midnight(); ok {
		path += fmt.Sprintf(".%d", timestamp)
	}
	if loc.page > 0 {
		path += fmt.Sprintf("(%v)", loc.page)
	}

	return path + ".json", nil
}

func escapeBadNames(name string) string {
	bad := [13]string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4",
		"LPT1", "LPT2", "LPT3", "LPT4", "LST"}

	upperName := strings.ToUpper(string(name))
	for _, s := range bad {
		if upperName == s {
			return "_" + name
		}
	}

	return name
}

func parseForPath(name string) string {
	escaped := url.PathEscape(string(name))
	escaped = strings.Replace(escaped, "%20", "+", -1)
	escaped = strings.Replace(escaped, "/", "+", -1)
	return escapeBadNames(escaped)
}

// TODO docu
type util struct {
	method string
	public bool
}

func Supertags() *util {
	return &util{
		method: "supertags",
		public: true,
	}
}

func APIKey() *util {
	return &util{
		method: "apikey",
		public: false,
	}
}

func SessionID() *util {
	return &util{
		method: "sessionid",
		public: false,
	}
}

func (u util) URL(apiKey string) (string, error) {
	return "", fail.WrapError(fail.Control,
		fmt.Errorf("'%v' cannot be used as a URL", u.method))
}

func (u util) Path() (string, error) {
	if u.public {
		return fmt.Sprintf("data/util/%v.json", u.method), nil
	}
	return fmt.Sprintf(".lastfm/util/%v.json", u.method), nil
}

type userData struct {
	method string
	name   string
}

func AllDayPlays(user string) (*userData, error) {
	if err := checkUserName(user); err != nil {
		return nil, err
	}
	return &userData{
		method: "alldayplays",
		name:   user,
	}, nil
}

func Bookmark(user string) (*userData, error) {
	if err := checkUserName(user); err != nil {
		return nil, err
	}
	return &userData{
		method: "bookmark",
		name:   user,
	}, nil
}

func (u userData) URL(apiKey string) (string, error) {
	return "", fail.WrapError(fail.Control,
		fmt.Errorf("'%v' cannot be used as a URL", u.method))
}

func (u userData) Path() (string, error) {
	return fmt.Sprintf(".lastfm/user/%v/%v.json", u.name, u.method), nil
}
