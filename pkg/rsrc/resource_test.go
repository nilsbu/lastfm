package rsrc

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/fail"
)

func TestLastFMURL(t *testing.T) {
	base := "http://ws.audioscrobbler.com/2.0/?format=json&"
	cases := []struct {
		loc    Locator
		apiKey string
		url    string
		ok     bool
	}{
		{ // ok
			UserInfo("user1"), "a3ee123098128acf29ca9f0cf29ca9f0",
			base + "api_key=a3ee123098128acf29ca9f0cf29ca9f0&method=user.getInfo&user=user1",
			true,
		},
		{ // ok
			ArtistInfo("König"), "a3ee123098128acf29ca9f0cf29ca9f0",
			base + "api_key=a3ee123098128acf29ca9f0cf29ca9f0&method=artist.getInfo&artist=K%C3%B6nig",
			true,
		},
		{ // ok
			ArtistTags("dido"), "a3ee123098128acf29ca9f0cf29ca9f0",
			base + "api_key=a3ee123098128acf29ca9f0cf29ca9f0&method=artist.getTopTags&artist=dido",
			true,
		},
		{ // ok
			TagInfo("blub"), "a3ee123098128acf29ca9f0cf29ca9f0",
			base + "api_key=a3ee123098128acf29ca9f0cf29ca9f0&method=tag.getInfo&tag=blub",
			true,
		},
		{ // invalid API key
			UserInfo("user1"), "a3ee",
			"",
			false,
		},
		{ // invalid API key
			UserInfo("user1"), "a3ee1NON-HEX-CHARACTERS0c29ca9f0",
			"",
			false,
		},
		{ // ok
			History("abc", 1, ToDay(86400)), "a3ee123098128acf29ca9f0cf29ca9f0",
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
	cases := []struct {
		loc  Locator
		path string
		// path is always ok, since input is considered valid
	}{
		{
			UserInfo("user1"),
			".lastfm/data/user.getInfo/user1.json",
		},
		{ // name must be escaped for Windows
			UserInfo("aux"),
			".lastfm/data/user.getInfo/_aux.json",
		},
		{
			History("abc", 1, ToDay(86400)),
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
		{SessionInfo(), ".lastfm/util/session.json"},
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

func TestUserDataURL(t *testing.T) {
	_, err := AllDayPlays("user1").URL("a3ee123098128acf29ca9f0cf29ca9f0")
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
	cases := []struct {
		loc  Locator
		path string
		// path is always ok, since input is considered valid
	}{
		{AllDayPlays("user1"), ".lastfm/user/user1/alldayplays.json"},
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
