package io

// Midnight is a unix time at midnight. It defaults to -1.
type Midnight int

// Page is a page of a multi-page resource. It defaults to 0.
type Page int

// Name is the name of a user, artist or tag
type Name string

// Resource is a general descriptor for local files or Last.fm URLs.
type Resource struct {
	// TODO replace string with custom types
	main   string
	method string
	name   Name
	page   Page
	time   Midnight
}

// NewUserInfo returns the Resource for "user.getInfo".
func NewUserInfo(name Name) *Resource {
	rsrc := new(Resource)
	rsrc.main = "user"
	rsrc.method = "getInfo"
	rsrc.name = name
	rsrc.time = -1
	return rsrc
}

// NewUserRecentTracks returns the Resource for "user.getRecentTracks".
func NewUserRecentTracks(name Name, page Page, time Midnight) *Resource {
	rsrc := new(Resource)
	rsrc.main = "user"
	rsrc.method = "getRecentTracks"
	rsrc.name = name
	rsrc.page = page
	rsrc.time = time - time%86400
	return rsrc
}

// NewArtistInfo returns the Resource for "artist.getInfo".
func NewArtistInfo(name Name) *Resource {
	rsrc := new(Resource)
	rsrc.main = "artist"
	rsrc.method = "getInfo"
	rsrc.name = name
	rsrc.time = -1
	return rsrc
}
