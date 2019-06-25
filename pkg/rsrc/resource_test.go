package rsrc

import "testing"

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
		{ // ok
			ArtistTags("M&M"), "a3ee123098128acf29ca9f0cf29ca9f0",
			base + "api_key=a3ee123098128acf29ca9f0cf29ca9f0&method=artist.getTopTags&artist=M%26M",
			true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			url, err := c.loc.URL(c.apiKey)
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
	cases := []struct {
		loc  Locator
		path string
		// path is always ok, since input is considered valid
	}{
		{
			UserInfo("user2"),
			".lastfm/raw/user.getInfo/60/25/d18fe48abd45168528f18a82e265dd98d421a7084aa09f61b341703901a3.json",
		},
		{
			UserInfo("user1"),
			".lastfm/raw/user.getInfo/0a/04/1b9462caa4a31bac3567e0b6e6fd9100787db2ab433d96f6d178cabfce90.json",
		},
		{ // name must be escaped for Windows
			UserInfo("aux"),
			".lastfm/raw/user.getInfo/32/1f/68140efca2b301c8c9e9cd67f0e0e3f89a6b24ca923c13bade1ee6552073.json",
		},
		{
			History("abc", 1, ToDay(2*86400)),
			".lastfm/raw/user.getRecentTracks/abc/86400/1970-01-03T00-00-00-1.json",
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
	cases := []struct {
		loc Locator
		// URL is never ok
	}{
		{APIKey()},
		{SessionInfo()},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			if _, err := c.loc.URL("a3ee123098128acf29ca9f0cf29ca9f0"); err == nil {
				t.Error("util resources should not yield a valid URL")
			}
		})
	}
}

func TestUtilPath(t *testing.T) {
	cases := []struct {
		loc  Locator
		path string
		// path is always ok, since input is considered valid
	}{
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
	cases := []struct {
		loc Locator
		// URL is never ok
	}{
		{AllDayPlays("user1")},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			if _, err := c.loc.URL("a3ee123098128acf29ca9f0cf29ca9f0"); err == nil {
				t.Error("user resources should not yield a valid URL")
			}
		})
	}
}

func TestUserDataPath(t *testing.T) {
	cases := []struct {
		loc  Locator
		path string
		// path is always ok, since input is considered valid
	}{
		{AllDayPlays("user1"), ".lastfm/user/user1/alldayplays.json"},
		{ArtistCorrections("user1"), ".lastfm/user/user1/artistcorrections.json"},
		{SupertagCorrections("user1"), ".lastfm/user/user1/supertagcorrections.json"},
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
