package unpack

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

// LastfmError wraps an error returned  by Last.fm.
type LastfmError struct {
	Code    int
	Message string
}

func (err *LastfmError) Error() string {
	return fmt.Sprintf("LastFM error (code = %v): %v", err.Code, err.Message)
}

func (err *LastfmError) IsFatal() bool {
	return err.Code >= 8
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

type obHistorySingle struct {
	obHistory
}

// LoadHistoryDayPage loads a page of a user's played tracks.
func LoadHistoryDayPage(
	user string, page int, day rsrc.Day, r rsrc.Reader) (*HistoryDayPage, error) {
	data, err := obtain(&obHistory{user, page, day}, r)
	if err != nil {
		data, err = obtain(&obHistorySingle{obHistory{user, page, day}}, r)
		if err != nil {
			return nil, err
		}
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

func (o *obHistorySingle) deserializer() interface{} {
	return &jsonUserRecentTrackSingle{}
}

func (o *obHistorySingle) interpret(raw interface{}) (interface{}, error) {
	data := raw.(*jsonUserRecentTrackSingle)

	urt := &jsonUserRecentTracks{
		RecentTracks: jsonRecentTracks{
			Track: []jsonTrack{data.RecentTracks.Track},
			Attr:  data.RecentTracks.Attr,
		}}

	return &HistoryDayPage{
		countPlays(urt),
		data.RecentTracks.Attr.TotalPages}, nil
}

// ArtistInfo contains information about an artist
type ArtistInfo struct {
	Name      string
	Listeners int64
	PlayCount int64
}

type obArtistInfo struct {
	name string
}

// LoadArtistInfo reads information of an artist.
func LoadArtistInfo(artist string, r rsrc.Reader) (*ArtistInfo, error) {
	data, err := obtain(&obArtistInfo{artist}, r)
	if err != nil {
		return nil, err
	}
	info := data.(*ArtistInfo)
	return info, nil
}

func (o *obArtistInfo) locator() rsrc.Locator {
	return rsrc.ArtistInfo(o.name)
}

func (o *obArtistInfo) deserializer() interface{} {
	return &jsonArtistInfo{}
}

func (o *obArtistInfo) interpret(raw interface{}) (interface{}, error) {
	ai := raw.(*jsonArtistInfo)
	return &ArtistInfo{
		Name:      ai.Artist.Name,
		Listeners: ai.Artist.Stats.Listeners,
		PlayCount: ai.Artist.Stats.PlayCount,
	}, nil
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

type obTagInfo struct {
	name string
}

// LoadTagInfo loads tag information.
func LoadTagInfo(tag string, buf *CachedLoader) (*charts.Tag, error) {
	data, err := buf.Load(&obTagInfo{tag})
	if err != nil {
		return nil, err
	}
	return data.(*charts.Tag), nil
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

// TrackInfo contains general information about a track.
type TrackInfo struct {
	Artist    string
	Track     string
	Duration  int
	Listeners int64
	Playcount int64
}

type obTrackInfo struct {
	artist string
	track  string
}

// LoadTrackInfo reads the track information.
func LoadTrackInfo(artist, track string, buf *CachedLoader) (TrackInfo, error) {
	data, err := buf.Load(&obTrackInfo{artist, track})
	if err != nil {
		return TrackInfo{}, err
	}
	info := data.(TrackInfo)
	info.Artist = artist
	info.Track = track
	return info, nil
}

// WriteTrackInfo writes the track information.
func WriteTrackInfo(artist, track string, tags TrackInfo, w rsrc.Writer) error {
	return deposit(tags, &obTrackInfo{artist: artist, track: track}, w)
}

func (o *obTrackInfo) locator() rsrc.Locator {
	return rsrc.TrackInfo(o.artist, o.track)
}

func (o *obTrackInfo) deserializer() interface{} {
	return &jsonTrackInfo{}
}

func (o *obTrackInfo) interpret(raw interface{}) (interface{}, error) {
	ti := raw.(*jsonTrackInfo).Track

	info := TrackInfo{
		Duration:  ti.Duration / 1000,
		Listeners: ti.Listeners,
		Playcount: ti.Playcount,
	}
	return info, nil
}

func (o *obTrackInfo) raw(obj interface{}) interface{} {
	info := obj.(TrackInfo)
	js := jsonTrackInfo{Track: jsonTrackTrack{
		Duration:  info.Duration * 1000,
		Listeners: info.Listeners,
		Playcount: info.Playcount,
	}}
	return js
}
