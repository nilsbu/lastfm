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
	r rsrc.Reader) ([][]charts.Song, error) {

	if until == nil {
		return nil, errors.New("parameter 'until' is no valid Day")
	} else if user.Registered == nil {
		return nil, errors.New("user has no valid registration date")
	}

	registered := user.Registered.Midnight()

	days := int((until.Midnight() - registered) / 86400)
	result := make([][]charts.Song, days+1)
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

// loadDayPlaysResult is the result of loadDayPlays.
type loadDayPlaysResult struct {
	DayPlays []charts.Song
	Page     int
	Err      error
}

func loadDayPlays(
	user string,
	time rsrc.Day,
	r rsrc.Reader,
) ([]charts.Song, error) {
	firstPage, err := unpack.LoadHistoryDayPage(user, 1, time, r)
	if err != nil {
		return nil, err
	} else if firstPage.Pages == 1 {
		return firstPage.Plays, nil
	}

	pages := make([][]charts.Song, firstPage.Pages)
	pages[0] = firstPage.Plays

	back := make(chan loadDayPlaysResult)
	for page := 2; page <= len(pages); page++ {
		go func(page int) {
			histPage, tmpErr := unpack.LoadHistoryDayPage(user, page, time, r)
			if tmpErr != nil {
				back <- loadDayPlaysResult{nil, page, tmpErr}
			} else {
				back <- loadDayPlaysResult{histPage.Plays, page, nil}
			}
		}(page)
	}

	for p := 1; p < len(pages); p++ {
		dpr := <-back
		if dpr.Err != nil {
			return nil, dpr.Err
		}

		pages[dpr.Page-1] = dpr.DayPlays
	}
	close(back)

	plays := []charts.Song{}
	for _, page := range pages {
		plays = append(plays, page...)
	}
	return plays, err
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

	summed := []map[string]float64{}
	for _, day := range newPlays {
		page := map[string]float64{}
		for _, song := range day {
			if _, ok := page[song.Artist]; ok {
				page[song.Artist]++
			} else {
				page[song.Artist] = 1
			}
		}
		summed = append(summed, page)
	}

	return append(oldPlays, summed...), err
}
