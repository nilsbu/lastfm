package organize

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type trackInfoResult struct {
	idx  int
	info unpack.TrackInfo
	err  error
}

// LoadTrackInfos lods track infos for all songs.
func LoadTrackInfos(
	songs []charts.Song, r rsrc.Reader,
) ([]unpack.TrackInfo, error) {
	trackInfos := make([]unpack.TrackInfo, len(songs))

	cache := unpack.NewCached(r)

	feedback := make(chan *trackInfoResult)
	for i, song := range songs {
		go func(i int, song charts.Song) {
			infos, err := unpack.LoadTrackInfo(song.Artist, song.Title, cache)
			feedback <- &trackInfoResult{i, infos, err}
		}(i, song)
	}

	quit := make(chan bool)
	err := &MultiError{"could not load track infos", []error{}}
	go func() {
		for range songs {
			res := <-feedback
			if res.err != nil {
				err.Errs = append(err.Errs, res.err)
			} else {
				trackInfos[res.idx] = res.info
			}
		}
		quit <- true
	}()

	<-quit
	if len(err.Errs) > 0 {
		return trackInfos, err
	}

	return trackInfos, nil
}
