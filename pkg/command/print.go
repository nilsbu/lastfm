package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/config"
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
	"github.com/pkg/errors"
)

type printTotal struct {
	by   string
	name string
	n    int
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	sums := charts.Compile(plays).Sum()

	out, err := getOutCharts(cmd.by, cmd.name, sums, s)
	if err != nil {
		return err
	}

	f := &format.Charts{
		Charts:    out,
		Column:    -1,
		Count:     cmd.n,
		Numbered:  true,
		Precision: 0,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}

type printFade struct {
	by   string
	name string
	n    int
	hl   float64
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	fade := charts.Compile(plays).Fade(cmd.hl)

	out, err := getOutCharts(cmd.by, cmd.name, fade, s)
	if err != nil {
		return err
	}

	f := &format.Charts{
		Charts:    out,
		Column:    -1,
		Count:     cmd.n,
		Numbered:  true,
		Precision: 2,
	}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}

func getOutCharts(
	by, name string,
	cha charts.Charts,
	r rsrc.Reader) (charts.Charts, error) {
	if name == "" {
		switch by {
		case "all":
			return cha, nil
		case "super":
			tags, err := organize.LoadArtistTags(cha.Keys(), r)
			if err != nil {
				return nil, err
			}

			return cha.Supertags(tags, config.Supertags), nil
		default:
			return nil, fmt.Errorf("chart type '%v' not supported", by)
		}
	} else {
		var container map[string]charts.Charts
		switch by {
		case "all":
			return nil, errors.New("name must be empty for chart type 'all'")
		case "super":
			tags, err := organize.LoadArtistTags(cha.Keys(), r)
			if err != nil {
				return nil, err
			}

			container = cha.SplitBySupertag(tags, config.Supertags)
		default:
			return nil, fmt.Errorf("chart type '%v' not supported", by)
		}

		out, ok := container[name]
		if !ok {
			return nil, fmt.Errorf("name '%v' not found", name)
		}

		return out, nil
	}
}

type printTags struct {
	artist string
}

func (cmd printTags) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {

	tags, err := unpack.LoadArtistTags(cmd.artist, s)
	if err != nil {
		return err
	}

	col := make(charts.Column, len(tags))
	for i, tag := range tags {
		col[i] = charts.Score{Name: tag.Name, Score: float64(tag.Count)}
	}

	f := &format.Column{
		Column:   col,
		Numbered: true}

	err = d.Display(f)
	if err != nil {
		return err
	}

	return nil
}
