package organize

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/rsrc"
	"github.com/nilsbu/lastfm/unpack"
)

// TODO name / what is this file

// LoadAPIKey loads an the API key.
func LoadAPIKey(r io.Reader) (apiKey rsrc.Key, err error) {
	data, err := r.Read(rsrc.APIKey())
	if err != nil {
		return
	}

	unm := &unpack.APIKey{}
	err = json.Unmarshal(data, unm)
	if err != nil {
		return
	}
	if unm.Key == "" {
		return "", errors.New("No valid API key was read")
	}

	return rsrc.Key(unm.Key), nil
}

// WriteAllDayPlays writes a list of day plays.
func WriteAllDayPlays(
	plays []unpack.DayPlays,
	name rsrc.Name,
	w io.Writer) (err error) {
	jsonData, _ := json.Marshal(plays)

	rs, err := rsrc.AllDayPlays(name)
	if err != nil {
		return err
	}
	return w.Write(jsonData, rs)
}

// ReadAllDayPlays reads a list of day plays.
func ReadAllDayPlays(
	name rsrc.Name,
	r io.Reader) (plays []unpack.DayPlays, err error) {
	rs, err := rsrc.AllDayPlays(name)
	if err != nil {
		return nil, err
	}
	jsonData, err := r.Read(rs)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &plays)
	return
}

// ReadBookmark read a bookmark for a user's saved daily plays.
// TODO Bookmarks should use time.Time
func ReadBookmark(user rsrc.Name, r io.Reader) (utc int64, err error) {
	rs, err := rsrc.Bookmark(user)
	if err != nil {
		return 0, err
	}
	data, err := r.Read(rs)
	if err != nil {
		return 0, err
	}

	bookmark := &unpack.Bookmark{}
	err = json.Unmarshal(data, bookmark)
	if err != nil {
		return 0, err
	}

	return bookmark.UTC, nil
}

// WriteBookmark writes a bookmark for a user's saved daily plays.
func WriteBookmark(utc int64, user rsrc.Name, w io.Writer) error {
	bookmark := unpack.Bookmark{
		UTC:        utc,
		TimeString: time.Unix(utc, 0).UTC().Format("2006-01-02 15:04:05 +0000 UTC"),
	}

	data, _ := json.Marshal(bookmark)
	rs, err := rsrc.Bookmark(user)
	if err != nil {
		return err
	}
	err = w.Write(data, rs)
	return err
}
