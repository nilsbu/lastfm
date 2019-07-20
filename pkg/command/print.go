package command

import (
	"fmt"
	"time"

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

type printCharts struct {
	keys       string //defaults to "artist"
	by         string
	name       string
	percentage bool
	normalized bool
	entry      float64
	n          int
}

// TODO document
type printChartsDescriptor interface {
	Accumulate(c charts.Charts) charts.Charts
	PrintCharts() printCharts
}

func (cmd printCharts) PrintCharts() printCharts {
	return cmd
}

func (cmd printCharts) getPartition(
	session *unpack.SessionInfo,
	r rsrc.Reader,
	cha charts.Charts,
) (year charts.Partition, err error) {
	switch cmd.by {
	case "all":
		return
	case "year":
		entry := cmd.entry
		if entry == 0 {
			entry = 2
		}
		year = cha.GetYearPartition(entry)
		return
	case "super":
		keys := []string{}
		for _, key := range cha.Keys {
			keys = append(keys, key.ArtistName())
		}
		tags, err := organize.LoadArtistTags(keys, r)
		if err != nil {
			for _, e := range err.(*organize.MultiError).Errs {
				switch e.(type) {
				case *unpack.LastfmError:
					// TODO can this be tested?
					if e.(*unpack.LastfmError).IsFatal() {
						return nil, err
					}
				default:
					return nil, err
				}
			}
		}

		corrections, _ := unpack.LoadSupertagCorrections(session.User, r)

		return charts.FirstTagPartition(tags, config.Supertags, corrections), nil
	case "country":
		keys := []string{}
		for _, key := range cha.Keys {
			keys = append(keys, key.ArtistName())
		}
		tags, err := organize.LoadArtistTags(keys, r)
		if err != nil {
			for _, e := range err.(*organize.MultiError).Errs {
				switch e.(type) {
				case *unpack.LastfmError:
					// TODO can this be tested?
					if e.(*unpack.LastfmError).IsFatal() {
						return nil, err
					}
				default:
					return nil, err
				}
			}
		}

		corrections := map[string]string{}

		return charts.FirstTagPartition(tags, config.Countries, corrections), nil
	default:
		return nil, fmt.Errorf("chart type '%v' not supported", cmd.by)
	}
}

func getOutCharts(
	session *unpack.SessionInfo,
	pcd printChartsDescriptor,
	r rsrc.Reader,
) (charts.Charts, error) {
	cmd := pcd.PrintCharts()

	user, err := unpack.LoadUserInfo(session.User, r)
	if err != nil {
		return charts.Charts{}, errors.Wrap(err, "failed to load user info")
	}

	plays, err := unpack.LoadSongHistory(session.User, r)
	if err != nil {
		return charts.Charts{}, err
	}

	var cha charts.Charts
	if cmd.keys == "" || cmd.keys == "artist" {
		cha = charts.ArtistsFromSongs(plays, user.Registered)
	} else if cmd.keys == "song" {
		cha = charts.CompileSongs(plays, user.Registered)
	}

	replace, err := unpack.LoadArtistCorrections(session.User, r)
	if err == nil {
		// TODO correct does not work for songs
		cha = cha.Correct(replace)
	}

	partition, err := cmd.getPartition(session, r, cha)
	if err != nil {
		return charts.Charts{}, err
	}

	if cmd.normalized {
		nm := charts.GaussianNormalizer{
			Sigma:      30,
			MirrorBack: false}
		cha = nm.Normalize(cha)
	}

	accCharts := pcd.Accumulate(cha)

	if cmd.name == "" {
		if partition == nil {
			return accCharts, nil
		}

		return accCharts.Group(partition), nil
	}

	if partition == nil {
		return charts.Charts{}, errors.New("name must be empty for chart type 'all'")
	}

	out, ok := accCharts.Split(partition)[cmd.name]
	if !ok {
		return charts.Charts{}, fmt.Errorf("name '%v' not found", cmd.name)
	}

	return out, nil
}

type printTotal struct {
	printCharts
	date time.Time
}

func (cmd printTotal) Accumulate(c charts.Charts) charts.Charts {
	return c.Sum()
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	col := -1
	var null time.Time
	null = time.Time{}
	if cmd.date != null {
		col = out.Headers.Index(rsrc.DayFromTime(cmd.date))
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     out,
		Column:     col,
		Count:      cmd.n,
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printFade struct {
	printCharts
	hl   float64
	date time.Time
}

func (cmd printFade) Accumulate(c charts.Charts) charts.Charts {
	return c.Fade(cmd.hl)
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	col := -1
	var null time.Time
	null = time.Time{}
	if cmd.date != null {
		col = out.Headers.Index(rsrc.DayFromTime(cmd.date))
	}

	f := &format.Charts{
		Charts:     out,
		Column:     col,
		Count:      cmd.n,
		Numbered:   true,
		Precision:  2,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printPeriod struct {
	printCharts
	period string
}

func (cmd printPeriod) Accumulate(c charts.Charts) charts.Charts {
	return c.Sum()
}

func (cmd printPeriod) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	period, err := charts.Period(cmd.period)
	if err != nil {
		return err
	}

	col := out.Interval(period)
	sumTotal := col.Sum()
	col = col.Top(cmd.n)

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Column{
		Column:     col,
		Numbered:   true,
		Percentage: cmd.percentage,
		Precision:  prec,
		SumTotal:   sumTotal,
	}

	return d.Display(f)
}

type printInterval struct {
	printCharts
	begin  time.Time
	before time.Time
}

func (cmd printInterval) Accumulate(c charts.Charts) charts.Charts {
	return c.Sum()
}

func (cmd printInterval) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	interval := charts.Interval{
		Begin:  rsrc.DayFromTime(cmd.begin),
		Before: rsrc.DayFromTime(cmd.before),
	}

	col := out.Interval(interval)
	sumTotal := col.Sum()
	col = col.Top(cmd.n)

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Column{
		Column:     col,
		Numbered:   true,
		Percentage: cmd.percentage,
		Precision:  prec,
		SumTotal:   sumTotal,
	}

	return d.Display(f)
}

type printFadeMax struct {
	printCharts
	hl float64
}

func (cmd printFadeMax) Accumulate(c charts.Charts) charts.Charts {
	return c.Fade(cmd.hl)
}

func (cmd printFadeMax) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	col := out.Max()
	sumTotal := col.Sum()
	col = col.Top(cmd.n)

	prec := 0
	if cmd.normalized {
		prec = 2
	}
	f := &format.Column{
		Column:     col,
		Numbered:   true,
		Percentage: cmd.percentage,
		Precision:  prec,
		SumTotal:   sumTotal,
	}

	return d.Display(f)
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

	return d.Display(f)
}
