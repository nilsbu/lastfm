package organize

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// TODO name / what is this file

// UpdateAllDayPlays loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateAllDayPlays(
	user *unpack.User,
	until rsrc.Day,
	s store.Store,
) (plays []unpack.PlayCount, err error) {
	registeredDay, ok := user.Registered.Midnight()
	if !ok {
		return nil, fmt.Errorf("user '%v' has no valid registration date",
			user.Name)
	}
	begin := registeredDay

	oldPlays, err := unpack.LoadAllDayPlays(user.Name, s)
	if err != nil {
		oldPlays = []unpack.PlayCount{}
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
