package organize

import (
	"fmt"
	"strings"

	async "github.com/nilsbu/async"
	"github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type tagResult struct {
	artist string
	tags   []charts.Tag
	err    error
}

// LoadArtistTags loads the tags for all given artists.
func LoadArtistTags(artists []string, r rsrc.Reader,
) (map[string][]charts.Tag, error) {
	tagLoader := unpack.NewCached(r)

	artistTags := make(map[string][]charts.Tag)
	feedback := make(chan *tagResult)
	for _, artist := range artists {
		go func(artist string) {
			tags, err := loadArtistTags(artist, r, tagLoader)
			feedback <- &tagResult{artist, tags, err}
		}(artist)
	}

	quit := make(chan bool)
	err := &async.MultiError{Msg: "could not load tags", Errs: []error{}}
	go func() {
		for range artists {
			res := <-feedback
			if res.err != nil {
				err.Errs = append(err.Errs, res.err)
			}
			artistTags[res.artist] = res.tags
		}
		quit <- true
	}()

	<-quit
	if len(err.Errs) > 0 {
		return artistTags, err
	}

	return artistTags, nil
}

func loadArtistTags(
	artist string,
	r rsrc.Reader,
	tl unpack.Loader,
) ([]charts.Tag, error) {

	tags, err := unpack.LoadArtistTags(artist, tl)
	if err != nil {
		switch err := err.(type) {
		case *unpack.LastfmError:
			err.Message = fmt.Sprintf(
				"could not load tags of '%v': %v",
				artist,
				err.Message)
			return nil, err
		default:
			return nil, errors.Wrapf(err, "could not load tags of '%v'", artist)
		}
	}

	wtags := make([]charts.Tag, len(tags))
	feedback := make(chan error)
	for i, tag := range tags {
		go func(i int, tag unpack.TagCount) {
			ti, terr := unpack.LoadTagInfo(tag.Name, tl)
			if terr == nil {
				wtags[i] = *ti
				wtags[i].Weight = tag.Count

				// all tags are lower case
				wtags[i].Name = strings.ToLower(wtags[i].Name)
			}
			feedback <- errors.Wrapf(terr, "could not load %v", tag.Name)
		}(i, tag)
	}

	err = nil
	for range tags {
		if terr := <-feedback; terr != nil {
			err = terr
		}
	}
	if err != nil {
		return nil, err
	}

	return wtags, nil
}
