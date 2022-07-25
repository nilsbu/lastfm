package organize

import (
	"errors"
	"fmt"

	async "github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// UpdateHistory loads saved daily plays from preprocessed all day plays and
// reads the remaining days from raw data. The last saved day gets reloaded.
func UpdateHistory(
	user *unpack.User,
	end rsrc.Day,
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

	oldPlays, err := LoadPreparedHistory(user.Name, user.Registered, endCached, s)
	if err != nil {
		return nil, err
	}

	if len(oldPlays) > 0 {
		days := rsrc.Between(user.Registered, endCached).Days()
		oldPlays = oldPlays[:days]
	}

	if end == nil {
		return nil, errors.New("'end' is not a valid day")
	}

	if rsrc.Between(end, endCached).Days() > 0 {
		days := rsrc.Between(user.Registered, endCached).Days() - 1
		return oldPlays[:days], nil
	}

	newPlays, err := LoadHistory(
		unpack.User{Name: user.Name, Registered: endCached},
		end, f, cache) // TODO make fresh optional
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

func LoadPreparedHistory(user string, begin, end rsrc.Day, r rsrc.Reader) ([][]info.Song, error) {
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
	end rsrc.Day,
	io rsrc.IO,
	l unpack.Loader) ([][]info.Song, error) {
	if end == nil {
		return nil, errors.New("parameter 'end' is no valid Day")
	} else if user.Registered == nil {
		return nil, errors.New("user has no valid registration date")
	} else {
		return loadHistory(user.Name, user.Registered, end, io, l)
	}
}

func loadHistory(
	user string,
	begin, end rsrc.Day,
	io rsrc.IO,
	l unpack.Loader) ([][]info.Song, error) {

	if days := rsrc.Between(begin, end).Days(); days <= 0 {
		return nil, nil
	} else {
		result := make([][]info.Song, days)
		errs := async.Pie(days, func(i int) error {
			date := begin.AddDate(0, 0, i)
			dp, err := loadDayPlays(user, date, io, l)
			result[i] = dp
			return err
		})
		return result, errs
	}
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

	errs := async.Pie(len(pages)-1, func(i int) error {
		page := i + 2
		if histPage, err := unpack.LoadHistoryDayPage(user, page, time, unpack.NewCacheless(io)); err != nil {
			return err
		} else {
			err := attachDuration(histPage.Plays, cache)
			pages[page-1] = histPage.Plays
			return err
		}
	})
	if errs != nil {
		return nil, errs
	}

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

// BackupUpdateHistory overwrites the prepared history by re-fetching the data.
// A number of days before the current bookmark, specified by delta, won't be re-fetched.
func BackupUpdateHistory(userName string, delta int, s io.Store) error {
	var bookmark, backup rsrc.Day
	if user, err := unpack.LoadUserInfo(userName, unpack.NewCacheless(s)); err != nil {
		return err
	} else if bookmark, err = unpack.LoadBookmark(userName, s); err != nil {
		return err
	} else if backup, err = unpack.LoadBackupBookmark(userName, s); err != nil {
		backup = user.Registered
	}
	end := bookmark.AddDate(0, 0, -delta)
	cache := unpack.NewCached(s)
	if songs, err := loadHistory(userName, backup, end, io.FreshStore(s), cache); err != nil {
		return err
	} else if len(songs) == 0 {
		return nil
	} else if err := unpack.WriteBackupBookmark(end, userName, s); err != nil {
		return err
	} else {
		return nil
	}
}
