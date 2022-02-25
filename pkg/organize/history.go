package organize

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// UpdateHistory loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateHistory(
	user *unpack.User,
	until rsrc.Day, // TODO change to end/before
	s store.Store,
) (plays [][]charts.Song, err error) {
	if user.Registered == nil {
		return nil, fmt.Errorf("user '%v' has no valid registration date",
			user.Name)
	}
	registeredDay := user.Registered.Midnight()
	begin := user.Registered

	cache := unpack.NewCached(s)

	bookmark, err := unpack.LoadBookmark(user.Name, s)
	if err == nil {
		begin = bookmark
	}

	oldPlays, err := loadDays(user.Name, user.Registered, begin, s)
	if err != nil {
		return nil, err
	}

	if len(oldPlays) > 0 {
		days := int((begin.Midnight() - registeredDay) / 86400)
		oldPlays = oldPlays[:days]
	}

	if until == nil {
		return nil, errors.New("'until' is not a valid day")
	}

	if begin.Midnight() > until.Midnight()+86400 {
		days := int((begin.Midnight()-registeredDay)/86400) - 1
		return oldPlays[:days], nil
	}

	newPlays, err := loadHistory(
		unpack.User{Name: user.Name, Registered: begin},
		until, store.Fresh(s), cache) // TODO make fresh optional

	for _, plays := range newPlays {
		err = attachDuration(plays, cache)
		if err != nil {
			return nil, err
		}
	}

	return append(oldPlays, newPlays...), err
}

func loadDays(user string, begin, end rsrc.Day, r rsrc.Reader) ([][]charts.Song, error) {
	days := (end.Midnight() - begin.Midnight()) / 86400
	plays := make([][]charts.Song, days)

	err := Pi(int(days), func(i int) error {
		if songs, err := unpack.LoadDayHistory(user, begin.AddDate(0, 0, i), r); err == nil {
			plays[i] = songs
			return nil
		} else {
			return err
		}
	})

	return plays, err
}

// loadHistory load plays from all days since the registration of the user.
func loadHistory(
	user unpack.User,
	until rsrc.Day,
	r rsrc.Reader,
	l unpack.Loader) ([][]charts.Song, error) {

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
			dp, err := loadDayPlays(user.Name, date, r, l)
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
	r rsrc.Reader, cache unpack.Loader,
) ([]charts.Song, error) {
	firstPage, err := unpack.LoadHistoryDayPage(user, 1, time, unpack.NewCacheless(r))
	if err != nil {
		return nil, err
	}
	err = attachDuration(firstPage.Plays, cache)
	if err != nil {
		return nil, err
	}

	if firstPage.Pages < 2 {
		return firstPage.Plays, nil
	}

	pages := make([][]charts.Song, firstPage.Pages)
	pages[0] = firstPage.Plays

	back := make(chan loadDayPlaysResult)
	for page := 2; page <= len(pages); page++ {
		go func(page int) {
			histPage, tmpErr := unpack.LoadHistoryDayPage(user, page, time, unpack.NewCacheless(r))
			if tmpErr != nil {
				back <- loadDayPlaysResult{nil, page, tmpErr}
			} else {
				err := attachDuration(histPage.Plays, cache)
				back <- loadDayPlaysResult{histPage.Plays, page, err}
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

// TODO put somewhere else?

// Pi executes f in parallel n times
func Pi(n int, f func(int) error) error {
	errs := make([]error, n)
	hasError := false

	back := make(chan bool, n)

	for i := 0; i < n; i++ {
		go func(i int) {
			if !hasError {
				if err := f(i); err != nil {
					hasError = true
					errs[i] = err
				}
			}
			back <- true
		}(i)
	}

	for i := 0; i < n; i++ {
		<-back
	}

	if !hasError {
		return nil
	}

	merr := &MultiError{
		Msg:  "error while executing in parallel",
		Errs: []error{},
	}
	for _, err := range errs {
		if err != nil {
			merr.Errs = append(merr.Errs, err)
		}
	}
	return merr
}

func attachDuration(songs []charts.Song, cache unpack.Loader) error {
	return Pi(len(songs), func(i int) error {
		if info, err := unpack.LoadTrackInfo(songs[i].Artist, songs[i].Title, cache); err == nil {
			songs[i].Duration = float64(info.Duration) / 60.0
			return nil
		} else {
			return err
		}
	})
}
