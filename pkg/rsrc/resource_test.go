package rsrc

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
)

func TestUserInfo(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{"aA", true},                // ok
		{"a%", false},               // forbidden symbol
		{"x", false},                // too short
		{"abcdef0123456789", false}, // too long
		{"0asdf", false},            // first letter no [A-z]
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			loc, err := UserInfo(c.name)

			if err != nil {
				if f, ok := err.(fail.Threat); ok {
					if f.Severity() != fail.Critical {
						t.Error("severity must be 'Critical':", err)
					}
				} else {
					t.Error("error must implement Threat but does not:", err)
				}
				if c.ok {
					t.Error("unexpected error:", err)
				}
			} else if err == nil && !c.ok {
				t.Errorf("name '%v' should not have been accepted", c.name)
			}

			if err == nil {
				if loc.name != c.name {
					t.Errorf("got name '%v', expected '%v'", loc.name, c.name)
				}
				if _, ok := loc.day.Midnight(); ok {
					t.Error("must not have a valid midnight")
				}
				if loc.page > 0 {
					t.Error("must not have a valid page")
				}
				// assume other fields without check
			}
		})
	}
}

func TestHistory(t *testing.T) {
	cases := []struct {
		name string
		page int
		day  Day
		ok   bool
	}{
		{"aA", 1, ToDay(10*86400 + 5000), true}, // ok
		{"1nvalid", 1, ToDay(0), false},         // invalid name
		{"name", 0, ToDay(0), false},            // invalid page
		{"name", 3, NoDay(), false},             // no date
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			loc, err := History(c.name, c.page, c.day)

			if err != nil {
				if f, ok := err.(fail.Threat); ok {
					if f.Severity() != fail.Critical {
						t.Error("severity must be 'Critical':", err)
					}
				} else {
					t.Error("error must implement Threat but does not:", err)
				}
				if c.ok {
					t.Error("unexpected error:", err)
				}
			} else if err == nil && !c.ok {
				t.Errorf("name '%v' should not have been accepted", c.name)
			}

			if err == nil {
				if loc.name != c.name {
					t.Errorf("got name '%v', expected '%v'", loc.name, c.name)
				}
				if loc.page <= 0 {
					t.Error("page must be positive")
				}
				if loc.page != c.page {
					t.Errorf("got page '%d', expected '%d'", loc.page, c.page)
				}
				if loc.day != c.day {
					rsMid, _ := loc.day.Midnight()
					cMid, _ := c.day.Midnight()
					t.Errorf("got midnight '%d', expected '%d'", rsMid, cMid)
				}
			}
		})
	}
}

func TestLastFMURL(t *testing.T) {
	base := "http://ws.audioscrobbler.com/2.0/?format=json&"
	userInfo, _ := UserInfo("user1")
	history, _ := History("abc", 1, ToDay(86400))

	cases := []struct {
		loc    *lastFM
		apiKey string
		url    string
		ok     bool
	}{
		{ // ok
			userInfo, "a3ee123098128acf29ca9f0cf29ca9f0",
			base + "api_key=a3ee123098128acf29ca9f0cf29ca9f0&method=user.getInfo&user=user1",
			true,
		},
		{ // invalid API key
			userInfo, "a3ee",
			"",
			false,
		},
		{ // invalid API key
			userInfo, "a3ee1NON-HEX-CHARACTERS0c29ca9f0",
			"",
			false,
		},
		{ // ok
			history, "a3ee123098128acf29ca9f0cf29ca9f0",
			base + "api_key=a3ee123098128acf29ca9f0cf29ca9f0&method=user.getRecentTracks&user=abc&page=1&from=86399&to=172800&limit=200",
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			url, err := c.loc.URL(c.apiKey)

			if err != nil {
				if f, ok := err.(fail.Threat); ok {
					if f.Severity() != fail.Critical {
						t.Error("severity must be 'Critical':", err)
					}
				} else {
					t.Error("error must implement Threat but does not:", err)
				}
				if c.ok {
					t.Error("unexpected error:", err)
				}
			} else if err == nil && !c.ok {
				t.Errorf("URL() should have thrown an error but did not")
			}

			if err == nil {
				if url != c.url {
					t.Errorf("unexpected url:\n got      '%v',\n expected '%v'",
						url, c.url)
				}

			}
		})
	}
}

func TestLastFMPath(t *testing.T) {
	userInfo, _ := UserInfo("user1")
	badUserInfo, _ := UserInfo("aux")
	history, _ := History("abc", 1, ToDay(86400))

	cases := []struct {
		loc  *lastFM
		path string
		// path is always ok, since input is considered valid
	}{
		{
			userInfo,
			".lastfm/data/user.getInfo/user1.json",
		},
		{ // name must be escaped for Windows
			badUserInfo,
			".lastfm/data/user.getInfo/_aux.json",
		},
		{
			history,
			".lastfm/data/user.getRecentTracks/abc.86400(1).json",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			path, err := c.loc.Path()

			if err != nil {
				t.Error("unexpected error:", err)
			} else {
				if path != c.path {
					t.Errorf("unexpected path:\n got      '%v',\n expected '%v'",
						path, c.path)
				}
			}
		})
	}
}

func TestUtilURL(t *testing.T) {
	_, err := APIKey().URL("a3ee123098128acf29ca9f0cf29ca9f0")
	if err == nil {
		t.Error("util resources should not yield a valid URL")
	}
	if f, ok := err.(fail.Threat); ok {
		if f.Severity() != fail.Control {
			t.Error("severity must be 'Control':", err)
		}
	} else {
		t.Error("error must implement Threat but does not:", err)
	}
}

func TestUtilPath(t *testing.T) {
	cases := []struct {
		loc  *util
		path string
		// path is always ok, since input is considered valid
	}{
		{Supertags(), "data/util/supertags.json"},
		{APIKey(), ".lastfm/util/apikey.json"},
		{SessionID(), ".lastfm/util/sessionid.json"},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			path, err := c.loc.Path()

			if err != nil {
				t.Error("unexpected error:", err)
			} else {
				if path != c.path {
					t.Errorf("unexpected path:\n got      '%v',\n expected '%v'",
						path, c.path)
				}
			}
		})
	}
}

func TestUserDataConstructors(t *testing.T) {
	cases := []struct {
		function func(string) (*userData, error)
		name     string
		ok       bool
	}{
		{AllDayPlays, "aA", true},
		{AllDayPlays, "a%", false},
		{Bookmark, "aAasldfhk", true},
		{Bookmark, "a%", false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			loc, err := c.function(c.name)

			if err != nil {
				if f, ok := err.(fail.Threat); ok {
					if f.Severity() != fail.Critical {
						t.Error("severity must be 'Critical':", err)
					}
				} else {
					t.Error("error must implement Threat but does not:", err)
				}
				if c.ok {
					t.Error("unexpected error:", err)
				}
			} else if err == nil && !c.ok {
				t.Errorf("name '%v' should not have been accepted", c.name)
			}
			if err == nil {
				if loc.name != c.name {
					t.Errorf("got name '%v', expected '%v'", loc.name, c.name)
				}
				// assume method without check
			}
		})
	}
}

func TestUserDataURL(t *testing.T) {
	allDayPlays, _ := AllDayPlays("user1")
	_, err := allDayPlays.URL("a3ee123098128acf29ca9f0cf29ca9f0")
	if err == nil {
		t.Error("user data should not yield a valid URL")
	}
	if f, ok := err.(fail.Threat); ok {
		if f.Severity() != fail.Control {
			t.Error("severity must be 'Control':", err)
		}
	} else {
		t.Error("error must implement Threat but does not:", err)
	}
}

func TestUserDataPath(t *testing.T) {
	allDayPlays, _ := AllDayPlays("user1")
	cases := []struct {
		loc  *userData
		path string
		// path is always ok, since input is considered valid
	}{
		{allDayPlays, ".lastfm/user/user1/alldayplays.json"},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			path, err := c.loc.Path()

			if err != nil {
				t.Error("unexpected error:", err)
			} else {
				if path != c.path {
					t.Errorf("unexpected path:\n got      '%v',\n expected '%v'",
						path, c.path)
				}
			}
		})
	}
}
