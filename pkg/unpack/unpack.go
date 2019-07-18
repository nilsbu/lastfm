package unpack

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type resource interface {
	locator() rsrc.Locator
}

type deserializer interface {
	deserializer() interface{}
	interpret(raw interface{}) (interface{}, error)
}

type serializer interface {
	raw(obj interface{}) interface{}
}

type obtainer interface {
	resource
	deserializer
}

func deserialize(o deserializer, data []byte) (interface{}, error) {
	raw := o.deserializer()
	if err := json.Unmarshal(data, raw); err != nil {
		return nil, errors.Wrap(err, "could not deserialize")
	}

	return o.interpret(raw)
}

func obtain(o obtainer, r rsrc.Reader) (interface{}, error) {
	data, err := r.Read(o.locator())
	if err != nil {
		return nil, err
	}

	if errMsg, err := deserialize(&obError{}, data); err == nil {
		d := errMsg.(*LastfmError)
		if d.Code > 0 {
			return nil, d
		}
	}

	return deserialize(o, data)
}

type depositer interface {
	resource
	serializer
}

func deposit(obj interface{}, d depositer, w rsrc.Writer) error {
	raw := d.raw(obj)

	data, _ := json.Marshal(raw)

	return w.Write(data, d.locator())
}
