package rsrc

import (
	"testing"
)

func TestUserInfo(t *testing.T) {
	cases := []struct {
		name Name
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
			rs, err := UserInfo(c.name)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Errorf("name '%v' should not have been accepted", c.name)
			}
			if err == nil {
				if rs.name != c.name {
					t.Errorf("got name '%v', expected '%v'", rs.name, c.name)
				}
				if _, ok := rs.day.Midnight(); ok {
					t.Error("must not have a valid midnight")
				}
				if rs.page > 0 {
					t.Error("must not have a valid page")
				}
				// assume other fields without check
			}
		})
	}
}

func TestHistory(t *testing.T) {
	cases := []struct {
		name Name
		page Page
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
			rs, err := History(c.name, c.page, c.day)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Errorf("name '%v' should not have been accepted", c.name)
			}
			if err == nil {
				if rs.name != c.name {
					t.Errorf("got name '%v', expected '%v'", rs.name, c.name)
				}
				if rs.page <= 0 {
					t.Error("page must be positive")
				}
				if rs.page != c.page {
					t.Errorf("got page '%d', expected '%d'", rs.page, c.page)
				}
				if rs.day != c.day {
					rsMid, _ := rs.day.Midnight()
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
		rs     *lastFM
		apiKey Key
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
			url, err := c.rs.URL(c.apiKey)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
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
		rs   *lastFM
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
			path, err := c.rs.Path()

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
	if _, err := APIKey().URL("a3ee123098128acf29ca9f0cf29ca9f0"); err == nil {
		t.Error("util resources should not yield a valid URL")
	}
}

func TestUtilPath(t *testing.T) {
	cases := []struct {
		rs   *util
		path string
		// path is always ok, since input is considered valid
	}{
		{APIKey(), ".lastfm/util/apikey.json"},
		{Supertags(), "data/util/supertags.json"},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			path, err := c.rs.Path()

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
		function func(Name) (*userData, error)
		name     Name
		ok       bool
	}{
		{AllDayPlays, "aA", true},
		{AllDayPlays, "a%", false},
		{Bookmark, "aAasldfhk", true},
		{Bookmark, "a%", false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			rs, err := c.function(c.name)

			if err != nil && c.ok {
				t.Error("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Errorf("name '%v' should not have been accepted", c.name)
			}
			if err == nil {
				if rs.name != c.name {
					t.Errorf("got name '%v', expected '%v'", rs.name, c.name)
				}
				// assume method without check
			}
		})
	}
}

func TestUserDataURL(t *testing.T) {
	allDayPlays, _ := AllDayPlays("user1")
	if _, err := allDayPlays.URL("a3ee123098128acf29ca9f0cf29ca9f0"); err == nil {
		t.Error("user data should not yield a valid URL")
	}
}

func TestUserDataPath(t *testing.T) {
	allDayPlays, _ := AllDayPlays("user1")
	cases := []struct {
		rs   *userData
		path string
		// path is always ok, since input is considered valid
	}{
		{allDayPlays, ".lastfm/user/user1/alldayplays.json"},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			path, err := c.rs.Path()

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
