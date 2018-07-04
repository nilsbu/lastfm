package organize

import (
	"encoding/json"
	"errors"
	"fmt"

	errs "github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// TODO name / what is this file

// LoadAPIKey loads an the API key.
func LoadAPIKey(r rsrc.Reader) (apiKey string, err error) {
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

	return unm.Key, nil
}

// SessionID describes a session.
type SessionID string // TODO should be struct or

// LoadSessionID loads a session ID.
func LoadSessionID(r rsrc.Reader) (SessionID, error) {
	data, err := r.Read(rsrc.SessionID())
	if err != nil {
		return "", err
	}

	unm := &unpack.SessionID{}
	err = json.Unmarshal(data, unm)
	if err != nil {
		return "", err
	}
	if unm.User == "" {
		return "", errors.New("No valid user was read")
	}

	return SessionID(unm.User), nil
}

// WriteAllDayPlays writes a list of day plays.
func WriteAllDayPlays(
	plays []unpack.DayPlays,
	user string,
	w rsrc.Writer) (err error) {
	jsonData, _ := json.Marshal(plays)

	loc, err := rsrc.AllDayPlays(user)
	if err != nil {
		return err
	}
	return w.Write(jsonData, loc)
}

// ReadAllDayPlays reads a list of day plays.
func ReadAllDayPlays(
	user string,
	r rsrc.Reader) (plays []unpack.DayPlays, err error) {
	loc, err := rsrc.AllDayPlays(user)
	if err != nil {
		return nil, err
	}
	jsonData, err := r.Read(loc)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &plays)
	return
}

// UpdateAllDayPlays loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateAllDayPlays(
	user *unpack.User,
	until rsrc.Day,
	s store.Store,
) (plays []unpack.DayPlays, err error) {
	registeredDay, ok := user.Registered.Midnight()
	if !ok {
		return nil, fmt.Errorf("user '%v' has no valid registration date",
			user.Name)
	}
	begin := registeredDay

	oldPlays, err := ReadAllDayPlays(user.Name, s)
	if err != nil {
		oldPlays = []unpack.DayPlays{}
	} else if len(oldPlays) > 0 {
		begin = registeredDay + int64(86400*(len(oldPlays)-1))
		oldPlays = oldPlays[:len(oldPlays)-1]
	}

	midn, ok := until.Midnight()
	if !ok {
		return nil, errors.New("'until' is not a valid day")
	}
	if begin > midn+86400 {
		days := int((begin-registeredDay)/86400) - 1
		return oldPlays[:days], nil
	}

	newPlays, err := LoadAllDayPlays(
		unpack.User{Name: user.Name, Registered: rsrc.ToDay(begin)},
		until, store.NewUpToDate(s))

	return append(oldPlays, newPlays...), err
}

func LoadUser(user string, r rsrc.Reader) (*unpack.User, error) {
	name, err := rsrc.UserInfo(user)
	if err != nil {
		return nil, err
	}

	data, err := r.Read(name)
	if err != nil {
		return nil, err
	}

	userRaw := unpack.UserInfo{}
	err = json.Unmarshal(data, &userRaw)
	if err != nil {
		return nil, err
	}

	return unpack.GetUser(&userRaw), nil
}

type TagCount struct {
	Name  string
	Count int
}

func ReadArtistTags(artist string, r rsrc.Reader) ([]TagCount, error) {
	rs, err := rsrc.ArtistTags(artist)
	if err != nil {
		return nil, err
	}

	data, err := r.Read(rs)
	if err != nil {
		return nil, err
	}

	at := unpack.ArtistTags{}
	err = json.Unmarshal(data, &at)
	if err != nil {
		return nil, errs.Wrap(err, "")
	}

	len := len(at.TopTags.Tags)
	if len == 0 {
		return nil, fmt.Errorf("no tags were read for '%v'", artist)
	}

	tags := make([]TagCount, len)
	for i, tag := range at.TopTags.Tags {
		tags[i] = TagCount{Name: tag.Name, Count: tag.Count}
	}

	return tags, nil
}
