package organize

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type similarResult struct {
	artist  string
	similar map[string]float32
	err     error
}

// LoadArtistSimilar loads the similar artists for all given artists.
func LoadArtistSimilar(artists []charts.Key, r rsrc.Reader,
) (map[string]map[string]float32, error) {
	similarLoader := unpack.NewCachedSimilarLoader(r)

	artistSimilars := make(map[string]map[string]float32)
	feedback := make(chan *similarResult)
	for _, artist := range artists {
		go func(artist string) {
			sim, err := similarLoader.LoadArtistSimilar(artist)
			feedback <- &similarResult{artist, sim, err}
		}(artist.ArtistName())
	}

	err := &MultiError{"could not load tags", []error{}}
	for range artists {
		res := <-feedback
		if res.err != nil {
			err.Errs = append(err.Errs, res.err)
		}
		artistSimilars[res.artist] = res.similar
	}

	if len(err.Errs) > 0 {
		return artistSimilars, err
	}

	return artistSimilars, nil
}
