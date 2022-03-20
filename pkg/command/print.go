package command

import (
	"fmt"
	"time"

	async "github.com/nilsbu/async"
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
	Accumulate(c charts.LazyCharts) charts.LazyCharts
	PrintCharts() printCharts
}

func (cmd printCharts) PrintCharts() printCharts {
	return cmd
}

func getOutCharts(
	session *unpack.SessionInfo,
	pcd printChartsDescriptor,
	r rsrc.Reader,
) (charts.LazyCharts, error) {
	cmd := pcd.PrintCharts()

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(r))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load user info")
	}

	bookmark, err := unpack.LoadBookmark(session.User, r)
	if err != nil {
		return nil, err
	}

	corrections, err := unpack.LoadArtistCorrections(session.User, r)
	if err != nil {
		return nil, err
	}

	days := int((bookmark.Midnight() - user.Registered.Midnight()) / 86400)
	plays := make([][]charts.Song, days+1)
	err = async.Pie(days+1, func(i int) error {
		day := user.Registered.AddDate(0, 0, i)
		if songs, err := unpack.LoadDayHistory(session.User, day, r); err == nil {
			for j, song := range songs {
				if c, ok := corrections[song.Artist]; ok {
					songs[j].Artist = c
				}
			}
			plays[i] = songs
			return nil
		} else {
			return err
		}
	})
	if err != nil {
		return nil, err
	}

	var base charts.LazyCharts
	switch {
	case cmd.keys == "song" && cmd.duration:
		base = charts.SongsDuration(plays)
	case cmd.keys == "song" && !cmd.duration:
		base = charts.Songs(plays)
	case cmd.duration:
		base = charts.ArtistsDuration(plays)
	default:
		base = charts.Artists(plays)
	}

	gaussian := charts.Gaussian(base, 7, 2*7+1, true, false)
	normalized := charts.Normalize(gaussian)

	if cmd.normalized {
		base = normalized
	}

	accCharts := pcd.Accumulate(base)

	partition, err := cmd.getPartition(session, r, gaussian, accCharts, user.Registered)
	if err != nil {
		return nil, err
	}

	if cmd.name == "" {
		if partition == nil {
			return accCharts, nil
		}

		return charts.Group(accCharts, partition), nil
	}

	if partition == nil {
		return nil, errors.New("name must be empty for chart type 'all'")
	}

	found := false
	for _, partition := range partition.Partitions() {
		if partition.Key() == cmd.name {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("name '%v' is no partition", cmd.name)
	}

	return charts.Subset(accCharts, partition, charts.KeyTitle(cmd.name)), nil
}

func (cmd printCharts) getPartition(
	session *unpack.SessionInfo,
	r rsrc.Reader,
	gaussian, normalized charts.LazyCharts,
	registered rsrc.Day,
) (charts.Partition, error) {
	switch cmd.by {
	case "all":
		return nil, nil
	case "year":
		return charts.YearPartition(gaussian, normalized, registered), nil
	case "total":
		return charts.TotalPartition(normalized.Titles()), nil
	case "super":
		tags, err := loadArtistTags(normalized, r)
		if err != nil {
			return nil, err
		}

		corrections, _ := unpack.LoadSupertagCorrections(session.User, r)

		return charts.FirstTagPartition(tags, config.Supertags, corrections), nil
	case "country":
		tags, err := loadArtistTags(normalized, r)
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

type printTotal struct {
	printCharts
	date time.Time
}

func (cmd printTotal) Accumulate(c charts.LazyCharts) charts.LazyCharts {
	return charts.Sum(c)
}

func (cmd printTotal) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	cha, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	col := -1
	null := time.Time{}
	if cmd.date != null {
		col = int((rsrc.DayFromTime(cmd.date).Midnight() - user.Registered.Midnight()) / 86400)
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     cha,
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

func (cmd printFade) Accumulate(c charts.LazyCharts) charts.LazyCharts {
	return charts.Fade(c, cmd.hl)
}

func (cmd printFade) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	cha, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	col := -1
	null := time.Time{}
	if cmd.date != null {
		col = int((rsrc.DayFromTime(cmd.date).Midnight() - user.Registered.Midnight()) / 86400)
	}

	f := &format.Charts{
		Charts:     cha,
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

func (cmd printPeriod) Accumulate(c charts.LazyCharts) charts.LazyCharts {
	return c
}

func (cmd printPeriod) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	cha, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	rnge, err := charts.ParseRange(cmd.period, user.Registered, cha.Len())
	if err != nil {
		return errors.Wrap(err, "invalid range")
	}

	cha = charts.Sum(charts.Interval(cha, rnge))

	col := -1

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     cha,
		Column:     col,
		Count:      cmd.n,
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

type printInterval struct {
	printCharts
	begin  time.Time
	before time.Time
}

func (cmd printInterval) Accumulate(c charts.LazyCharts) charts.LazyCharts {
	return c
}

func (cmd printInterval) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	cha, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}

	rnge, err := charts.CroppedRange(
		rsrc.DayFromTime(cmd.begin),
		rsrc.DayFromTime(cmd.before),
		user.Registered, cha.Len())
	if err != nil {
		return errors.Wrap(err, "invalid range")
	}

	cha = charts.Sum(charts.Interval(cha, rnge))

	col := -1

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}
	f := &format.Charts{
		Charts:     cha,
		Column:     col,
		Count:      cmd.n,
		Numbered:   true,
		Precision:  prec,
		Percentage: cmd.percentage,
	}

	return d.Display(f)
}

// type printFadeMax struct {
// 	printCharts
// 	hl float64
// }

// func (cmd printFadeMax) Accumulate(c charts.LazyCharts) charts.LazyCharts {
// 	return charts.Fade(c, cmd.hl)
// }

// func (cmd printFadeMax) Execute(
// 	session *unpack.SessionInfo, s store.Store, d display.Display) error {
// 	cha, err := getOutCharts2(session, cmd, s)
// 	if err != nil {
// 		return err
// 	}

// 	max := charts.Max(cha)
// 	sumTotal := charts.Sum(max)
// 	col = col.Top(cmd.n)

// 	prec := 0
// 	if cmd.normalized {
// 		prec = 2
// 	}
// 	f := &format.Column{
// 		Column:     col,
// 		Numbered:   true,
// 		Percentage: cmd.percentage,
// 		Precision:  prec,
// 		SumTotal:   sumTotal,
// 	}

// 	return d.Display(f)
// }

type printTags struct {
	artist string
}

func (cmd printTags) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {

	tags, err := unpack.LoadArtistTags(cmd.artist, unpack.NewCacheless(s))
	if err != nil {
		return err
	}

	col := make(map[string][]float64, len(tags))
	for _, tag := range tags {
		col[tag.Name] = []float64{float64(tag.Count)}
	}

	f := &format.Column{
		Column:     charts.FromMap(col),
		Numbered:   true,
		Precision:  0,
		Percentage: false,
		SumTotal:   0,
	}

	return d.Display(f)
}
