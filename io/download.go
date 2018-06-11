package io

import (
	"fmt"
	"net/http"
	"net/url"
)

// FmtURL formats the Last.FM URL.
func FmtURL(rsrc *Resource, apiKey string) string {
	base := "http://ws.audioscrobbler.com/2.0/"
	params := "?format=json&api_key=%v&method=%v.%v&%v=%v"

	url := base + fmt.Sprintf(params, apiKey,
		rsrc.main, rsrc.method, rsrc.main, url.PathEscape(string(rsrc.name)))

	if rsrc.page > 0 {
		url += fmt.Sprintf("&page=%d", int(rsrc.page))
	}

	if rsrc.time > -1 {
		url += fmt.Sprintf("&from=%d&to=%d",
			int(rsrc.time)-1, int(rsrc.time)+86400)
	}

	return url
}

// Download retrieves the content of the given URL.
func Download(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
