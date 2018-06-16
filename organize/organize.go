package organize

import (
	"encoding/json"
	"errors"

	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/unpack"
)

// TODO name / what is this file

// LoadAPIKey loads an the API key.
func LoadAPIKey(r io.Reader) (key io.APIKey, err error) {
	rsrc := io.NewAPIKey()
	data, err := r.Read(rsrc)
	if err != nil {
		return
	}

	unm := &unpack.APIKey{}
	err = json.Unmarshal(data, unm)
	if err != nil {
		return
	}
	if unm.Key == "" {
		return "", errors.New("No valid API key was read")
	}

	return io.APIKey(unm.Key), nil
}

// WriteAllDayPlays writes a list of day plays.
func WriteAllDayPlays(
	plays []unpack.DayPlays,
	name io.Name,
	w io.Writer) (err error) {
	jsonData, _ := json.Marshal(plays)
	err = w.Write(jsonData, io.NewAllDayPlays(name))
	return
}

// ReadAllDayPlays reads a list of day plays.
func ReadAllDayPlays(
	name io.Name,
	r io.Reader) (plays []unpack.DayPlays, err error) {
	jsonData, err := r.Read(io.NewAllDayPlays(name))
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &plays)
	return
}
