package io

// UserRecentTracks is unmarshaled JSON data from user.getRecentTracks.
type UserRecentTracks struct {
	Recenttracks struct {
		Track []Track `json:"track"`
		Attr  struct {
			User       string `json:"user"`
			Page       int    `json:"page,string"`
			PerPage    int    `json:"perPage,string"`
			TotalPages int    `json:"totalPages,string"`
			Total      int    `json:"total,string"`
		} `json:"@attr"`
	} `json:"recenttracks"`
}

// Track is unmarshaled JSON data from user.getRecentTracks' track tag.
type Track struct {
	Artist     Artist  `json:"artist"`
	Name       string  `json:"name"`
	Streamable string  `json:"streamable"`
	Mbid       string  `json:"mbid"`
	Album      Album   `json:"album"`
	URL        string  `json:"url"`
	Image      []Image `json:"image"`
	Date       Date    `json:"date"`
}

// Artist is unmarshaled JSON data from user.getRecentTracks' artist tag.
type Artist struct {
	MBID string `json:"mbid"`
	Text string `json:"#text"`
}

// Album is unmarshaled JSON data from user.getRecentTracks' album tag.
type Album struct {
	Artist
}
