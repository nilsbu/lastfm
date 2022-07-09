package organize

import (
	"errors"
	"fmt"

	async "github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// UpdateHistory loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateHistory(
	user *unpack.User,
	until rsrc.Day, // TODO change to end/before
	s, f rsrc.IO,
) (plays [][]info.Song, err error) {
	if user.Registered == nil {
		return nil, fmt.Errorf("user '%v' has no valid registration date",
			user.Name)
	}
	cache := unpack.NewCached(s)

	endCached := user.Registered
	if bookmark, err := unpack.LoadBookmark(user.Name, s); err == nil {
		endCached = bookmark
	}

	oldPlays, err := loadDays(user.Name, user.Registered, endCached, s)
	if err != nil {
		return nil, err
	}

	if len(oldPlays) > 0 {
		days := rsrc.Between(user.Registered, endCached).Days()
		oldPlays = oldPlays[:days]
	}

	if until == nil {
		return nil, errors.New("'until' is not a valid day")
	}

	if rsrc.Between(until.AddDate(0, 0, 1), endCached).Days() > 0 {
		days := rsrc.Between(user.Registered, endCached).Days() - 1
		return oldPlays[:days], nil
	}

	newPlays, err := LoadHistory(
		unpack.User{Name: user.Name, Registered: endCached},
		until, f, cache) // TODO make fresh optional
	if err != nil {
		return nil, err
	}

	for _, plays := range newPlays {
		err = attachDuration(plays, cache)
		if err != nil {
			return nil, err
		}
	}

	return append(oldPlays, newPlays...), err
}

func loadDays(user string, begin, end rsrc.Day, r rsrc.Reader) ([][]info.Song, error) {
	days := rsrc.Between(begin, end).Days()
	plays := make([][]info.Song, days)

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

// LoadHistory load plays from all days since the registration of the user.
func LoadHistory(
	user unpack.User,
	until rsrc.Day,
	io rsrc.IO,
	l unpack.Loader) ([][]info.Song, error) {
	if until == nil {
		return nil, errors.New("parameter 'until' is no valid Day")
	} else if user.Registered == nil {
		return nil, errors.New("user has no valid registration date")
	} else {
		return loadHistory(user.Name, user.Registered, until, io, l)
	}
}

func loadHistory(
	user string,
	begin, until rsrc.Day,
	io rsrc.IO,
	l unpack.Loader) ([][]info.Song, error) {

	days := rsrc.Between(begin, until).Days()
	result := make([][]info.Song, days+1)
	errs := async.Pie(len(result), func(i int) error {
		date := begin.AddDate(0, 0, i)
		dp, err := loadDayPlays(user, date, io, l)
		result[i] = dp
		return err
	})
	return result, errs
}

// loadDayPlaysResult is the result of loadDayPlays.
type loadDayPlaysResult struct {
	DayPlays []info.Song
	Page     int
	Err      error
}

func loadDayPlays(
	user string,
	time rsrc.Day,
	io rsrc.IO, cache unpack.Loader,
) ([]info.Song, error) {
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

	pages := make([][]info.Song, firstPage.Pages)
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

	plays := []info.Song{}
	for _, page := range pages {
		plays = append(plays, page...)
		unpack.WriteDayHistory(plays, user, time, io)
	}
	return plays, err
}

func attachDuration(songs []info.Song, cache unpack.Loader) error {
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
