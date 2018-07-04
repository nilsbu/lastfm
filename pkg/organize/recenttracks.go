package organize

import (
	"errors"

	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type HistoryDay map[string]int

// LoadAllDayPlays load plays from all days since the registration of the user.
// TODO consistant naming scheme between read, write, download, load, get etc.
func LoadAllDayPlays(
	user unpack.User,
	until rsrc.Day,
	r rsrc.Reader) ([]HistoryDay, error) {

	untilMdn, uOK := until.Midnight()
	registered, rOK := user.Registered.Midnight()
	if !uOK {
		return nil, errors.New("parameter 'until' is no valid Day")
	} else if !rOK {
		return nil, errors.New("user has no valid registration date")
	}
	days := int((untilMdn - registered) / 86400)
	result := make([]HistoryDay, days+1)
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
	DayPlays map[string]int
	Err      error
}

func loadDayPlays(
	user string,
	time rsrc.Day,
	r rsrc.Reader,
) (map[string]int, error) {
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
			histPage.Plays[k] += v
		}
	}
	close(back)

	return histPage.Plays, err
}
