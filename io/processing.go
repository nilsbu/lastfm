package io

import (
	"encoding/json"
)

func getPages(data string) (page int, err error) {
	dat := UserRecentTracks{}
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		return 0, err
	}

	return dat.Recenttracks.Attr.TotalPages, nil
}
