package io

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	return url
}

// download retrieves the content of the given URL.
func download(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, err
}
