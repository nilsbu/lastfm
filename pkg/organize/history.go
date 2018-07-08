package organize

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// LoadHistory load plays from all days since the registration of the user.
func LoadHistory(
	user unpack.User,
	until rsrc.Day,
	r rsrc.Reader) ([]charts.Charts, error) {

	untilMdn, uOK := until.Midnight()
	registered, rOK := user.Registered.Midnight()
	if !uOK {
		return nil, errors.New("parameter 'until' is no valid Day")
	} else if !rOK {
		return nil, errors.New("user has no valid registration date")
	}
	days := int((untilMdn - registered) / 86400)
	result := make([]charts.Charts, days+1)
	feedback := make(chan error)
	for i := range result {
		go func(i int) {
			date := rsrc.ToDay(int64(i*86400) + registered)
			dp, err := loadDayPlays(user.Name, date, r)
			if err == nil {
				result[i] = dp
			}
			feedback <- err
		}(i)
	}
	fail := []error{}
	for i := 0; i <= days; i++ {
		err := <-feedback
		if err != nil {
			fail = append(fail, err)
		}
	}
	if len(fail) > 0 {
		return nil, fail[0] // TODO return all errors
	}
	return result, nil
}

// LoadDayPlaysResult is the result of loadDayPlays.
type LoadDayPlaysResult struct {
	DayPlays charts.Charts
	Err      error
}

func loadDayPlays(
	user string,
	time rsrc.Day,
	r rsrc.Reader,
) (charts.Charts, error) {
	histPage, err := unpack.LoadHistoryDayPage(user, 1, time, r)
	if err != nil {
		return nil, err
	} else if histPage.Pages == 1 {
		return histPage.Plays, nil
	}

	back := make(chan LoadDayPlaysResult)
	for page := 1; page < histPage.Pages; page++ {
		go func(page int) {
			histPage, tmpErr := unpack.LoadHistoryDayPage(user, page+1, time, r)
			if tmpErr != nil {
				back <- LoadDayPlaysResult{nil, tmpErr}
			} else {
				back <- LoadDayPlaysResult{histPage.Plays, nil}
			}
		}(page)
	}

	for p := 1; p < histPage.Pages; p++ {
		dpr := <-back
		if dpr.Err != nil {
			return nil, dpr.Err
		}

		for k, v := range dpr.DayPlays {
			if _, ok := histPage.Plays[k]; ok {
				histPage.Plays[k][0] += v[0]
			} else {
				histPage.Plays[k] = v
			}
		}
	}
	close(back)

	return histPage.Plays, err
}

// UpdateHistory loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateHistory(
	user *unpack.User,
	until rsrc.Day,
	s store.Store,
) (plays []charts.Charts, err error) {
	registeredDay, ok := user.Registered.Midnight()
	if !ok {
		return nil, fmt.Errorf("user '%v' has no valid registration date",
			user.Name)
	}
	begin := registeredDay

	oldPlays, err := unpack.LoadAllDayPlays(user.Name, s)
	if err != nil {
		oldPlays = []charts.Charts{}
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

	newPlays, err := LoadHistory(
		unpack.User{Name: user.Name, Registered: rsrc.ToDay(begin)},
		until, store.NewUpToDate(s))

	return append(oldPlays, newPlays...), err
}
