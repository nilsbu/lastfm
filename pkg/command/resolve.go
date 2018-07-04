package command

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nilsbu/lastfm/pkg/organize"
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
	case "print":
		return resolvePrint(params, sid)
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
		return sessionStart{sid: sid, user: params[0]}, nil
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

func resolvePrint(
	params []string, sid organize.SessionID) (cmd command, err error) {
	if len(params) < 1 {
		return nil, errors.New("'print' requires further parameters")
	}

	if sid == "" {
		// TODO should not be needed when no user is required
		return nil, errors.New("'print' requires a running session")
	}

	first, params := params[0], params[1:]

	switch first {
	case "total":
		if len(params) < 1 {
			return printTotal{sid: sid}, nil
		} else if len(params) == 1 {
			n, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, fmt.Errorf("'%v' must be an int", params[0])
			}
			return printTotal{sid: sid, n: n}, nil
		} else {
			return nil, errors.New(
				"'print total' accepts no more than one additional parameter")
		}
	case "fade":
		if len(params) < 1 {
			return nil, errors.New("'print total' accepts no more than one additional parameter")
		} else if len(params) > 2 {
			return nil, errors.New(
				"'print total' accepts no more than two additional parameter")
		}

		hl, err := strconv.ParseFloat(params[0], 64)
		if err != nil {
			return nil, fmt.Errorf("'%v' must be a float", params[0])
		}

		var n int
		if len(params) == 2 {
			n, err = strconv.Atoi(params[1])
			if err != nil {
				return nil, fmt.Errorf("'%v' must be an int", params[1])
			}
		}
		return printFade{sid: sid, hl: hl, n: n}, nil
	case "tags":
		if len(params) == 1 {
			return printTags{params[0]}, nil
		} else {
			return nil, errors.New("'print tags' requires exactly one parameter")
		}
	default:
		return nil, fmt.Errorf("parameter '%v' is not supported", first)
	}
}
