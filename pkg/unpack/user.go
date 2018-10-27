package unpack

import (
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type obAllDayPlays struct {
	user string
}

// LoadAllDayPlays loads the pre-processed history of a user, called alldayplays.
func LoadAllDayPlays(user string, r rsrc.Reader) ([]charts.Charts, error) {
	data, err := obtain(obAllDayPlays{user}, r)
	if err != nil {
		return nil, err
	}

	plays := data.([]map[string]float64)
	days := make([]charts.Charts, len(plays))

	for i := range plays {
		day := charts.Charts{}
		for name, value := range plays[i] {
			day[name] = []float64{value}
		}
		days[i] = day
	}

	return days, nil
}

// WriteAllDayPlays writed the pre-processed history of a user.
func WriteAllDayPlays(days []charts.Charts, user string, w rsrc.Writer) error {
	plays := make([]map[string]float64, len(days))
	for i := range days {
		day := map[string]float64{}
		for name, values := range days[i] {
			day[name] = values[0]
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
	data, err := obtain(obCorrections{user, rsrc.SupertagCorrections}, r)
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
