package command

import (
	"fmt"
	"sort"
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

		return charts.FirstTagPartition(tags, config.Countries, nil), nil
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

type printLifeExpectancy struct {
	printCharts
}

func (cmd printLifeExpectancy) Accumulate(c charts.Charts) charts.Charts {
	return c
}

func (cmd printLifeExpectancy) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	out, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	entryDates := out.FindEntryDatesDynamic(2.0)
	entry := make(map[string]rsrc.Day)
	for _, ed := range entryDates {
		entry[ed.Name] = ed.Date
	}

	var col charts.Column
	for i, key := range out.Keys {
		expectance := 0.0
		sum := 0.0
		for j, val := range out.Values[i] {
			expectance += float64(j) * val
			sum += val
		}
		if sum > 0 {
			expectance /= sum
		}
		if e, ok := entry[key.String()]; ok {
			expectance -= float64(out.Headers.Index(e))
			expectance /= 365.25
			col = append(col, charts.Score{Name: key.String(), Score: expectance})
		}
	}

	sort.Sort(col)

	f := &format.Column{
		Column:    col,
		Numbered:  true,
		Precision: 2,
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
	hl  float64
	min float64
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

	var col charts.Column
	var sumTotal float64
	if cmd.percentage {
		totalCmd := printTotal{printCharts: cmd.printCharts}
		tot, err := getOutCharts(session, totalCmd, s)
		if err != nil {
			return err
		}

		outMax := out.Max()
		total, err := tot.Column(-1)
		if err != nil {
			return err
		}

		var fademax charts.Column
		for i := 0; i < len(outMax); i++ {
			for _, score := range total {
				if outMax[i].Name == score.Name {
					if score.Score < cmd.min {
						break
					}
					fademax = append(fademax, charts.Score{
						Name:  score.Name,
						Score: outMax[i].Score / score.Score,
					})
					break
				}
			}
		}

		sort.Sort(fademax)

		col = fademax.Top(cmd.n)
		sumTotal = 1

	} else {
		col = out.Max()
		col = col.Top(cmd.n)

		if len(col) == 0 {
			sumTotal = 1
		} else {
			sumTotal = col[0].Score
		}
	}

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
		prec = 3
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

type printAfterEntry struct {
	printCharts
	days int
}

func (cmd printAfterEntry) Accumulate(c charts.Charts) charts.Charts {
	return c
}

func (cmd printAfterEntry) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {

	cha, err := getOutCharts(session, cmd, s)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}

	outCol := charts.Column{}

	if cmd.entry == 0 {
		cmd.entry = 2
	}

	cmd2 := cmd
	cmd2.normalized = false
	cmd2.by = "all"
	cmd2.name = ""
	cha2, _ := getOutCharts(session, cmd2, s)

	ed := cha2.FindEntryDatesDynamic(cmd.entry)

	out := cha.Sum()

	for _, date := range ed {
		d := out.Headers.Index(date.Date) + cmd.days

		if d >= out.Headers.Len() {
			continue
		} else {
			var v float64
			var found bool
			for i, k := range out.Keys {
				if date.Name == k.ArtistName() {
					v = out.Values[i][d]
					found = true
					break
				}
			}

			if found {
				outCol = append(outCol, charts.Score{
					Name:  date.Name,
					Score: v,
				})
			}
		}
	}

	sort.Sort(outCol)

	f := &format.Column{
		Column:    outCol,
		Numbered:  true,
		Precision: prec,
	}

	return d.Display(f)
}

//////

func getOutCharts2(
	session *unpack.SessionInfo,
	pcd printChartsDescriptor,
	r rsrc.Reader,
	drop int,
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

	plays = plays[:len(plays)-drop]

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

type printLastDays struct {
	printCharts
	days int
}

func (cmd printLastDays) Accumulate(c charts.Charts) charts.Charts {
	return c.Sum()
}

func (cmd printLastDays) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {

	cha, err := getOutCharts2(session, cmd, s, 0)
	if err != nil {
		return err
	}

	cha2, err := getOutCharts2(session, cmd, s, cmd.days)
	if err != nil {
		return err
	}

	prec := 0
	if cmd.percentage || cmd.normalized {
		prec = 2
	}

	outCol := charts.Column{}

	for i, key := range cha.Keys {
		fJ := 0
		for j, key2 := range cha2.Keys {
			if key.String() == key2.String() {
				fJ = j
				break
			}
		}

		d := cha.Values[i][len(cha.Values[i])-1] - cha2.Values[fJ][len(cha2.Values[fJ])-1]
		outCol = append(outCol, charts.Score{
			Name:  key.ArtistName(),
			Score: d,
		})
	}

	sort.Sort(outCol)
	outCol = outCol.Top(cmd.n)

	f := &format.Column{
		Column:    outCol,
		Numbered:  true,
		Precision: prec,
	}

	return d.Display(f)
}

///////////////////

type simpleKey string

func (s simpleKey) String() string {
	return string(s)
}

func (s simpleKey) ArtistName() string {
	return string(s)
}

func (s simpleKey) FullTitle() string {
	return string(s)
}

func getOutCharts3(
	session *unpack.SessionInfo,
	pcd printChartsDescriptor,
	r rsrc.Reader,
) (map[string]charts.Charts, error) {
	cmd := pcd.PrintCharts()

	user, err := unpack.LoadUserInfo(session.User, r)
	if err != nil {
		return map[string]charts.Charts{}, errors.Wrap(err, "failed to load user info")
	}

	plays, err := unpack.LoadSongHistory(session.User, r)
	if err != nil {
		return map[string]charts.Charts{}, err
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
		return map[string]charts.Charts{}, err
	}

	if cmd.normalized {
		nm := charts.GaussianNormalizer{
			Sigma:      7,
			MirrorBack: false}
		cha = nm.Normalize(cha)
	}

	accCharts := pcd.Accumulate(cha)

	if partition == nil {
		return map[string]charts.Charts{}, errors.New("name must be empty for chart type 'all'")
	}

	return accCharts.Split(partition), nil
}

type printTotalX struct {
	printCharts
	date time.Time
}

func (cmd printTotalX) Accumulate(c charts.Charts) charts.Charts {
	return c.Sum()
}

func (cmd printTotalX) Execute(
	session *unpack.SessionInfo, s store.Store, d display.Display) error {
	cha1, err := getOutCharts3(session, cmd, s)
	if err != nil {
		return err
	}

	corrections, _ := unpack.LoadSupertagCorrections(session.User, s)

	days := []map[string]float64{}
	for _, cha := range cha1 {
		for i := 0; i < cha.Len(); i++ {
			days = append(days, map[string]float64{})
		}
		break
	}

	for key, cha := range cha1 {
		tags, err := loadArtistTags(cha, s)
		if err != nil {
			return err
		}

		cc := cha.Group(charts.FirstTagPartition(tags, config.Supertags, corrections))

		for i, key2 := range cc.Keys {
			v2 := cc.Values[i]
			k := fmt.Sprintf("%v - %v", key, key2)
			for j, v := range v2 {
				if _, ok := days[j][k]; ok {
					days[j][k] = days[j][k] + v
				} else {
					days[j][k] = v
				}
			}
		}
	}

	user, err := unpack.LoadUserInfo(session.User, s)
	if err != nil {
		return errors.Wrap(err, "failed to load user info")
	}
	out := charts.CompileArtists(days, user.Registered)

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

////

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
