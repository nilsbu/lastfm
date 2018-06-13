package unpack

// UserRecentTracks is unmarshaled JSON data from user.getRecentTracks.
type UserRecentTracks struct {
	RecentTracks recentTracks `json:"recenttracks"`
}

type recentTracks struct {
	Track []Track          `json:"track"`
	Attr  recentTracksAttr `json:"@attr"`
}

type recentTracksAttr struct {
	User       string `json:"user"`
	Page       int    `json:"page,string"`
	PerPage    int    `json:"perPage,string"`
	TotalPages int    `json:"totalPages,string"`
	Total      int    `json:"total,string"`
}

// Date is an unmarshaled JSON date tag, that contains a unix timestamp.
type Date struct {
	UTS int `json:"uts,string"`
	// Not included: #text
}

// Track is unmarshaled JSON data from user.getRecentTracks' track tag.
type Track struct {
	Artist Text   `json:"artist"`
	Name   string `json:"name"`
	Album  Text   `json:"album"`
	Date   Date   `json:"date"`
	// Not included: streamable, mbid, url, image
}

// Text is an unmarshaled JSON text tag with an omitted MBID attribute.
type Text struct {
	Str string `json:"#text"`
	// Not included: mbid
}

// APIKey is an unmarshaled JSON tag for an API key.
type APIKey struct {
	Key string `json:"apikey"`
}
