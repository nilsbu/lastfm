package command

import (
	"fmt"
	"strconv"
	"strings"

	async "github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/config"
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

// TODO test and move
type Web interface {
	Execute(steps []string) (charts.LazyCharts, error)
	Registered() rsrc.Day
}

type web struct {
	charts   map[string]charts.LazyCharts
	baseType string
	vars     vars
	session  *unpack.SessionInfo
	store    store.Store
}

func newWeb(session *unpack.SessionInfo, s store.Store) Web {
	return &web{
		session: session,
		store:   s,
	}
}

func (w *web) Registered() rsrc.Day {
	return w.vars.user.Registered
}

func (w *web) Execute(steps []string) (charts.LazyCharts, error) {
	if w.baseType == "" {
		if err := w.load(); err != nil {
			return nil, err
		}
	}
	if w.baseType != steps[0] {
		w.calcDaily(steps[0])
	}

	var err error
	parent := w.charts["daily"]
	for _, step := range steps[1:] {
		parent, err = w.step(step, parent)
		if err != nil {
			return nil, errors.Wrapf(err, "during step '%v'", step)
		}
	}

	return parent, nil
}

func (w *web) load() error {
	err := async.Pe([]func() error{
		func() error {
			var err error
			w.vars.user, err = unpack.LoadUserInfo(w.session.User, unpack.NewCacheless(w.store))
			return errors.Wrap(err, "failed to load user info")
		},
		func() error {
			var err error
			w.vars.corrections, err = unpack.LoadArtistCorrections(w.session.User, w.store)
			return err
		},
		func() error {
			var err error
			w.vars.bookmark, err = unpack.LoadBookmark(w.session.User, w.store)
			return err
		},
	})
	if err != nil {
		return err
	}

	days := int((w.vars.bookmark.Midnight() - w.vars.user.Registered.Midnight()) / 86400)
	w.vars.plays = make([][]charts.Song, days+1)
	return async.Pie(days+1, func(i int) error {
		day := w.vars.user.Registered.AddDate(0, 0, i)
		if songs, err := unpack.LoadDayHistory(w.vars.user.Name, day, w.store); err == nil {
			for j, song := range songs {
				if c, ok := w.vars.corrections[song.Artist]; ok {
					songs[j].Artist = c
				}
			}
			w.vars.plays[i] = songs
			return nil
		} else {
			return err
		}
	})
}

func (w *web) calcDaily(s string) {
	w.charts = map[string]charts.LazyCharts{}
	switch {
	case strings.Contains(s, "songs duration"):
		w.charts["base"] = charts.SongsDuration(w.vars.plays)
	case strings.Contains(s, "song"):
		w.charts["base"] = charts.Songs(w.vars.plays)
	case strings.Contains(s, "artist duration"):
		w.charts["base"] = charts.ArtistsDuration(w.vars.plays)
	default:
		w.charts["base"] = charts.Artists(w.vars.plays)
	}

	w.charts["gaussian"] = charts.Cache(charts.Gaussian(w.charts["base"], 7, 2*7+1, true, false))
	w.charts["normalized"] = charts.Normalize(w.charts["gaussian"])

	if strings.Contains(s, "normalized") {
		w.charts["daily"] = charts.Id(w.charts["normalized"])
	} else {
		w.charts["daily"] = charts.Id(w.charts["base"])
	}
}

func (w *web) step(step string, parent charts.LazyCharts) (charts.LazyCharts, error) {
	split := strings.Split(step, " ")
	switch split[0] {
	case "id":
		return charts.Id(parent), nil

	case "sum":
		return charts.Sum(parent), nil

	case "max":
		return charts.Max(parent), nil

	case "normalize":
		return charts.Normalize(parent), nil

	case "fade":
		hl, _ := strconv.ParseFloat(split[1], 64)
		return charts.Fade(parent, hl), nil

	case "multiply":
		s, _ := strconv.ParseFloat(split[1], 64)
		return charts.Multiply(parent, s), nil

	case "group":
		partition, err := w.getPartition(split[1], w.charts["gaussian"], parent)
		if err != nil {
			return nil, err
		} else {
			return charts.Group(parent, partition), nil
		}

	case "split":
		partition, err := w.getPartition(split[1], w.charts["gaussian"], parent)
		if err != nil {
			return nil, err
		} else {
			if !partitionCongains(partition, split[2]) {
				return nil, fmt.Errorf("name '%v' is no partition", split[2])
			} else {
				return charts.Subset(parent, partition, charts.KeyTitle(split[2])), nil
			}
		}

	case "day":
		col := int((rsrc.ParseDay(split[1]).Midnight() - w.vars.user.Registered.Midnight()) / 86400)
		return charts.Column(parent, col), nil

	case "period":
		rnge, err := charts.ParseRange(split[1], w.vars.user.Registered, parent.Len())
		if err != nil {
			return nil, errors.Wrap(err, "invalid range")
		} else {
			return charts.Interval(parent, rnge), nil
		}

	case "periods":
		rnge, err := charts.ParseRanges(split[1], w.vars.user.Registered, parent.Len())
		if err != nil {
			return nil, errors.Wrap(err, "invalid range")
		} else {
			return charts.Intervals(parent, rnge, charts.Id), nil
		}

	case "interval":
		rnge, err := charts.CroppedRange(
			rsrc.ParseDay(split[1]),
			rsrc.ParseDay(split[2]),
			w.vars.user.Registered, parent.Len())
		if err != nil {
			return nil, errors.Wrap(err, "invalid range")
		} else {
			return charts.Interval(parent, rnge), nil
		}

	default:
		return nil, errors.New("step does not exist")
	}
}

func partitionCongains(partition charts.Partition, name string) bool {
	found := false
	for _, partition := range partition.Partitions() {
		if partition.Key() == name {
			found = true
			break
		}
	}
	return found
}

func (w *web) getPartition(
	step string,
	gaussian, parent charts.LazyCharts,
) (charts.Partition, error) {
	switch step {
	case "all":
		return nil, nil
	case "year":
		return charts.YearPartition(gaussian, parent, w.vars.user.Registered), nil
	case "total":
		return charts.TotalPartition(parent.Titles()), nil
	case "super":
		tags, err := loadArtistTags(parent, w.store)
		if err != nil {
			return nil, err
		}

		corrections, _ := unpack.LoadSupertagCorrections(w.session.User, w.store)

		return charts.FirstTagPartition(tags, config.Supertags, corrections), nil
	case "country":
		tags, err := loadArtistTags(parent, w.store)
		if err != nil {
			return nil, err
		}

		corrections, _ := unpack.LoadCountryCorrections(w.session.User, w.store)

		return charts.FirstTagPartition(tags, config.Countries, corrections), nil
	default:
		return nil, fmt.Errorf("chart type '%v' not supported", step)
	}
}

func loadArtistTags(
	cha charts.LazyCharts,
	r rsrc.Reader,
) (map[string][]charts.Tag, error) {
	keys := []string{}

	for _, key := range cha.Titles() {
		keys = append(keys, key.Artist())
	}

	tags, err := organize.LoadArtistTags(keys, r)
	if err != nil {
		for _, e := range err.(*async.MultiError).Errs {
			switch e := e.(type) {
			case *unpack.LastfmError:
				// TODO can this be tested?
				if e.IsFatal() {
					return nil, err
				}
			default:
				return nil, err
			}
		}
	}

	return tags, nil
}
