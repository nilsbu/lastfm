package unpack

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type obtainer interface {
	locator() rsrc.Locator
	deserializer() interface{}
	interpret(raw interface{}) (interface{}, error)
}

func obtain(o obtainer, r rsrc.Reader) (interface{}, error) {
	data, err := r.Read(o.locator())
	if err != nil {
		return nil, err
	}

	raw := o.deserializer()
	err = json.Unmarshal(data, raw)
	if err != nil {
		return nil, errors.Wrap(err, "could not deserialize")
	}

	return o.interpret(raw)
}
