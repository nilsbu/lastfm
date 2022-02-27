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
	keys       string // defaults to "artist"
	by         string
	name       string
	percentage bool
	normalized bool
	duration   bool
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
) (charts.Partition, error) {
	switch cmd.by {
	case "all":
		return nil, nil
	case "year":
		entry := cmd.entry
		if entry == 0 {
			entry = 2
		}
		return cha.GetYearPartition(entry), nil
	case "total":
		return charts.TotalPartition{}, nil
	case "super":
		tags, err := loadArtistTags(cha, r)
		if err != nil {
			return nil, err
		}

		corrections, _ := unpack.LoadSupertagCorrections(session.User, r)

		return charts.FirstTagPartition(tags, config.Supertags, corrections), nil
	case "country":
		tags, err := loadArtistTags(cha, r)
		if err != nil {
			return nil, err
		}

		corrections, _ := unpack.LoadCountryCorrections(session.User, r)

		return charts.FirstTagPartition(tags, config.Countries, corrections), nil
	default:
		return nil, fmt.Errorf("chart type '%v' not supported", cmd.by)
	}
}

func loadArtistTags(
	cha charts.Charts,
	r rsrc.Reader,
) (map[string][]charts.Tag, error) {
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

	return tags, nil
}

func getOutCharts(
	session *unpack.SessionInfo,
	pcd printChartsDescriptor,
	r rsrc.Reader,
) (charts.Charts, error) {
	cmd := pcd.PrintCharts()

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(r))
	if err != nil {
		return charts.Charts{}, errors.Wrap(err, "failed to load user info")
	}

	bookmark, err := unpack.LoadBookmark(session.User, r)
	if err != nil {
		return charts.Charts{}, err
	}

	days := int((bookmark.Midnight() - user.Registered.Midnight()) / 86400)
	plays := make([][]charts.Song, days+1)
	for i := 0; i < days+1; i++ {
		day := user.Registered.AddDate(0, 0, i)
		if songs, err := unpack.LoadDayHistory(session.User, day, r); err == nil {
			plays[i] = songs
		} else {
			return charts.Charts{}, err
		}
	}
	if err != nil {
		return charts.Charts{}, err
	}

	var cha charts.Charts
	if cmd.keys == "song" || cmd.duration {
		cha = charts.CompileSongs(plays, user.Registered)
		if cmd.duration {
			cha, err = normalizeDuration(cha, r)
			if err != nil {
				return charts.Charts{}, err
			}
		}
		if cmd.keys != "song" {
			p := charts.NewArtistNamePartition(cha)
			cha = cha.Group(p)
		}
	} else if cmd.keys == "" || cmd.keys == "artist" {
		cha = charts.ArtistsFromSongs(plays, user.Registered)
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
			Sigma:      7,
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

func normalizeDuration(cha charts.Charts, r rsrc.Reader) (charts.Charts, error) {
	sd := charts.SongDurations{}
	for _, key := range cha.Keys {
		info := key.(charts.Song)
		if info.Duration > 0 {
			if _, ok := sd[info.Artist]; !ok {
				sd[info.Artist] = make(map[string]float64)
			}
			sd[info.Artist][info.Title] = info.Duration
		}
	}
	sd[""] = make(map[string]float64)
	sd[""][""] = 4.0
	return sd.Normalize(cha), nil
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

type printDay struct {
	printTotal
	// date time.Time
}

func (cmd printDay) Accumulate(c charts.Charts) charts.Charts {
	return c
}

func (cmd printDay) Execute(
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

type printTags struct {
	artist string
}

func (cmd printTags) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {

	tags, err := unpack.LoadArtistTags(cmd.artist, unpack.NewCacheless(s))
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
