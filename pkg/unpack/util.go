package unpack

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/nilsbu/lastfm/pkg/rsrc"
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

// SessionInfo contains information about a running session.
type SessionInfo struct {
	User    string
	Options map[string]string
}

type obSessionInfo struct{}

// LoadSessionInfo loads information about a session, if one is
// running.
func LoadSessionInfo(r rsrc.Reader) (*SessionInfo, error) {
	data, err := obtain(obSessionInfo{}, r)
	if err != nil {
		return nil, err
	}
	session := data.(*SessionInfo)
	return session, nil
}

func WriteSessionInfo(session *SessionInfo, w rsrc.Writer) error {
	return deposit(session, obSessionInfo{}, w)
}

func (o obSessionInfo) locator() rsrc.Locator {
	return rsrc.SessionInfo()
}

func (o obSessionInfo) deserializer() interface{} {
	return &jsonSessionInfo{}
}

func (o obSessionInfo) interpret(raw interface{}) (interface{}, error) {
	session := raw.(*jsonSessionInfo)
	if session.User == "" {
		return "", errors.New("could not read session")
	}

	options := make(map[string]string)
	for _, opt := range session.Options {
		options[opt.Name] = opt.Value
	}

	return &SessionInfo{User: session.User, Options: options}, nil
}

func (o obSessionInfo) raw(obj interface{}) interface{} {
	session := obj.(*SessionInfo)

	options := []jsonSessionOption{}
	if session.Options != nil {
		for k, v := range session.Options {
			options = append(options, jsonSessionOption{Name: k, Value: v})
		}
	}
	// sort mainly for test stability
	sort.Slice(options, func(i, j int) bool { return options[i].Name < options[j].Name })

	return &jsonSessionInfo{User: session.User, Options: options}
}
