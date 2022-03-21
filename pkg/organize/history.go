package organize

import (
	"errors"
	"fmt"

	async "github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// UpdateHistory loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateHistory(
	user *unpack.User,
	until rsrc.Day, // TODO change to end/before
	s io.Store,
) (plays [][]charts.Song, err error) {
	if user.Registered == nil {
		return nil, fmt.Errorf("user '%v' has no valid registration date",
			user.Name)
	}
	registeredDay := user.Registered.Midnight()
	endCached := user.Registered

	cache := unpack.NewCached(s)

	bookmark, err := unpack.LoadBookmark(user.Name, s)
	if err == nil {
		endCached = bookmark
	}

	oldPlays, err := loadDays(user.Name, user.Registered, endCached, s)
	if err != nil {
		return nil, err
	}

	if len(oldPlays) > 0 {
		days := int((endCached.Midnight() - registeredDay) / 86400)
		oldPlays = oldPlays[:days]
	}

	if until == nil {
		return nil, errors.New("'until' is not a valid day")
	}

	if endCached.Midnight() > until.Midnight()+86400 {
		days := int((endCached.Midnight()-registeredDay)/86400) - 1
		return oldPlays[:days], nil
	}

	newPlays, err := loadHistory(
		unpack.User{Name: user.Name, Registered: endCached},
		until, io.FreshStore(s), cache) // TODO make fresh optional

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

	err := async.Pie(int(days), func(i int) error {
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
	io rsrc.IO,
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
			dp, err := loadDayPlays(user.Name, date, io, l)
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
	io rsrc.IO, cache unpack.Loader,
) ([]charts.Song, error) {
	firstPage, err := unpack.LoadHistoryDayPage(user, 1, time, unpack.NewCacheless(io))
	if err != nil {
		return nil, err
	}
	err = attachDuration(firstPage.Plays, cache)
	if err != nil {
		return nil, err
	}

	if firstPage.Pages < 2 {
		unpack.WriteDayHistory(firstPage.Plays, user, time, io)
		return firstPage.Plays, nil
	}

	pages := make([][]charts.Song, firstPage.Pages)
	pages[0] = firstPage.Plays

	back := make(chan loadDayPlaysResult)
	for page := 2; page <= len(pages); page++ {
		go func(page int) {
			histPage, tmpErr := unpack.LoadHistoryDayPage(user, page, time, unpack.NewCacheless(io))
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
		unpack.WriteDayHistory(plays, user, time, io)
	}
	return plays, err
}

func attachDuration(songs []charts.Song, cache unpack.Loader) error {
	async.Pie(len(songs), func(i int) error {
		if info, err := unpack.LoadTrackInfo(songs[i].Artist, songs[i].Title, cache); err == nil {
			songs[i].Duration = float64(info.Duration) / 60.0
			return nil
		} else {
			return err
		}
	})

	// not all tracks are found, ignore this
	return nil
}
