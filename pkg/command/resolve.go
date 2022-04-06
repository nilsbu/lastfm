package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	kind  string
}

type params []*param

type option struct {
	param
	value string
}

type options map[string]*option

var cmdRoot = node{
	nodes: map[string]node{
		"lastfm":     cmdLastfm,
		"lastfm-csv": cmdLastfm,
		"lastfm-srv": cmdLastfm,
	},
}

var cmdLastfm = node{
	cmd: exeHelp,
	nodes: map[string]node{
		"help":     cmdHelp,
		"print":    cmdPrint,
		"session":  cmdSession,
		"table":    cmdTable,
		"timeline": {cmd: exeTimeline},
		"update":   cmdUpdate,
	},
}

var cmdHelp = node{
	cmd: exeHelp,
}

var cmdPrint = node{
	nodes: nodes{
		"fade":     node{cmd: exePrintFade},
		"period":   node{cmd: exePrintPeriod},
		"interval": node{cmd: exePrintInterval},
		"fademax":  node{cmd: exePrintFadeMax},
		"tags":     node{cmd: exePrintTags},
		"total":    node{cmd: exePrintTotal},
		"periods":  node{cmd: exePrintPeriods},
		"fades":    node{cmd: exePrintFades},
		// TODO Add command for raw execution
	},
}

var cmdTable = node{
	nodes: nodes{
		"fade":   node{cmd: exeTableFade},
		"period": node{cmd: exeTablePeriods},
		"total":  node{cmd: exeTableTotal},
	},
}

var cmdSession = node{
	cmd: exeSessionInfo,
	nodes: nodes{
		"info":   node{cmd: exeSessionInfo},
		"start":  node{cmd: exeSessionStart},
		"stop":   node{cmd: exeSessionStop},
		"config": node{cmd: exeSessionConfig},
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
	descr: "prints a user's top artists in fading charts", // TODO
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printFade{printCharts: printCharts{
			keys:       opts["keys"].(string),
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			hl:   params[0].(float64),
			date: opts["date"].(time.Time),
		}
	},
	params: params{parHL},
	options: options{
		"keys":       optChartsKeys,
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
		"date":       optDate,
	},
	session: true,
}

var exePrintPeriod = &cmd{
	descr: "", // TODO
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printPeriod{printCharts: printCharts{
			keys:       opts["keys"].(string),
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			period: params[0].(string),
		}
	},
	params: params{&param{
		"period",
		"", // TODO
		"string",
	}},
	options: options{
		"keys":       optChartsKeys,
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
	},
	session: true,
}

var exePrintInterval = &cmd{
	descr: "", // TODO
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printInterval{printCharts: printCharts{
			keys:       opts["keys"].(string),
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			begin:  params[0].(time.Time),
			before: params[1].(time.Time),
		}
	},
	params: params{&param{
		"begin",
		"", // TODO
		"time",
	}, &param{
		"before",
		"", // TODO
		"time",
	}},
	options: options{
		"keys":       optChartsKeys,
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
	},
	session: true,
}

var exePrintFadeMax = &cmd{
	descr: "", // TODO
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printFadeMax{printCharts: printCharts{
			keys:       opts["keys"].(string),
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: false, // Disabled since it makes no sense here
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			hl: params[0].(float64),
		}
	},
	params: params{parHL},
	options: options{
		"keys":       optChartsKeys,
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
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
	descr: "tables a user's top artists by total number of plays",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printTotal{printCharts: printCharts{
			keys:       opts["keys"].(string),
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			date: opts["date"].(time.Time),
		}
	},
	options: options{
		"keys":       optChartsKeys,
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
		"date":       optDate,
	},
	session: true,
}

var exePrintPeriods = &cmd{
	descr: "prints a user's top artists by total number of plays in the specified periods",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printPeriods{printCharts: printCharts{
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			period: params[0].(string),
		}
	},
	params: params{&param{
		"period",
		"period descriptor, format: '[0-9]*[yMd]'",
		"string",
	}},
	options: options{
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
	},
	session: true,
}

var exePrintFades = &cmd{
	descr: "prints a user's top artists by total number of plays in the specified periods",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printFades{printCharts: printCharts{
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			hl:     params[0].(float64),
			period: params[1].(string),
		}
	},
	params: params{parHL, &param{
		"period",
		"period descriptor, format: '[0-9]*[yMd]'",
		"string",
	}},
	options: options{
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
	},
	session: true,
}

var exeTimeline = &cmd{
	descr: "timeline of events",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return printTimeline{
			n: params[0].(int),
		}
	},
	params:  params{{"n", "top n artists will be evaluated each day", "int"}},
	session: true,
}

var exeTableFade = &cmd{
	descr: "tables a user's top artists in fading charts",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return tableFade{printCharts: printCharts{
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			hl:   params[0].(float64),
			step: opts["step"].(int),
		}
	},
	params: params{parHL},
	options: options{
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
		"step":       optStep,
	},
	session: true,
}

var exeTablePeriods = &cmd{
	descr: "tables a user's top artists by total number of plays in the specified periods",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return tablePeriods{printCharts: printCharts{
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			period: params[0].(string),
		}
	},
	params: params{&param{
		"period",
		"period descriptor, format: '[0-9]*[yMd]'",
		"string",
	}},
	options: options{
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
	},
	session: true,
}

var exeTableTotal = &cmd{
	descr: "tables a user's top artists by total number of plays",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return tableTotal{printCharts: printCharts{
			by:         opts["by"].(string),
			name:       opts["name"].(string),
			n:          opts["n"].(int),
			percentage: opts["%"].(bool),
			normalized: opts["normalized"].(bool),
			duration:   opts["duration"].(bool),
			entry:      opts["entry"].(float64),
		},
			step: opts["step"].(int),
		}
	},
	options: options{
		"by":         optChartType,
		"name":       optGenericName,
		"n":          optArtistCount,
		"%":          optChartsPercentage,
		"normalized": optChartsNormalized,
		"duration":   optChartsDuration,
		"entry":      optChartsEntry,
		"step":       optStep,
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

var exeSessionConfig = &cmd{
	descr: "set an option that will overwrite the defauls in future commands",
	get: func(params []interface{}, opts map[string]interface{}) command {
		return sessionConfig{
			option: params[0].(string),
			value:  params[1].(string),
		}
	},
	params: params{parOptionName, parOptionValue},
}

var parOptionName = &param{
	"option name",
	"a name of an option",
	"string",
}

var parOptionValue = &param{
	"user value",
	"a value of an option as a string",
	"string",
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

var parHL = &param{
	"half-life",
	"span of days over which a 'scrobble' loses half its value",
	"float",
}

// TODO name any key (see above) of option are duplicate
var optChartType = &option{
	param{"by",
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
	param{"n",
		"number of artists",
		"int"},
	"10",
}

var optChartsDuration = &option{
	param{"duration",
		"if charts are compiled by song duration",
		"bool"},
	"false",
}

var optChartsPercentage = &option{
	param{"%",
		"if charts are in percentage",
		"bool"},
	"false",
}

var optChartsKeys = &option{
	param{"keys",
		"keys of the charts ('artist' or 'song')",
		"string"},
	"",
}

var optChartsNormalized = &option{
	param{"normalized",
		"if charts are normalized",
		"bool"},
	"false",
}

var optChartsEntry = &option{
	param{"entry",
		"threshold for charts entry",
		"float"},
	"0",
}

var optDate = &option{
	param{"date",
		"a date in the format YYYY-MM-DD",
		"time"},
	"",
}

var optStep = &option{
	param{"step",
		"date step", // TODO
		"int"},
	"1",
}

var storableOptions = []*option{
	optChartType,
	optGenericName,
	optArtistCount,
	optChartsDuration,
	optChartsPercentage,
	optChartsKeys,
	optChartsNormalized,
	optChartsEntry,
	optDate,
	optStep,
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

	sessionOptions := map[string]string{}
	if session != nil && session.Options != nil {
		sessionOptions = session.Options
	}
	params, opts, err := parseArguments(args, tree.cmd, sessionOptions)
	if err != nil {
		return nil, err
	}

	return tree.cmd.get(params, opts), nil

}

func parseArguments(args []string, cmd *cmd, sessionOptions map[string]string,
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
			// special case: bool needs no '-k=v' format
			key := args[i][1:]
			if opt, ok := cmd.options[key]; ok && opt.param.kind == "bool" {
				rawOpts[key] = "true"
				continue
			} else {
				return nil, nil, fmt.Errorf(
					"option must be of format '-key=value', '%v' is not",
					args[i])
			}
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
			// fill missing options from session defaults
			if value, ok := sessionOptions[key]; ok {
				rawOpts[key] = value
			} else {
				// or from option defauls
				rawOpts[key] = opt.value
			}
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
	case "bool":
		value, err = strconv.ParseBool(arg)
	case "time":
		if arg == "" {
			value = time.Time{}
		} else {
			value, err = time.Parse("2006-01-02", arg)
		}
	default:
		// Cannot be reached
	}

	return value, err
}
