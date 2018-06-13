package unpack

import "encoding/json"

// NewUserRecentTracks unmarshales a JSON result of user.getRecentTracks.
func NewUserRecentTracks(data []byte) (urt *UserRecentTracks, err error) {
	urt = &UserRecentTracks{}
	err = json.Unmarshal([]byte(data), urt)
	return
}
