package io

import (
	"fmt"
	"net/url"
	"strings"
)

// APIKey is an API key for Last.fm.
type APIKey string

// fmtURL formats the Last.FM URL.
func fmtURL(rsrc *Resource, key APIKey) string {
	base := "http://ws.audioscrobbler.com/2.0/"
	params := "?format=json&api_key=%v&method=%v.%v&%v=%v"

	url := base + fmt.Sprintf(params, key,
		rsrc.main, rsrc.method, rsrc.main, url.PathEscape(string(rsrc.name)))

	if rsrc.page > 0 {
		url += fmt.Sprintf("&page=%d", int(rsrc.page))
	}

	if rsrc.time > -1 {
		url += fmt.Sprintf("&from=%d&to=%d",
			int(rsrc.time)-1, int(rsrc.time)+86400)
	}

	if rsrc.method == "getRecentTracks" {
		url += "&limit=200"
	}

	return url
}

func escapeBadNames(name Name) Name {
	bad := [13]string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4",
		"LPT1", "LPT2", "LPT3", "LPT4", "LST"}

	upperName := strings.ToUpper(string(name))
	for _, s := range bad {
		if upperName == s {
			return Name("_") + name
		}
	}

	return name
}

func parseForPath(name Name) Name {
	escaped := url.PathEscape(string(name))
	escaped = strings.Replace(escaped, "%20", "+", -1)
	escaped = strings.Replace(escaped, "/", "+", -1)
	return escapeBadNames(Name(escaped))
}

// fmtPath returns the relative path for a resource.
func fmtPath(rsrc *Resource) string {
	path := fmt.Sprintf(".lastfm/%v/", rsrc.domain)
	if rsrc.domain == Raw {
		path += fmt.Sprintf("%v.%v/%v",
			rsrc.main, rsrc.method, parseForPath(rsrc.name))

		if rsrc.time > -1 {
			path += fmt.Sprintf(".%d", rsrc.time)
		}
		if rsrc.page > 0 {
			path += fmt.Sprintf("(%v)", rsrc.page)
		}
	} else if rsrc.domain == Util {
		path += rsrc.method
	} else if rsrc.domain == User {
		path += fmt.Sprintf("%v/%v", rsrc.name, rsrc.method)
	}

	return path + ".json"
}
