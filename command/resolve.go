package command

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/rsrc"
)

func resolve(args []string) (cmd command, err error) {
	if len(args) < 1 {
		return nil, errors.New("args does not contain the program name")
	}

	first, params := args[0], args[1:]

	switch first {
	case "lastfm":
		return resolveLastfm(params)
	default:
		return nil, fmt.Errorf("program '%v' is not supported", first)
	}
}

func resolveLastfm(params []string) (cmd command, err error) {
	if len(params) < 1 {
		return help{}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	case "help":
		return help{}, nil
	case "session":
		return resolveSession(params)
	default:
		return nil, fmt.Errorf("command '%v' is not supported", first)
	}
}

func resolveSession(params []string) (cmd command, err error) {
	if len(params) < 1 {
		return sessionInfo{}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	case "info":
		if len(params) > 0 {
			return nil, errors.New("'session info' requires no further parameters")
		}
		return sessionInfo{}, nil
	case "start":
		if len(params) < 1 {
			return nil, errors.New("'session start' requires a user name")
		}
		if len(params) > 1 {
			return nil, errors.New("params %v are superfluous")
		}
		return sessionStart{user: rsrc.Name(params[0])}, nil
	case "stop":
		if len(params) > 0 {
			return nil, errors.New("'session stop' requires no further parameters")
		}
		return sessionStop{}, nil
	default:
		return nil, fmt.Errorf("parameter '%v' is not supported", first)
	}
}
