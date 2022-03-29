package charts

import (
	"github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type load struct {
	r     rsrc.Reader
	user  string
	plays [][]info.Song
}

func LoadArtists(user string, r rsrc.Reader) Charts {
	l := load{
		r: r, user: user,
	}

	return new(l.songs,
		func(s info.Song) Title { return ArtistTitle(s.Artist) },
		func(s info.Song) float64 { return 1.0 })
}

func LoadArtistsDuration(user string, r rsrc.Reader) Charts {
	l := load{
		r: r, user: user,
	}

	return new(l.songs,
		func(s info.Song) Title { return ArtistTitle(s.Artist) },
		fDuration)
}

func LoadSongs(user string, r rsrc.Reader) Charts {
	l := load{
		r: r, user: user,
	}

	return new(l.songs,
		func(s info.Song) Title { return SongTitle(s) },
		func(s info.Song) float64 { return 1.0 })
}

func LoadSongsDuration(user string, r rsrc.Reader) Charts {
	l := load{
		r: r, user: user,
	}

	return new(l.songs,
		func(s info.Song) Title { return SongTitle(s) },
		fDuration)
}

func (w *load) load() error {
	var user *unpack.User
	var bookmark rsrc.Day
	var corrections map[string]string

	err := async.Pe([]func() error{
		func() error {
			var err error
			user, err = unpack.LoadUserInfo(w.user, unpack.NewCacheless(w.r))
			return errors.Wrap(err, "failed to load user info")
		},
		func() error {
			var err error
			corrections, err = unpack.LoadArtistCorrections(w.user, w.r)
			return err
		},
		func() error {
			var err error
			bookmark, err = unpack.LoadBookmark(w.user, w.r)
			return err
		},
	})
	if err != nil {
		return err
	}

	days := int((bookmark.Midnight() - user.Registered.Midnight()) / 86400)
	w.plays = make([][]info.Song, days+1)
	err = async.Pie(days+1, func(i int) error {
		day := user.Registered.AddDate(0, 0, i)
		if songs, err := unpack.LoadDayHistory(w.user, day, w.r); err == nil {
			for j, song := range songs {
				if c, ok := corrections[song.Artist]; ok {
					songs[j].Artist = c
				}
			}
			w.plays[i] = songs
			return nil
		} else {
			return err
		}
	})
	return err
}

func (w *load) songs() ([][]info.Song, error) {
	if w.plays == nil {
		err := w.load()
		return w.plays, err
	}
	return w.plays, nil
}
