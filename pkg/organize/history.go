package organize

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// LoadHistory load plays from all days since the registration of the user.
func LoadHistory(
	user unpack.User,
	until rsrc.Day,
	r rsrc.Reader) ([]map[string]float64, error) {

	if until == nil {
		return nil, errors.New("parameter 'until' is no valid Day")
	} else if user.Registered == nil {
		return nil, errors.New("user has no valid registration date")
	}

	registered := user.Registered.Midnight()

	days := int((until.Midnight() - registered) / 86400)
	result := make([]map[string]float64, days+1)
	feedback := make(chan error)
	for i := range result {
		go func(i int) {
			date := user.Registered.AddDate(0, 0, i)
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
	DayPlays map[string]float64
	Err      error
}

func loadDayPlays(
	user string,
	time rsrc.Day,
	r rsrc.Reader,
) (map[string]float64, error) {
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
				histPage.Plays[k] += v
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
	until rsrc.Day, // TODO change to end/before
	s store.Store,
) (plays []map[string]float64, err error) {
	if user.Registered == nil {
		return nil, fmt.Errorf("user '%v' has no valid registration date",
			user.Name)
	}
	registeredDay := user.Registered.Midnight()
	begin := registeredDay

	oldPlays, err := unpack.LoadAllDayPlays(user.Name, s)
	if err != nil {
		oldPlays = []map[string]float64{}
	} else if len(oldPlays) > 0 {
		// TODO cleanup the use of time in this function
		begin = user.Registered.AddDate(0, 0, len(oldPlays)-1).Midnight()
	}

	bookmark, err := unpack.LoadBookmark(user.Name, s)
	if err == nil && bookmark.Midnight() < begin {
		begin = bookmark.Midnight()
	}

	if len(oldPlays) > 0 {
		days := int((begin - registeredDay) / 86400)
		oldPlays = oldPlays[:days]
	}

	if until == nil {
		return nil, errors.New("'until' is not a valid day")
	}

	if begin > until.Midnight()+86400 {
		days := int((begin-registeredDay)/86400) - 1
		return oldPlays[:days], nil
	}

	newPlays, err := LoadHistory(
		unpack.User{Name: user.Name, Registered: rsrc.ToDay(begin)},
		until, store.Fresh(s))

	return append(oldPlays, newPlays...), err
}
