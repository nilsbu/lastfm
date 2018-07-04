package unpack

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type obUserInfo struct {
	name string
}

// LoadUserInfo loads a user's registration date. It is returned along with the
// name.
func LoadUserInfo(name string, r rsrc.Reader) (*User, error) {
	data, err := obtain(&obUserInfo{name}, r)
	if err != nil {
		return nil, err
	}
	user := data.(*User)
	return user, nil
}

func (o *obUserInfo) locator() rsrc.Locator {
	return rsrc.UserInfo(o.name)
}

func (o *obUserInfo) deserializer() interface{} {
	return &jsonUserInfo{}
}

func (o *obUserInfo) interpret(raw interface{}) (interface{}, error) {
	ui := raw.(*jsonUserInfo)

	utc := ui.User.Registered.UTC
	return &User{ui.User.Name, rsrc.ToDay(utc)}, nil
}
