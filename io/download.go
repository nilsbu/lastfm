package io

import (
	"fmt"
	"net/http"
	"strconv"
)

// FmtURL formats the Last.FM URL.
func FmtURL(rsrc *Resource, apiKey string) string {
	base := "http://ws.audioscrobbler.com/2.0/"
	params := "?format=json&api_key=%v&method=%v.%v&%v=%v"

	url := base + fmt.Sprintf(params, apiKey, rsrc.main, rsrc.method, rsrc.main, rsrc.name)

	if rsrc.page > 0 {
		url += "&page=" + strconv.Itoa(int(rsrc.page))
	}

	for k, v := range rsrc.params {
		url += fmt.Sprintf("&%v=%v", k, v)
	}

	return url
}

// Download retrieves the content of the given URL.
func Download(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
