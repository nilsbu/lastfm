package command

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nilsbu/lastfm/pkg/unpack"
)

func resolve(args []string, session *unpack.SessionInfo) (cmd command, err error) {
	if len(args) < 1 {
		return nil, errors.New("args does not contain the program name")
	}

	first, params := args[0], args[1:]

	switch first {
	case "lastfm":
		return resolveLastfm(params, session)
	default:
		return nil, fmt.Errorf("program '%v' is not supported", first)
	}
}

func resolveLastfm(
	params []string, session *unpack.SessionInfo) (cmd command, err error) {
	if len(params) < 1 {
		return help{}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	case "help":
		return help{}, nil
	case "session":
		return resolveSession(params, session)
	case "update":
		return resolveUpdate(params, session)
	case "print":
		return resolvePrint(params, session)
	default:
		return nil, fmt.Errorf("command '%v' is not supported", first)
	}
}

func resolveSession(
	params []string, session *unpack.SessionInfo) (cmd command, err error) {
	if len(params) < 1 {
		return sessionInfo{session}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	case "info":
		if len(params) > 0 {
			return nil, errors.New("'session info' requires no further parameters")
		}
		return sessionInfo{session}, nil
	case "start":
		if len(params) < 1 {
			return nil, errors.New("'session start' requires a user name")
		} else if len(params) > 1 {
			return nil, errors.New("params %v are superfluous")
		}
		return sessionStart{session: session, user: params[0]}, nil
	case "stop":
		if len(params) > 0 {
			return nil, errors.New("'session stop' requires no further parameters")
		}
		return sessionStop{session}, nil
	default:
		return nil, fmt.Errorf("parameter '%v' is not supported", first)
	}
}

func resolveUpdate(
	params []string, session *unpack.SessionInfo) (cmd command, err error) {
	if session == nil {
		return nil, errors.New("'update' requires a running session or further parameters")
	}

	if len(params) < 1 {
		return updateHistory{session}, nil
	}

	first, params := params[0], params[1:]

	switch first {
	// TODO more update commands
	default:
		return nil, fmt.Errorf("parameter '%v' is not supported", first)
	}
}

func resolvePrint(
	params []string, session *unpack.SessionInfo) (cmd command, err error) {
	if len(params) < 1 {
		return nil, errors.New("'print' requires further parameters")
	}

	if session == nil {
		// TODO should not be needed when no user is required
		return nil, errors.New("'print' requires a running session")
	}

	first, params := params[0], params[1:]

	switch first {
	case "total":
		if len(params) < 1 {
			return printTotal{session: session}, nil
		} else if len(params) == 1 {
			n, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, fmt.Errorf("'%v' must be an int", params[0])
			}
			return printTotal{session: session, n: n}, nil
		} else {
			return nil, errors.New(
				"'print total' accepts no more than one additional parameter")
		}
	case "fade":
		if len(params) < 1 {
			return nil, errors.New("'print fade' needs one more additional parameter")
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
		return printFade{session: session, hl: hl, n: n}, nil
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
