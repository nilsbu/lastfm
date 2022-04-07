package pipeline

import (
	"fmt"
	"strconv"
	"strings"

	async "github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/config"
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/info"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type dynamic interface {
	Exec() (interface{}, error)
}

type once struct {
	f      func() (interface{}, error)
	ran    bool
	result interface{}
	err    error
}

func newDynamic(f func() (interface{}, error)) dynamic {
	return &once{f: f}
}

func (d *once) Exec() (interface{}, error) {
	if !d.ran {
		d.result, d.err = d.f()
		d.ran = true
	}
	return d.result, d.err
}

// TODO test Pipeline
// TODO cleanup Pipeline
type Pipeline interface {
	Execute(steps []string) (charts.Charts, error)
	Registered() rsrc.Day
	Session() *unpack.SessionInfo
}

type pipeline struct {
	graph     graph
	bookmarks map[string][]string
	vars      dynamic
	session   *unpack.SessionInfo
	store     io.Store
}

type vars struct {
	user        *unpack.User
	bookmark    rsrc.Day
	corrections map[string]string
	plays       [][]info.Song
}

func New(session *unpack.SessionInfo, s io.Store) Pipeline {
	pl := &pipeline{
		graph:     *newGraph(10),
		bookmarks: map[string][]string{},
		session:   session,
		store:     s,
	}

	pl.vars = newDynamic(func() (interface{}, error) {
		return pl.load()
	})
	return pl
}

func (w *pipeline) Registered() rsrc.Day {
	v, err := w.vars.Exec()
	if err != nil {
		return nil
	} else {
		return v.(*vars).user.Registered
	}
}

func (w *pipeline) Session() *unpack.SessionInfo {
	return w.session
}

func (w *pipeline) Execute(steps []string) (charts.Charts, error) {
	if w.session.User == "" {
		return nil, fmt.Errorf("no user name given, session might not be properly initialized")
	}

	// Ensure that gaussian exists, might be needed for year partition
	w.bookmarks["gaussian"] = []string{steps[0], "gaussian", "cache"}
	_, err := w.runSteps(w.bookmarks["gaussian"])
	if err != nil {
		return nil, err
	}

	return w.runSteps(steps)
}

func (w *pipeline) runSteps(steps []string) (charts.Charts, error) {
	var parent charts.Charts
	var err error
	for i, step := range steps {
		p := w.graph.get(steps[:i+1])
		if p != nil {
			parent = p
		} else {
			if i == 0 {
				parent, err = w.root(steps[0])
			} else {
				parent, err = w.step(step, parent)
			}
			if err != nil {
				return nil, errors.Wrapf(err, "during step '%v'", step)
			}
			w.graph.set(steps[:i+1], parent)
		}
	}

	return parent, err
}

func (w *pipeline) load() (*vars, error) {
	v := &vars{}

	err := async.Pe([]func() error{
		func() error {
			var err error
			v.user, err = unpack.LoadUserInfo(w.session.User, unpack.NewCacheless(w.store))
			return errors.Wrap(err, "failed to load user info")
		},
		func() error {
			var err error
			v.corrections, err = unpack.LoadArtistCorrections(w.session.User, w.store)
			return err
		},
		func() error {
			var err error
			v.bookmark, err = unpack.LoadBookmark(w.session.User, w.store)
			return err
		},
	})
	if err != nil {
		return nil, err
	}

	days := rsrc.Between(v.user.Registered, v.bookmark).Days()
	v.plays = make([][]info.Song, days+1)
	err = async.Pie(days+1, func(i int) error {
		day := v.user.Registered.AddDate(0, 0, i)
		if songs, err := unpack.LoadDayHistory(v.user.Name, day, w.store); err == nil {
			for j, song := range songs {
				if c, ok := v.corrections[song.Artist]; ok {
					songs[j].Artist = c
				}
			}
			v.plays[i] = songs
			return nil
		} else {
			return err
		}
	})
	return v, err
}

func (w *pipeline) root(s string) (charts.Charts, error) {
	var c charts.Charts
	switch s {
	case "songsduration":
		c = charts.LoadSongsDuration(w.session.User, w.store)
	case "songs":
		c = charts.LoadSongs(w.session.User, w.store)
	case "artistsduration":
		c = charts.LoadArtistsDuration(w.session.User, w.store)
	default:
		c = charts.LoadArtists(w.session.User, w.store)
	}
	return w.graph.set([]string{s}, c), nil
}

func (w *pipeline) step(step string, parent charts.Charts) (charts.Charts, error) {
	// TODO w.vars isn't always needed here
	vv, err := w.vars.Exec()
	if err != nil {
		return nil, err
	}
	v := vv.(*vars)

	split := strings.Split(step, ",")
	switch split[0] {
	case "id":
		return charts.Id(parent), nil

	case "cache":
		return charts.Cache(parent), nil

	case "sum":
		return charts.Sum(parent), nil

	case "max":
		return charts.Max(parent), nil

	case "normalize":
		return charts.Normalize(parent), nil

	case "gaussian":
		return charts.Gaussian(parent, 7, 2*7+1, true, false), nil
	case "fade":
		hl, _ := strconv.ParseFloat(split[1], 64)
		return charts.Fade(parent, hl), nil

	case "multiply":
		s, _ := strconv.ParseFloat(split[1], 64)
		return charts.Multiply(parent, s), nil

	case "group":
		gaussian, _ := w.runSteps(w.bookmarks["gaussian"])
		partition, err := w.getPartition(split[1], gaussian, parent)
		if err != nil {
			return nil, err
		} else {
			return charts.Group(parent, partition), nil
		}

	case "split":
		gaussian, _ := w.runSteps(w.bookmarks["gaussian"])
		partition, err := w.getPartition(split[1], gaussian, parent)
		if err != nil {
			return nil, err
		} else {
			if !partitionContains(partition, split[2]) {
				return nil, fmt.Errorf("name '%v' is no partition", split[2])
			} else {
				return charts.Subset(parent, partition, charts.KeyTitle(split[2])), nil
			}
		}

	case "day":
		col := rsrc.Between(v.user.Registered, rsrc.ParseDay(split[1])).Days()
		return charts.Column(parent, col), nil

	case "period":
		rnge, err := charts.ParseRange(split[1], v.user.Registered, parent.Len())
		if err != nil {
			return nil, errors.Wrap(err, "invalid range")
		} else {
			return charts.Interval(parent, rnge), nil
		}

	case "periods":
		rnge, err := charts.ParseRanges(split[1], v.user.Registered, parent.Len())
		if err != nil {
			return nil, errors.Wrap(err, "invalid range")
		} else {
			return charts.Intervals(parent, rnge, charts.Sum), nil
		}

	case "step":
		rnge, err := charts.ParseRanges(split[1], v.user.Registered, parent.Len())
		if err != nil {
			return nil, errors.Wrap(err, "invalid range")
		} else {
			return charts.Intervals(parent, rnge, charts.Id), nil
		}

	case "interval":
		rnge, err := charts.CroppedRange(
			rsrc.ParseDay(split[1]),
			rsrc.ParseDay(split[2]),
			v.user.Registered, parent.Len())
		if err != nil {
			return nil, errors.Wrap(err, "invalid range")
		} else {
			return charts.Interval(parent, rnge), nil
		}

	case "top":
		n, _ := strconv.Atoi(split[1])
		titles, _ := charts.Top(parent, n)
		return charts.Only(parent, titles), nil

	case "column":
		i, _ := strconv.Atoi(split[1])
		return charts.Column(parent, i), nil
	default:
		return nil, errors.New("step does not exist")
	}
}

func partitionContains(partition charts.Partition, name string) bool {
	found := false

	partitions, _ := partition.Partitions()
	for _, partition := range partitions {
		if partition.Key() == name {
			found = true
			break
		}
	}
	return found
}

func (w *pipeline) getPartition(
	step string,
	gaussian, parent charts.Charts,
) (charts.Partition, error) {
	switch step {
	case "all":
		return nil, nil
	case "year":
		vv, err := w.vars.Exec()
		if err != nil {
			return nil, err
		}

		return charts.YearPartition(gaussian, parent, vv.(*vars).user.Registered)
	case "total":
		return charts.TotalPartition(parent.Titles()), nil
	case "super":
		corrections, _ := unpack.LoadSupertagCorrections(w.session.User, w.store)
		return charts.TagPartition(parent, config.Supertags, corrections, w.store), nil
	case "country":
		corrections, _ := unpack.LoadCountryCorrections(w.session.User, w.store)
		return charts.TagPartition(parent, config.Countries, corrections, w.store), nil
	case "tags":
		titles := parent.Titles()
		artists := make([]string, len(titles))
		for i := range titles {
			artists[i] = titles[i].Artist()
		}

		at, _ := organize.LoadArtistTags(artists, w.store)
		tags := make([][]info.Tag, len(titles))
		for i, title := range titles {
			tags[i] = at[title.String()]
		}

		return charts.TagWeightPartition(titles, tags, config.Blacklist()), nil
	case "groups":
		replacements, err := unpack.LoadGroups(w.session.User, w.store)
		if err != nil {
			return nil, err
		}

		return charts.PartialReplacements(parent.Titles(), replacements), nil
	default:
		return nil, fmt.Errorf("chart type '%v' not supported", step)
	}
}
