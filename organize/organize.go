package organize

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/unpack"
)

// TODO name / what is this file

// LoadAPIKey loads an the API key.
func LoadAPIKey(r io.Reader) (key io.APIKey, err error) {
	rsrc := io.NewAPIKey()
	data, err := r.Read(rsrc)
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

	return io.APIKey(unm.Key), nil
}

// WriteAllDayPlays writes a list of day plays.
func WriteAllDayPlays(
	plays []unpack.DayPlays,
	name io.Name,
	w io.Writer) (err error) {
	jsonData, _ := json.Marshal(plays)
	err = w.Write(jsonData, io.NewAllDayPlays(name))
	return
}

// ReadAllDayPlays reads a list of day plays.
func ReadAllDayPlays(
	name io.Name,
	r io.Reader) (plays []unpack.DayPlays, err error) {
	jsonData, err := r.Read(io.NewAllDayPlays(name))
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &plays)
	return
}

// ReadBookmark read a bookmark for a user's saved daily plays.
// TODO Bookmarks should use time.Time
func ReadBookmark(user io.Name, r io.Reader) (utc int64, err error) {
	data, err := r.Read(io.NewBookmark(user))
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
func WriteBookmark(utc int64, user io.Name, w io.Writer) error {
	bookmark := unpack.Bookmark{
		UTC:        utc,
		TimeString: time.Unix(utc, 0).UTC().Format("2006-01-02 15:04:05 +0000 UTC"),
	}

	data, _ := json.Marshal(bookmark)
	err := w.Write(data, io.NewBookmark(user))
	return err
}

// UpdateAllDayPlays loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateAllDayPlays(
	user unpack.User,
	until io.Midnight,
	ioPool io.Pool, // Need Wrapper for Async readers ??
) (plays []unpack.DayPlays, err error) {
	registeredDay := user.Registered - user.Registered%86400
	begin := registeredDay
	fr := io.SeqReader(ioPool.ReadFile)

	oldPlays, err := ReadAllDayPlays(user.Name, fr)
	if err != nil {
		oldPlays = []unpack.DayPlays{}
	} else if len(oldPlays) > 0 {
		begin = registeredDay + io.Midnight(86400*(len(oldPlays)-1))
		oldPlays = oldPlays[:len(oldPlays)-1]
	}

	if begin > until+86400 {
		days := int((begin-registeredDay)/86400) - 1
		return oldPlays[:days], nil
	}

	newPlays, err := LoadAllDayPlays(
		unpack.User{Name: user.Name, Registered: begin},
		until, io.AsyncDownloadGetter(ioPool))

	return append(oldPlays, newPlays...), err
}
