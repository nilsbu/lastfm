package command

// TODO Re-add timeline
import (
	"fmt"
	"sort"

	"github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type printTimeline struct {
	n int
}

func (cmd printTimeline) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {

	user, err := unpack.LoadUserInfo(session.User, unpack.NewCacheless(s))
	if err != nil {
		return err
	}

	fcmd := &printFade{
		printCharts: printCharts{
			by:         "all",
			keys:       "artist",
			duration:   true,
			normalized: true,
		},
		hl: 365,
	}

	steps, err := fcmd.getSteps()
	if err != nil {
		return err
	}

	steps = setStep(steps, fmt.Sprintf("fade,%v", fcmd.hl))

	cha, err := pl.Execute(steps)
	if err != nil {
		return err
	}

	titles := cha.Titles()
	data, err := cha.Data(titles, 0, cha.Len())
	if err != nil {
		return err
	}

	tops := make([][]charts.Title, cha.Len())
	err = async.Pie(len(tops), func(i int) error {
		var err error
		tops[i], err = nTop(titles, data, cmd.n, i)
		return err
	})
	if err != nil {
		return err
	}

	totals := map[string]*top{}

	tmps := make([]*top, cmd.n)
	for i := range tmps {
		tmps[i] = &top{title: tops[0][i]}

		err := d.Display(&format.Message{
			Msg: fmt.Sprintf("%v: '%v' starts in the top %v", user.Registered, tops[0][i], cmd.n),
		})
		if err != nil {
			return err
		}
	}
	for i, tp := range tops {
		tmpp := []*top{}

		for _, tmp := range tmps {
			if in(tmp.title, tp) {
				tmpp = append(tmpp, tmp)
			} else {
				if old, ok := totals[tmp.title.Key()]; ok {
					totals[tmp.title.Key()] = &top{title: tmp.title, value: old.value + i - tmp.value}
				} else {
					totals[tmp.title.Key()] = &top{title: tmp.title, value: i - tmp.value}
				}

				err := d.Display(&format.Message{
					Msg: fmt.Sprintf("%v: '%v' was %vd in the top %v (since %v)",
						user.Registered.AddDate(0, 0, i-1), tmp.title, i-tmp.value, cmd.n, user.Registered.AddDate(0, 0, tmp.value)),
				})
				if err != nil {
					return err
				}
			}
		}
		for _, t := range tp {
			found := false
			for _, tmp := range tmpp {
				if t.Key() == tmp.title.Key() {
					found = true
					break
				}
			}
			if !found {
				tmpp = append(tmpp, &top{title: t, value: i})

				err := d.Display(&format.Message{
					Msg: fmt.Sprintf("%v: '%v' enters the top %vd", user.Registered.AddDate(0, 0, i), t, cmd.n),
				})
				if err != nil {
					return err
				}
			}
		}
		tmps = tmpp
	}
	for _, tmp := range tmps {
		if old, ok := totals[tmp.title.Key()]; ok {
			totals[tmp.title.Key()] = &top{title: tmp.title, value: old.value + cha.Len() - tmp.value}
		} else {
			totals[tmp.title.Key()] = &top{title: tmp.title, value: cha.Len() - tmp.value}
		}

		err = d.Display(&format.Message{
			Msg: fmt.Sprintf("'%v' is %vd on top (since %v)",
				tmp.title, cha.Len()-tmp.value, user.Registered.AddDate(0, 0, tmp.value)),
		})
		if err != nil {
			return err
		}
	}

	totalOrdered := make([]*top, len(totals))
	i := 0
	for _, v := range totals {
		totalOrdered[i] = v
		i++
	}

	sort.Slice(totalOrdered, func(i, j int) bool { return totalOrdered[i].value > totalOrdered[j].value })
	for _, top := range totalOrdered {
		err = d.Display(&format.Message{
			Msg: fmt.Sprintf("'%v' was in the top %v for a total of %vd", top.title, cmd.n, top.value),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type top struct {
	title charts.Title
	value int
}

func in(v charts.Title, vs []charts.Title) bool {
	for _, w := range vs {
		if v.Key() == w.Key() {
			return true
		}
	}
	return false
}

func nTop(fullTitles []charts.Title, data [][]float64, n, c int) ([]charts.Title, error) {
	m := n + 1
	if len(data) < n {
		m = len(data)
	}

	vs := make([]float64, m)
	ts := make([]charts.Title, m)
	i := 0
	for k, tv := range data {
		if tv[c] == 0 {
			continue
		}
		vs[i] = tv[c]
		ts[i] = fullTitles[k]
		for j := i; j > 0; j-- {
			if vs[j-1] < vs[j] {
				vs[j-1], vs[j] = vs[j], vs[j-1]
				ts[j-1], ts[j] = ts[j], ts[j-1]
			} else {
				break
			}
		}
		if i+1 < m {
			i++
		}
	}
	if len(ts) > n {
		ts = ts[:n]
	}
	if len(ts) > i && ts[i] == nil {
		ts = ts[:i]
	}

	return ts, nil
}
