package io

// Midnight is a unix time at midnight. It defaults to -1.
type Midnight int64

// Page is a page of a multi-page resource. It defaults to 0.
type Page int

// Name is the name of a user, artist or tag
type Name string

// Domain is used to differentiate between different kinds of resources.
type Domain string

// List of valid Domains.
const (
	Raw    Domain = "data"
	User          = "user"
	Util          = "util"
	Global        = "global"
)

// Resource is a general descriptor for local files or Last.fm URLs.
type Resource struct {
	domain Domain
	// TODO replace string with custom types.
	// TODO main is not a precise name.
	main   string
	method string
	name   Name
	page   Page
	time   Midnight
}

// NewUserInfo returns the Resource for "user.getInfo".
func NewUserInfo(name Name) *Resource {
	return &Resource{
		domain: Raw,
		main:   "user",
		method: "getInfo",
		name:   name,
		time:   -1,
	}
}

// NewUserRecentTracks returns the Resource for "user.getRecentTracks".
func NewUserRecentTracks(name Name, page Page, time Midnight) *Resource {
	return &Resource{
		domain: Raw,
		main:   "user",
		method: "getRecentTracks",
		name:   name,
		page:   page,
		time:   time - time%86400,
	}
}

// NewArtistInfo returns the Resource for "artist.getInfo".
func NewArtistInfo(name Name) *Resource {
	return &Resource{
		domain: Raw,
		main:   "artist",
		method: "getInfo",
		name:   name,
		time:   -1,
	}
}

// NewAPIKey returns the Resource for "artist.getInfo".
func NewAPIKey() *Resource {
	return &Resource{
		domain: Util,
		method: "apikey",
		time:   -1,
	}
}

func NewAllDayPlays(user Name) *Resource {
	return &Resource{
		domain: User,
		name:   user,
		method: "alldayplays",
		time:   -1,
	}
}
