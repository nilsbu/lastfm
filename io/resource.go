package io

import "strconv"

// Midnight is a unix time at midnight.
type Midnight int

// Page is a page of a multi-page resource. It defaults to 0.
type Page int

// Resource is a general descriptor for local files or Last.fm URLs.
type Resource struct {
	// TODO replace string with custom types
	main   string
	method string
	name   string
	page   Page
	params map[string]string
}

// NewUserInfo returns the Resource for "user.getInfo".
func NewUserInfo(name string) *Resource {
	rsrc := new(Resource)
	rsrc.main = "user"
	rsrc.method = "getInfo"
	rsrc.name = name
	rsrc.params = make(map[string]string)
	return rsrc
}

// NewUserRecentTracks returns the Resource for "user.getRecentTracks".
func NewUserRecentTracks(name string, page Page, time Midnight) *Resource {
	rsrc := new(Resource)
	rsrc.main = "user"
	rsrc.method = "getRecentTracks"
	rsrc.name = name
	rsrc.page = page
	time -= time % 86400
	rsrc.params = make(map[string]string)
	rsrc.params["from"] = strconv.Itoa(int(time) - 1)
	rsrc.params["to"] = strconv.Itoa(int(time) + 86400)
	return rsrc
}

// NewArtistInfo returns the Resource for "artist.getInfo".
func NewArtistInfo(name string) *Resource {
	rsrc := new(Resource)
	rsrc.main = "artist"
	rsrc.method = "getInfo"
	rsrc.name = name
	rsrc.params = make(map[string]string)
	return rsrc
}
