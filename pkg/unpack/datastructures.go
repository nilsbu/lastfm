package unpack

type jsonUserInfo struct {
	User jsonUser `json:"user"`
}

type jsonUser struct {
	Name       string `json:"name"`
	PlayCount  int    `json:"playcount"`
	Registered time   `json:"registered"`
	// Not Included: realname, image, url, country, age, gender, subscriber, type
	//               playlists, bootstrap
}

type jsonUserRecentTracks struct {
	RecentTracks jsonRecentTracks `json:"recenttracks"`
}

type jsonRecentTracks struct {
	Track []track          `json:"track"`
	Attr  recentTracksAttr `json:"@attr"`
}

type recentTracksAttr struct {
	User       string `json:"user"`
	Page       int    `json:"page,string"`
	PerPage    int    `json:"perPage,string"`
	TotalPages int    `json:"totalPages,string"`
	Total      int    `json:"total,string"`
}

type date struct {
	UTC int64 `json:"uts,string"`
	// Not included: #text
}

type time struct {
	UTC int64 `json:"unixtime"`
	// Not included: #text
}

type track struct {
	Artist text      `json:"artist"`
	Name   string    `json:"name"`
	Album  text      `json:"album"`
	Date   date      `json:"date"`
	Attr   trackAttr `json:"@attr"`
	// Not included: streamable, mbid, url, image
}

type trackAttr struct {
	NowPlaying bool `json:"nowplaying,string"`
}

type text struct {
	Str string `json:"#text"`
	// Not included: mbid
}

type ArtistInfo struct {
	Artist artistArtist `json:"artist"`
}

type artistArtist struct {
	Name  string      `json:"name"`
	Stats artistStats `json:"stats"`
}

type artistStats struct {
	Listeners int64 `json:"listeners,string"`
	PlayCount int64 `json:"playcount,string"`
}

type ArtistTags struct {
	TopTags topTags `json:"toptags"`
}

type topTags struct {
	Tags []tag `json:"tag"`
}

type tag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type jsonAPIKey struct {
	Key string `json:"apikey"`
}

// SessionID is an unmarshaled JSON session identifier.
type SessionID struct {
	User string `json:"user"`
}
