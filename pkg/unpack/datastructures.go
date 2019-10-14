package unpack

type jsonError struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

type jsonUserInfo struct {
	User jsonUser `json:"user"`
}

type jsonUser struct {
	Name       string   `json:"name"`
	PlayCount  int      `json:"playcount"`
	Registered jsonTime `json:"registered"`
	// Not Included: realname, image, url, country, age, gender, subscriber, type
	//               playlists, bootstrap
}

type jsonUserRecentTracks struct {
	RecentTracks jsonRecentTracks `json:"recenttracks"`
}

type jsonRecentTracks struct {
	Track []jsonTrack          `json:"track"`
	Attr  jsonRecentTracksAttr `json:"@attr"`
}

type jsonUserRecentTrackSingle struct {
	RecentTracks jsonRecentTrackSingle `json:"recenttracks"`
}

type jsonRecentTrackSingle struct {
	Track jsonTrack            `json:"track"`
	Attr  jsonRecentTracksAttr `json:"@attr"`
}

type jsonRecentTracksAttr struct {
	User       string `json:"user"`
	Page       int    `json:"page,string"`
	PerPage    int    `json:"perPage,string"`
	TotalPages int    `json:"totalPages,string"`
	Total      int    `json:"total,string"`
}

type jsonDate struct {
	UTC int64 `json:"uts,string"`
	// Not included: #text
}

type jsonTime struct {
	UTC int64 `json:"unixtime"`
	// Not included: #text
}

type jsonTrack struct {
	Artist jsonText      `json:"artist"`
	Name   string        `json:"name"`
	Album  jsonText      `json:"album"`
	Date   jsonDate      `json:"date"`
	Attr   jsonTrackAttr `json:"@attr"`
	// Not included: streamable, mbid, url, image
}

type jsonTrackAttr struct {
	NowPlaying bool `json:"nowplaying,string"`
}

type jsonText struct {
	Str string `json:"#text"`
	// Not included: mbid
}

type jsonArtistInfo struct {
	Artist jsonArtistArtist `json:"artist"`
}

type jsonArtistArtist struct {
	Name  string          `json:"name"`
	Stats jsonArtistStats `json:"stats"`
}

type jsonArtistStats struct {
	Listeners int64 `json:"listeners,string"`
	PlayCount int64 `json:"playcount,string"`
}

type jsonArtistTags struct {
	TopTags jsonTopTags `json:"toptags"`
}

type jsonTopTags struct {
	Tags []jsonTag      `json:"tag"`
	Attr jsonTopTagAttr `json:"@attr"`
}

type jsonTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type jsonTopTagAttr struct {
	Artist string `json:"artist"`
}

type jsonTagInfo struct {
	Tag jsonTagTag `json:"tag"`
}

type jsonTagTag struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
	Reach int64  `json:"reach"`
	// Not included: wiki
}

type jsonAPIKey struct {
	Key string `json:"apikey"`
}

type jsonSessionInfo struct {
	User string `json:"user"`
}

type jsonCorrections struct {
	Corrections map[string]string `json:"corrections"`
}

type jsonBookmark struct {
	NextDay string `json:"nextday"`
}
