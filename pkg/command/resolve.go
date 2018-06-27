package command

import (
	"errors"
	"fmt"

	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func resolve(args []string, sid organize.SessionID) (cmd command, err error) {
	if len(args) < 1 {
		return nil, errors.New("args does not contain the program name")
	}

	first, params := args[0], args[1:]

	switch first {
	case "lastfm":
		return resolveLastfm(params, sid)
	default:
		return nil, fmt.Errorf("program '%v' is not supported", first)
	}
}

func resolveLastfm(
	params []string, sid organize.SessionID) (cmd command, err error) {
	if len(params) < 1 {
		return help{}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	case "help":
		return help{}, nil
	case "session":
		return resolveSession(params, sid)
	case "update":
		return resolveUpdate(params, sid)
	default:
		return nil, fmt.Errorf("command '%v' is not supported", first)
	}
}

func resolveSession(
	params []string, sid organize.SessionID) (cmd command, err error) {
	if len(params) < 1 {
		return sessionInfo{sid}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	case "info":
		if len(params) > 0 {
			return nil, errors.New("'session info' requires no further parameters")
		}
		return sessionInfo{sid}, nil
	case "start":
		if len(params) < 1 {
			return nil, errors.New("'session start' requires a user name")
		} else if len(params) > 1 {
			return nil, errors.New("params %v are superfluous")
		}
		return sessionStart{sid: sid, user: rsrc.Name(params[0])}, nil
	case "stop":
		if len(params) > 0 {
			return nil, errors.New("'session stop' requires no further parameters")
		}
		return sessionStop{sid}, nil
	default:
		return nil, fmt.Errorf("parameter '%v' is not supported", first)
	}
}

func resolveUpdate(
	params []string, sid organize.SessionID) (cmd command, err error) {
	if sid == "" {
		return nil, errors.New("'update' requires a running session or further parameters")
	}

	if len(params) < 1 {
		return updateHistory{sid}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	// TODO more update commands
	default:
		return nil, fmt.Errorf("parameter '%v' is not supported", first)
	}
}
