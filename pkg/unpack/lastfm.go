package unpack

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

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
	utc, _ := user.Registered.Midnight()
	js := jsonUserInfo{User: jsonUser{
		Name:       user.Name,
		PlayCount:  0,
		Registered: jsonTime{UTC: utc},
	}}
	return js
}

// HistoryDayPage is a single page of a day of a user's played tracks.
type HistoryDayPage struct {
	Plays charts.Charts
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

func countPlays(urt *jsonUserRecentTracks) charts.Charts {
	plays := charts.Charts{}
	for _, track := range urt.RecentTracks.Track {
		if !track.Attr.NowPlaying {
			if cnt, ok := plays[track.Artist.Str]; ok {
				cnt[0]++
			} else {
				plays[track.Artist.Str] = []float64{1}
			}
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

// CachedTagLoader if a buffer that stores tag information.
type CachedTagLoader struct {
	reader      rsrc.Reader
	requestChan chan tagRequest
}

type tagRequest struct {
	name string
	back chan tagResult
}

type tagResult struct {
	name string
	tag  *charts.Tag
	err  error
}

type obTagInfo struct {
	name string
}

// NewCachedTagLoader creates a buffer which can read and store tag information.
func NewCachedTagLoader(r rsrc.Reader) *CachedTagLoader {
	buf := &CachedTagLoader{
		reader:      r,
		requestChan: make(chan tagRequest),
	}

	go buf.worker()
	return buf
}

func (buf *CachedTagLoader) worker() {
	resultChan := make(chan tagResult)

	requests := make(map[string][]tagRequest)
	tagCounts := make(map[string]tagResult)

	for {
		select {
		case request := <-buf.requestChan:
			if tc, ok := tagCounts[request.name]; ok {
				request.back <- tc
				close(request.back)
			} else {

				requests[request.name] = append(requests[request.name], request)
				go func(request tagRequest) {
					data, err := obtain(&obTagInfo{request.name}, buf.reader)
					if err != nil {
						resultChan <- tagResult{request.name, nil, err}
					} else {
						tag := data.(*charts.Tag)
						resultChan <- tagResult{request.name, tag, nil}
					}
				}(request)
			}

		case tag := <-resultChan:
			tagCounts[tag.name] = tag

			for _, request := range requests[tag.name] {
				request.back <- tag
				close(request.back)
			}

			requests[tag.name] = nil
		}
	}
}

// LoadTagInfo loads tag information.
func (buf *CachedTagLoader) LoadTagInfo(artist string) (*charts.Tag, error) {
	back := make(chan tagResult)

	buf.requestChan <- tagRequest{
		name: artist,
		back: back,
	}

	result := <-back
	if result.err != nil {
		return nil, result.err
	}

	return result.tag, nil
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
