package unpack

import (
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/pkg/errors"
)

type obAPIKey struct{}

// LoadAPIKey loads the API key. It returns an error if the key could not be
// read or is invalid.
func LoadAPIKey(r rsrc.Reader) (key string, err error) {
	data, err := obtain(obAPIKey{}, r)
	if err != nil {
		return "", err
	}
	key = data.(string)
	return key, nil
}

func (o obAPIKey) locator() rsrc.Locator {
	return rsrc.APIKey()
}

func (o obAPIKey) deserializer() interface{} {
	return &jsonAPIKey{}
}

func (o obAPIKey) interpret(raw interface{}) (interface{}, error) {
	key := raw.(*jsonAPIKey)
	if err := rsrc.CheckAPIKey(key.Key); err != nil {
		return "", errors.Wrap(err, "API key could not be read")
	}

	return key.Key, nil
}
