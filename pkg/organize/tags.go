package organize

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type multiError struct {
	msg  string
	errs []error
}

func (err *multiError) Error() string {
	msg := err.msg + ":"

	for _, e := range err.errs {
		msg += "\n  " + e.Error()
	}

	return msg
}

type tagResult struct {
	artist string
	tags   []charts.Tag
	err    error
}

// LoadArtistTags loads the tags for all given artists.
func LoadArtistTags(artists []string, r rsrc.Reader,
) (map[string][]charts.Tag, error) {
	tagLoader := unpack.NewCachedTagLoader(r)

	artistTags := make(map[string][]charts.Tag)
	feedback := make(chan *tagResult)
	for _, artist := range artists {
		go func(artist string) {
			tags, err := loadArtistTags(artist, r, tagLoader)
			feedback <- &tagResult{artist, tags, err}
		}(artist)
	}

	quit := make(chan bool)
	err := &multiError{"could not load tags", []error{}}
	go func() {
		for range artists {
			if res := <-feedback; res.err != nil {
				err.errs = append(err.errs, res.err)
			} else {
				artistTags[res.artist] = res.tags
			}
		}
		quit <- true
	}()

	<-quit
	if len(err.errs) > 0 {
		return nil, err
	}

	return artistTags, nil
}

func loadArtistTags(
	artist string,
	r rsrc.Reader,
	tl *unpack.CachedTagLoader,
) ([]charts.Tag, error) {

	tags, err := unpack.LoadArtistTags(artist, r)
	if err != nil {
		return nil, errors.Wrapf(err, "could not load tags of '%v'", artist)
	}

	wtags := make([]charts.Tag, len(tags))
	feedback := make(chan error)
	for i, tag := range tags {
		go func(i int, tag unpack.TagCount) {
			ti, terr := tl.LoadTagInfo(tag.Name)
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
