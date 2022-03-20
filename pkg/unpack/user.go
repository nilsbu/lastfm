package unpack

import (
	"fmt"
	"strconv"

	"github.com/nilsbu/lastfm/pkg/charts2"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type obBookmark struct {
	user string
}

// WriteBookmark writes a bookmark.
func WriteBookmark(bookmark rsrc.Day, user string, w rsrc.Writer) error {
	return deposit(bookmark, obBookmark{user}, w)
}

// LoadBookmark loads a bookmark.
func LoadBookmark(user string, r rsrc.Reader) (rsrc.Day, error) {
	data, err := obtain(&obBookmark{user}, r)
	if err != nil {
		return nil, err
	}
	return data.(rsrc.Day), nil
}

func (o obBookmark) locator() rsrc.Locator {
	return rsrc.Bookmark(o.user)
}

func (o obBookmark) deserializer() interface{} {
	return &jsonBookmark{}
}

func (o obBookmark) interpret(raw interface{}) (interface{}, error) {
	bookmark := raw.(*jsonBookmark)
	return rsrc.ParseDay(bookmark.NextDay), nil
}

func (o obBookmark) raw(obj interface{}) interface{} {
	t := obj.(rsrc.Day).Time()

	js := jsonBookmark{
		NextDay: fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()),
	}
	return js
}

type obAllDayPlays struct {
	user string
}

// LoadAllDayPlays loads the pre-processed history of a user, called alldayplays.
func LoadAllDayPlays(user string, r rsrc.Reader) ([]map[string]float64, error) {
	data, err := obtain(obAllDayPlays{user}, r)
	if err != nil {
		return nil, err
	}

	plays := data.([]map[string]float64)
	days := make([]map[string]float64, len(plays))

	for i := range plays {
		day := map[string]float64{}
		for name, value := range plays[i] {
			day[name] = value
		}
		days[i] = day
	}

	return days, nil
}

// WriteAllDayPlays writed the pre-processed history of a user.
func WriteAllDayPlays(days []map[string]float64, user string, w rsrc.Writer) error {
	plays := make([]map[string]float64, len(days))
	for i := range days {
		day := map[string]float64{}
		for name, values := range days[i] {
			day[name] = values
		}
		plays[i] = day
	}

	return deposit(plays, obAllDayPlays{user}, w)
}

func (o obAllDayPlays) locator() rsrc.Locator {
	return rsrc.AllDayPlays(o.user)
}

func (o obAllDayPlays) deserializer() interface{} {
	return &[]map[string]float64{}
}

func (o obAllDayPlays) interpret(raw interface{}) (interface{}, error) {
	return *raw.(*[]map[string]float64), nil
}

func (o obAllDayPlays) raw(obj interface{}) interface{} {
	return obj
}

type obDayHistory struct {
	user string
	day  rsrc.Day
}

// LoadDayHistory loads the pre-processed history of a user for a single day, called history.
func LoadDayHistory(user string, day rsrc.Day, r rsrc.Reader) ([]charts2.Song, error) {
	data, err := obtain(obDayHistory{user, day}, r)
	if err != nil {
		return nil, err
	}

	inSongs := data.([][]string)
	outSongs := make([]charts2.Song, len(inSongs))

	for i, song := range inSongs {
		duration, err := strconv.ParseFloat(song[3], 64)
		if err != nil {
			return nil, err
		}
		outSongs[i] = charts2.Song{
			Artist:   song[0],
			Title:    song[1],
			Album:    song[2],
			Duration: duration,
		}
	}

	return outSongs, nil
}

// WritDayHistory write the pre-processed history of a user for a single day.
func WriteDayHistory(songs []charts2.Song, user string, day rsrc.Day, w rsrc.Writer) error {
	outSongs := make([][]string, len(songs))
	for i, song := range songs {
		outSongs[i] = []string{song.Artist, song.Title, song.Album, fmt.Sprintf("%f", song.Duration)}
	}

	return deposit(outSongs, obDayHistory{user, day}, w)
}

func (o obDayHistory) locator() rsrc.Locator {
	return rsrc.DayHistory(o.user, o.day)
}

func (o obDayHistory) deserializer() interface{} {
	return &[][]string{}
}

func (o obDayHistory) interpret(raw interface{}) (interface{}, error) {
	return *raw.(*[][]string), nil
}

func (o obDayHistory) raw(obj interface{}) interface{} {
	return obj
}

type obCorrections struct {
	user string
	fn   func(string) rsrc.Locator
}

// LoadArtistCorrections loads corrections for artist names. The result is a map
// with the false names as keys and correct names as values.
func LoadArtistCorrections(user string, r rsrc.Reader,
) (map[string]string, error) {
	data, err := obtain(obCorrections{user, rsrc.ArtistCorrections}, r)
	if err != nil {
		return nil, err
	}

	corr := data.(map[string]string)
	return corr, nil
}

// LoadSupertagCorrections loads corrections for artist's supertags. The result
// is a map with the artist names as keys and intended supertags as values.
func LoadSupertagCorrections(user string, r rsrc.Reader,
) (map[string]string, error) {
	return loadTagCorrections(user, r, rsrc.SupertagCorrections)
}

// LoadCountryCorrections loads corrections for artist's country. The result
// is a map with the artist names as keys and intended supertags as values.
func LoadCountryCorrections(user string, r rsrc.Reader,
) (map[string]string, error) {
	return loadTagCorrections(user, r, rsrc.CountryCorrections)
}

func loadTagCorrections(
	user string,
	r rsrc.Reader,
	corrections func(string) rsrc.Locator,
) (map[string]string, error) {
	data, err := obtain(obCorrections{user, corrections}, r)
	if err != nil {
		return nil, err
	}

	corr := data.(map[string]string)
	return corr, nil
}

func (o obCorrections) locator() rsrc.Locator {
	return o.fn(o.user)
}

func (o obCorrections) deserializer() interface{} {
	return &jsonCorrections{}
}

func (o obCorrections) interpret(raw interface{}) (interface{}, error) {
	key := raw.(*jsonCorrections)

	return key.Corrections, nil
}
