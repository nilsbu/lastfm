package unpack

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/pkg/errors"
)

type LastfmError struct {
	Code    int
	Message string
}

// obError is a deserializer. It parses a resource as a jsonError.
type obError struct {
}

func (o *obError) deserializer() interface{} {
	return &jsonError{}
}

func (o *obError) interpret(raw interface{}) (interface{}, error) {
	e := raw.(*jsonError)

	return &LastfmError{e.Error, e.Message}, nil
}

// User contains relevant core information about a user.
type User struct {
	Name       string
	Registered rsrc.Day
}

type obUserInfo struct {
	name string
}

// LoadUserInfo loads a user's registration date. It is returned along with the
// name.
func LoadUserInfo(name string, r rsrc.Reader) (*User, error) {
	data, err := obtain(&obUserInfo{name}, r)
	if err != nil {
		return nil, err
	}
	user := data.(*User)
	return user, nil
}

// WriteUserInfo writes a user's registration date. The playcount is set to 0.
func WriteUserInfo(user *User, w rsrc.Writer) error {
	return deposit(user, &obUserInfo{user.Name}, w)
}

func (o *obUserInfo) locator() rsrc.Locator {
	return rsrc.UserInfo(o.name)
}

func (o *obUserInfo) deserializer() interface{} {
	return &jsonUserInfo{}
}

func (o *obUserInfo) interpret(raw interface{}) (interface{}, error) {
	ui := raw.(*jsonUserInfo)

	utc := ui.User.Registered.UTC
	return &User{ui.User.Name, rsrc.ToDay(utc)}, nil
}

func (o *obUserInfo) raw(obj interface{}) interface{} {
	user := obj.(*User)
	js := jsonUserInfo{User: jsonUser{
		Name:       user.Name,
		PlayCount:  0,
		Registered: jsonTime{UTC: user.Registered.Midnight()},
	}}
	return js
}

// HistoryDayPage is a single page of a day of a user's played tracks.
type HistoryDayPage struct {
	Plays []charts.Song
	Pages int
}

type obHistory struct {
	user string
	page int
	day  rsrc.Day
}

// LoadHistoryDayPage loads a page of a user's played tracks.
func LoadHistoryDayPage(
	user string, page int, day rsrc.Day, r rsrc.Reader) (*HistoryDayPage, error) {
	data, err := obtain(&obHistory{user, page, day}, r)
	if err != nil {
		return nil, err
	}
	hist := data.(*HistoryDayPage)
	return hist, nil
}

func (o *obHistory) locator() rsrc.Locator {
	return rsrc.History(o.user, o.page, o.day)
}

func (o *obHistory) deserializer() interface{} {
	return &jsonUserRecentTracks{}
}

func (o *obHistory) interpret(raw interface{}) (interface{}, error) {
	data := raw.(*jsonUserRecentTracks)

	return &HistoryDayPage{
		countPlays(data),
		data.RecentTracks.Attr.TotalPages}, nil
}

func countPlays(urt *jsonUserRecentTracks) []charts.Song {
	plays := []charts.Song{}
	for _, track := range urt.RecentTracks.Track {
		if !track.Attr.NowPlaying {

			plays = append(plays, charts.Song{
				Artist: track.Artist.Str,
				Title:  track.Name,
				Album:  track.Album.Str,
			})
		}
	}
	return plays
}

// TagCount assigns a tag a value.
type TagCount struct {
	Name  string
	Count int
}

type obArtistTags struct {
	name string
}

// LoadArtistTags reads the top tags of an artist.
func LoadArtistTags(artist string, r rsrc.Reader) ([]TagCount, error) {
	data, err := obtain(&obArtistTags{artist}, r)
	if err != nil {
		return nil, err
	}
	tags := data.([]TagCount)
	return tags, nil
}

// WriteArtistTags writes the top tags of an artist.
func WriteArtistTags(artist string, tags []TagCount, w rsrc.Writer) error {
	return deposit(tags, &obArtistTags{name: artist}, w)
}

func (o *obArtistTags) locator() rsrc.Locator {
	return rsrc.ArtistTags(o.name)
}

func (o *obArtistTags) deserializer() interface{} {
	return &jsonArtistTags{}
}

func (o *obArtistTags) interpret(raw interface{}) (interface{}, error) {
	at := raw.(*jsonArtistTags)

	len := len(at.TopTags.Tags)

	tags := make([]TagCount, len)
	for i, tag := range at.TopTags.Tags {
		tags[i] = TagCount{Name: tag.Name, Count: tag.Count}
	}
	return tags, nil
}

func (o *obArtistTags) raw(obj interface{}) interface{} {
	tags := obj.([]TagCount)
	jsTags := []jsonTag{}
	for _, tag := range tags {
		jsTags = append(jsTags, jsonTag{Name: tag.Name, Count: tag.Count})
	}

	js := jsonArtistTags{TopTags: jsonTopTags{
		Tags: jsTags,
		Attr: jsonTopTagAttr{Artist: ""}, // Artist name isn't available here
	}}
	return js
}

// cachedLoader if a buffer that stores data that can be obtained. It minimizes
// the number of external calls.
// TODO extract as file / module.
type cachedLoader struct {
	reader      rsrc.Reader
	requestChan chan cacheRequest
}

type cacheRequest struct {
	name     string
	obtainer obtainer
	back     chan cacheResult
}

type cacheResult struct {
	name string
	data interface{}
	err  error
}

type obTagInfo struct {
	name string
}

// newCachedLoader creates a buffer which can read and store data.
func newCachedLoader(r rsrc.Reader) *cachedLoader {
	buf := &cachedLoader{
		reader:      r,
		requestChan: make(chan cacheRequest),
	}

	go buf.worker()
	return buf
}

func (buf *cachedLoader) worker() {
	resultChan := make(chan cacheResult)

	requests := make(map[string][]cacheRequest)
	cacheMap := make(map[string]cacheResult)

	hasError := false

	for {
		select {
		case request := <-buf.requestChan:
			if hasError {
				request.back <- cacheResult{
					name: request.name,
					data: nil,
					err:  errors.New("abort due to previous error"),
				}
				close(request.back)
			} else if tc, ok := cacheMap[request.name]; ok {
				request.back <- tc
				close(request.back)
			} else if tc, ok := requests[request.name]; ok && tc != nil {
				requests[request.name] = append(requests[request.name], request)
			} else {
				requests[request.name] = append(requests[request.name], request)
				go func(request cacheRequest) {
					data, err := obtain(request.obtainer, buf.reader)
					if err != nil {
						hasError = true
						resultChan <- cacheResult{request.name, nil, err}
					} else {
						resultChan <- cacheResult{request.name, data, nil}
					}
				}(request)
			}

		case result := <-resultChan:
			cacheMap[result.name] = result

			for _, request := range requests[result.name] {
				request.back <- result
				close(request.back)
			}

			requests[result.name] = nil
		}
	}
}

// CachedTagLoader is a buffer that loads tag information with minimal amount of
// external calls.
type CachedTagLoader interface {
	LoadTagInfo(tag string) (*charts.Tag, error)
}

type cachedTagLoader struct {
	Loader *cachedLoader
}

// NewCachedTagLoader creates a buffer which can read and store tag information.
func NewCachedTagLoader(r rsrc.Reader) CachedTagLoader {
	return cachedTagLoader{newCachedLoader(r)}
}

// LoadTagInfo loads tag information.
func (buf cachedTagLoader) LoadTagInfo(tag string) (*charts.Tag, error) {
	back := make(chan cacheResult)

	buf.Loader.requestChan <- cacheRequest{
		name:     tag,
		obtainer: &obTagInfo{tag},
		back:     back,
	}

	result := <-back
	if result.err != nil {
		return nil, result.err
	}

	return result.data.(*charts.Tag), nil
}

// WriteTagInfo writes tag infos.
func WriteTagInfo(tag *charts.Tag, w rsrc.Writer) error {
	return deposit(tag, &obTagInfo{name: tag.Name}, w)
}

func (o *obTagInfo) locator() rsrc.Locator {
	return rsrc.TagInfo(o.name)
}

func (o *obTagInfo) deserializer() interface{} {
	return &jsonTagInfo{}
}

func (o *obTagInfo) interpret(raw interface{}) (interface{}, error) {
	tag := raw.(*jsonTagInfo)

	return &charts.Tag{
		Name:  tag.Tag.Name,
		Total: tag.Tag.Total,
		Reach: tag.Tag.Reach}, nil
}

func (o *obTagInfo) raw(obj interface{}) interface{} {
	tag := obj.(*charts.Tag)
	js := jsonTagInfo{Tag: jsonTagTag{
		Name:  tag.Name,
		Total: tag.Total,
		Reach: tag.Reach,
	}}
	return js
}

type obArtistSimilar struct {
	name string
}

func (o *obArtistSimilar) locator() rsrc.Locator {
	return rsrc.ArtistSimilar(o.name)
}

func (o *obArtistSimilar) deserializer() interface{} {
	return &jsonArtistSimilar{}
}

func (o *obArtistSimilar) interpret(raw interface{}) (interface{}, error) {
	inArtists := raw.(*jsonArtistSimilar).SimilarArtists

	outArtists := []SimilarArtist{}

	for _, inArtist := range inArtists.Matches {
		outArtists = append(outArtists,
			SimilarArtist{Name: inArtist.Name, Match: inArtist.Match})
	}

	return outArtists, nil
}

func (o *obArtistSimilar) raw(obj interface{}) interface{} {
	inSimilar := obj.([]SimilarArtist)

	jsSimilar := []jsonArtistSimilarMatch{}
	for _, sim := range inSimilar {
		jsSimilar = append(jsSimilar, jsonArtistSimilarMatch{
			Name:  sim.Name,
			Match: sim.Match,
		})
	}

	js := jsonArtistSimilar{SimilarArtists: jsonArtistMatches{Matches: jsSimilar}}
	return js
}

// SimilarArtist contains the matching score of a similar artist with respect
// to a requested artist.
type SimilarArtist struct {
	Name  string
	Match float32
}

// CachedSimilarLoader is a buffer that loads artists' similar artists with
// minimal amount of external calls.
type CachedSimilarLoader interface {
	LoadArtistSimilar(artist string) ([]SimilarArtist, error)
}

type cachedSimilarLoader struct {
	Loader *cachedLoader
}

// NewCachedSimilarLoader creates a buffer which can read and store similar
// artists.
func NewCachedSimilarLoader(r rsrc.Reader) CachedSimilarLoader {
	return cachedSimilarLoader{newCachedLoader(r)}
}

// LoadArtistSimilar loads an artist's similar artists.
func (buf cachedSimilarLoader) LoadArtistSimilar(
	artist string) ([]SimilarArtist, error) {
	back := make(chan cacheResult)

	buf.Loader.requestChan <- cacheRequest{
		name:     artist,
		obtainer: &obArtistSimilar{artist},
		back:     back,
	}

	result := <-back
	if result.err != nil {
		return nil, result.err
	}

	return result.data.([]SimilarArtist), nil
}

// WriteArtistSimilar writes the similar artists of an artist.
func WriteArtistSimilar(
	artist string,
	similar []SimilarArtist,
	w rsrc.Writer,
) error {
	return deposit(similar, &obArtistSimilar{name: artist}, w)
}
