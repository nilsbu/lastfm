package unpack

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type obAllDayPlays struct {
	user string
}

// LoadAllDayPlays loads the pre-processed history of a user, called alldayplays.
func LoadAllDayPlays(user string, r rsrc.Reader) ([]PlayCount, error) {
	data, err := obtain(obAllDayPlays{user}, r)
	if err != nil {
		return nil, err
	}
	plays := data.([]PlayCount)
	return plays, nil
}

// WriteAllDayPlays writed the pre-processed history of a user.
func WriteAllDayPlays(plays []PlayCount, user string, w rsrc.Writer) error {
	return deposite(plays, obAllDayPlays{user}, w)
}

func (o obAllDayPlays) locator() rsrc.Locator {
	return rsrc.AllDayPlays(o.user)
}

func (o obAllDayPlays) deserializer() interface{} {
	return &[]PlayCount{}
}

func (o obAllDayPlays) interpret(raw interface{}) (interface{}, error) {
	return *raw.(*[]PlayCount), nil
}

func (o obAllDayPlays) raw(obj interface{}) interface{} {
	return obj
}
