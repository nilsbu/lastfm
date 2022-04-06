package charts

import (
	"strings"

	"github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type userLoad struct {
	user string
	r    rsrc.Reader
}

type load struct {
	userLoad
	plays [][]info.Song
}

func LoadArtists(user string, r rsrc.Reader) Charts {
	l := load{userLoad: userLoad{user: user, r: r}}

	return new(l.songs,
		func(s info.Song) Title { return ArtistTitle(s.Artist) },
		func(s info.Song) float64 { return 1.0 })
}

func LoadArtistsDuration(user string, r rsrc.Reader) Charts {
	l := load{userLoad: userLoad{user: user, r: r}}

	return new(l.songs,
		func(s info.Song) Title { return ArtistTitle(s.Artist) },
		fDuration)
}

func LoadSongs(user string, r rsrc.Reader) Charts {
	l := load{userLoad: userLoad{user: user, r: r}}

	return new(l.songs,
		func(s info.Song) Title { return SongTitle(s) },
		func(s info.Song) float64 { return 1.0 })
}

func LoadSongsDuration(user string, r rsrc.Reader) Charts {
	l := load{userLoad: userLoad{user: user, r: r}}

	return new(l.songs,
		func(s info.Song) Title { return SongTitle(s) },
		fDuration)
}

func (w *load) load() error {
	var user *unpack.User
	var bookmark rsrc.Day
	var corrections map[string]string

	// TODO there's duplication with userLoad.span()
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

	days := rsrc.Between(user.Registered, bookmark).Days()
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

type tagPartition struct {
	titles      func() []Title
	decide      map[string]string
	corrections map[string]string
	r           rsrc.Reader
}

func TagPartition(parent Charts, decide, corrections map[string]string, r rsrc.Reader) Partition {
	return &tagPartition{
		titles:      parent.Titles,
		decide:      decide,
		corrections: corrections,
		r:           r,
	}
}

func (p *tagPartition) Titles(partition Title) ([]Title, error) {
	parentTitles := p.titles()
	partitions := make([]string, len(parentTitles))
	loader := unpack.NewCacheless(p.r)

	// TODO what to do with errors in tags?
	// An error occurs for the artist "Kamiyada+".
	async.Pie(len(parentTitles), func(i int) error {
		if partition, ok := p.corrections[parentTitles[i].Artist()]; ok {
			partitions[i] = partition
			return nil
		}

		tags, err := unpack.LoadArtistTags(parentTitles[i].Artist(), loader)
		if err != nil {
			return err
		} else {
			for _, tag := range tags {
				if partition, ok := p.decide[strings.ToLower(tag.Name)]; ok {
					partitions[i] = partition
					return nil
				}
			}
			partitions[i] = "-"
			return nil
		}
	})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }

	titles := make([]Title, len(parentTitles))
	j := 0
	for i, title := range parentTitles {
		if partitions[i] == partition.String() {
			titles[j] = title
			j++
		}
	}

	return titles[:j], nil
}

func (p *tagPartition) Partitions() ([]Title, error) {
	partitions := make([]Title, 0)
	for _, partition := range p.decide {
		found := false
		for _, p := range partitions {
			if partition == p.Key() {
				found = true
				break
			}
		}
		if !found {
			partitions = append(partitions, KeyTitle(partition))
		}
	}
	partitions = append(partitions, keyTitle("-"))
	return partitions, nil
}
