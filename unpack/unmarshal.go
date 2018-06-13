package unpack

import "encoding/json"

// UnmarshalUserRecentTracks unmarshales a JSON result of user.getRecentTracks.
func UnmarshalUserRecentTracks(data []byte) (urt *UserRecentTracks, err error) {
	urt = &UserRecentTracks{}
	err = json.Unmarshal([]byte(data), urt)
	return
}

// UnmarshalAPIKey unmarshals an API key from JSON data.
func UnmarshalAPIKey(data []byte) (key *APIKey, err error) {
	key = &APIKey{}
	err = json.Unmarshal([]byte(data), key)
	return
}
