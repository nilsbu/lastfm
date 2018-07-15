package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nilsbu/lastfm/pkg/unpack"
)

type node struct {
	cmd   *cmd
	nodes nodes
}

type nodes map[string]node

type cmd struct {
	descr   string
	get     func(params []interface{}, opts map[string]interface{}) command
	params  params
	options options
	session bool
}

type param struct {
	name  string
	descr string
	kind  string // TODO
}

type params []*param

type option struct {
	param
	value string
}

type options map[string]*option

var cmdRoot = node{
	nodes: map[string]node{
		"lastfm": cmdLastfm,
	},
}

var cmdLastfm = node{
	cmd: exeHelp,
	nodes: map[string]node{
		"cache":   cmdCache,
		"help":    cmdHelp,
		"print":   cmdPrint,
		"session": cmdSession,
		"update":  cmdUpdate,
	},
}

var cmdCache = node{
	cmd: &cmd{
		descr: "runs cache server",
		get: func(params []interface{}, opts map[string]interface{}) command {
			return cache{opts["p"].(int)}
		},
		options: options{"p": optCachePort},
	},
}

var cmdHelp = node{
	cmd: exeHelp,
}

var cmdPrint = node{
	nodes: nodes{
		"fade":  node{cmd: exePrintFade},
		"tags":  node{cmd: exePrintTags},
		"total": node{cmd: exePrintTotal},
	},
}

var cmdSession = node{
	cmd: exeSessionInfo,
	nodes: nodes{
		"info":  node{cmd: exeSessionInfo},
		"start": node{cmd: exeSessionStart},
		"stop":  node{cmd: exeSessionStop},
	},
}

var cmdUpdate = node{
	cmd: &cmd{
		descr: "updates a user's history",
		get: func(params []interface{}, opts map[string]interface{}) command {
			return updateHistory{}
		},
		session: true,
	},
}

var exeHelp = &cmd{
	descr: "gives help",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return help{}
	},
}

var exePrintFade = &cmd{
	descr: "prints a user's top artists in fading charts",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printFade{
			by:   opts["by"].(string),
			name: opts["name"].(string),
			n:    opts["n"].(int),
			hl:   params[0].(float64),
		}
	},
	params: params{&param{
		"half-life",
		"span of days over which a 'scrobble' loses half its value",
		"float",
	}},
	options: options{
		"by":   optChartType,
		"name": optGenericName,
		"n":    optArtistCount,
	},
	session: true,
}

var exePrintTags = &cmd{
	descr: "prints the top tags of an artist",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printTags{params[0].(string)}
	},
	params: params{parArtistName},
}

var exePrintTotal = &cmd{
	descr: "prints a user's top artists by total number of plays",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printTotal{
			by:   opts["by"].(string),
			name: opts["name"].(string),
			n:    opts["n"].(int),
		}
	},
	options: options{
		"by":   optChartType,
		"name": optGenericName,
		"n":    optArtistCount,
	},
	session: true,
}

var exeSessionInfo = &cmd{
	descr: "report information about the currently running session, if one is running",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return sessionInfo{}
	},
}

var exeSessionStart = &cmd{
	descr: "starts a session if none is currently running",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return sessionStart{params[0].(string)}
	},
	params: params{parUserName},
}

var exeSessionStop = &cmd{
	descr: "stops the currently running session if none is currently running",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return sessionStop{}
	},
}

var parPort = &param{
	"port",
	"a TCP port number",
	"int",
}

var parUserName = &param{
	"user name",
	"a Last.fm user name",
	"string",
}

var parArtistName = &param{
	"artist name",
	"the name of an artist",
	"string",
}

var optCachePort = &option{
	*parPort,
	"14003",
}

var optChartType = &option{
	param{"chart type",
		"'all' or 'super'",
		"string"}, // TODO make something like an enum
	"all",
}

var optGenericName = &option{
	param{"name",
		"some name", // TODO be more specific
		"string"},
	"",
}

var optArtistCount = &option{
	param{"count",
		"number of artists",
		"int"},
	"10",
}

func resolve(args []string, session *unpack.SessionInfo) (cmd command, err error) {
	return resolveTree(args, session, cmdRoot)
}

func resolveTree(
	args []string,
	session *unpack.SessionInfo,
	tree node,
) (command, error) {
	if len(args) > 0 && args[0][0] != '-' {
		if cmd, ok := tree.nodes[args[0]]; ok {
			return resolveTree(args[1:], session, cmd)
		}
	}

	if tree.cmd == nil {
		// TODO more details
		return nil, errors.New("command does not exist, are more arguments missing?")
	}

	if tree.cmd.session && session == nil {
		return nil, errors.New("command can only be executed when a session is running")
	}

	params, opts, err := parseArguments(args, tree.cmd)
	if err != nil {
		return nil, err
	}

	return tree.cmd.get(params, opts), nil

}

func parseArguments(args []string, cmd *cmd,
) (params []interface{},
	opts map[string]interface{},
	err error) {
	if len(args) < len(cmd.params) {
		return nil, nil, errors.New("too few params")
	}

	params = make([]interface{}, len(cmd.params))

	for i := 0; i < len(cmd.params); i++ {
		value, err := parseArgument(args[i], cmd.params[i].kind)
		if err != nil {
			return nil, nil, err
		}

		params[i] = value
	}

	rawOpts := make(map[string]string)
	opts = make(map[string]interface{})

	for i := len(cmd.params); i < len(args); i++ {
		if args[i][0] != '-' {
			return nil, nil, fmt.Errorf("parameter '%v' is unexpected", args[i])
		}

		idx := strings.Index(args[i], "=")
		if idx < 0 {
			return nil, nil,
				fmt.Errorf("option must be of format '-key=value', '%v' is not", args[i])
		}

		key := args[i][1:idx]
		_, ok := cmd.options[key]
		if !ok {
			return nil, nil, fmt.Errorf("option '%v' is not supported", key)
		}

		rawOpts[key] = args[i][idx+1:]
	}

	for key, opt := range cmd.options {
		if _, ok := rawOpts[key]; !ok {
			rawOpts[key] = opt.value
		}
	}

	for key, raw := range rawOpts {
		value, err := parseArgument(raw, cmd.options[key].kind)
		if err != nil {
			return nil, nil, err
		}

		opts[key] = value
	}

	return params, opts, nil
}

func parseArgument(arg, kind string) (value interface{}, err error) {
	switch kind {
	case "float":
		value, err = strconv.ParseFloat(arg, 64)
	case "int":
		value, err = strconv.Atoi(arg)
	case "string":
		value = arg
	default:
		// Cannot be reached
	}

	return
}
