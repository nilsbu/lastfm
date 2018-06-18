package unpack

// UserInfo is unmarshaled JSON data from user.getInfo.
type UserInfo struct {
	User userUser `json:"user"`
}

type userUser struct {
	Name       string `json:"name"`
	PlayCount  int    `json:"playcount"`
	Registered time   `json:"registered"`
	// Not Included: realname, image, url, country, age, gender, subscriber, type
	//               playlists, bootstrap
}

// UserRecentTracks is unmarshaled JSON data from user.getRecentTracks.
type UserRecentTracks struct {
	RecentTracks recentTracks `json:"recenttracks"`
}

type recentTracks struct {
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
	Artist text   `json:"artist"`
	Name   string `json:"name"`
	Album  text   `json:"album"`
	Date   date   `json:"date"`
	// Not included: streamable, mbid, url, image
}

type text struct {
	Str string `json:"#text"`
	// Not included: mbid
}

// APIKey is an unmarshaled JSON API key.
type APIKey struct {
	Key string `json:"apikey"`
}

// Bookmark is an unmarshaled JSON bookmark for daily plays.
type Bookmark struct {
	UTC        int64  `json:"unixtime"`
	TimeString string `json:"time"`
}
