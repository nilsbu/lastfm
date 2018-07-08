package command

import (
	"github.com/nilsbu/lastfm/config"
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type printTotal struct {
	n int
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	sums := charts.Compile(plays).Sum()
	f := &format.Charts{
		Charts:    sums,
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
	n  int
	hl float64
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	sums := charts.Compile(plays).Fade(cmd.hl)
	f := &format.Charts{
		Charts:    sums,
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

type printTotalSuper struct {
	n int
}

func (cmd printTotalSuper) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	plays, err := unpack.LoadAllDayPlays(session.User, s)
	if err != nil {
		return err
	}

	sums := charts.Compile(plays).Sum()

	tags, err := organize.LoadArtistTags(sums.Keys(), s)
	if err != nil {
		return err
	}

	tsums := sums.Supertags(tags, config.Supertags)

	f := &format.Charts{
		Charts:    tsums,
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
