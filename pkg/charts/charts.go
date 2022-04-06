package charts

import (
	"sort"

	"github.com/nilsbu/lastfm/pkg/info"
)

type Charts interface {
	Data(titles []Title, begin, end int) ([][]float64, error)

	Titles() []Title
	Len() int
}

type charts struct {
	songs   func() ([][]info.Song, error)
	key     func(info.Song) Title
	value   func(info.Song) float64
	jobChan chan compileJob
	titles  []Title
	values  map[string][]float64
}
type compileJob chan<- error

// Artists compiles LazyCharts in which all songs by an artist are grouped.
func Artists(songs [][]info.Song) Charts {
	return new(func() ([][]info.Song, error) { return songs, nil },
		func(s info.Song) Title { return ArtistTitle(s.Artist) },
		func(s info.Song) float64 { return 1.0 })
}

// ArtistsDuration compiles LazyCharts in which all songs by an artist are
// grouped. The songs are weighted by duration before they are summed up.
func ArtistsDuration(songs [][]info.Song) Charts {
	return new(func() ([][]info.Song, error) { return songs, nil },
		func(s info.Song) Title { return ArtistTitle(s.Artist) },
		fDuration)
}

func fDuration(s info.Song) float64 {
	if s.Duration == 0 {
		return 4
	} else {
		return s.Duration
	}
}

// Songs compiles LazyCharts in which all songs are listed separately.
func Songs(songs [][]info.Song) Charts {
	return new(func() ([][]info.Song, error) { return songs, nil },
		func(s info.Song) Title { return SongTitle(s) },
		func(s info.Song) float64 { return 1.0 })
}

// SongsDuration compiles LazyCharts in which all songs are listed separately.
// The songs are weighted by duration.
func SongsDuration(songs [][]info.Song) Charts {
	return new(func() ([][]info.Song, error) { return songs, nil },
		func(s info.Song) Title { return SongTitle(s) },
		fDuration)
}

func new(songs func() ([][]info.Song, error), key func(info.Song) Title, value func(info.Song) float64) *charts {
	job := make(chan compileJob)
	c := &charts{
		songs:   songs,
		key:     key,
		value:   value,
		jobChan: job,
	}
	go c.compileWorker()
	return c
}

func (c *charts) compileWorker() {
	for back := range c.jobChan {
		if c.songs != nil {
			back <- c.compile()
		} else {
			back <- nil
		}
	}
}

func (c *charts) compile() error {
	c.values = map[string][]float64{}
	c.titles = []Title{}

	songs, err := c.songs()
	if err != nil {
		return err
	}

	// TODO can this be parallelized?
	for d, day := range songs {
		for _, song := range day {
			k := c.key(song)
			if line, ok := c.values[k.Key()]; ok {
				line[d] += c.value(song)
			} else {
				c.titles = append(c.titles, k)
				c.values[k.Key()] = make([]float64, len(songs))
				c.values[k.Key()][d] = c.value(song)
			}
		}
	}

	c.songs = nil
	return nil
}

func (c *charts) await() error {
	back := make(chan error)
	c.jobChan <- back
	return <-back
}

func FromMap(data map[string][]float64) Charts {
	titles := make([]Title, len(data))
	i := 0
	for t := range data {
		titles[i] = KeyTitle(t)
		i++
	}
	sort.Slice(titles, func(i, j int) bool { return titles[i].Key() < titles[j].Key() })

	charts := &charts{
		values: data,
		titles: titles,
	}

	return charts
}

// Pair contains a Title and a line of values.
type Pair struct {
	Title  Title
	Values []float64
}

// TODO doc & test
func InOrder(data []Pair) Charts {
	titles := make([]Title, len(data))
	values := make(map[string][]float64, len(data))
	for i, pair := range data {
		titles[i] = pair.Title
		values[pair.Title.Key()] = pair.Values
	}

	charts := &charts{
		titles: titles,
		values: values,
	}

	return charts
}

func (l *charts) Data(titles []Title, begin, end int) ([][]float64, error) {
	if l.songs != nil {
		if err := l.await(); err != nil {
			return nil, err
		}
	}

	data := make([][]float64, len(titles))
	for i, t := range titles {
		data[i] = l.values[t.Key()][begin:end]
	}
	return data, nil
}

func (l *charts) Titles() []Title {
	if l.songs != nil {
		if err := l.await(); err != nil {
			return nil // TODO error gets lost
		}
	}

	// assumption: noone touches the return value
	return l.titles
}

func (l *charts) Len() int {
	if l.songs != nil {
		l.await() // TOOD error gets lost
	}

	for _, line := range l.values {
		return len(line)
	}
	return -1
}

type chartsNode struct {
	parent Charts
}

func (l chartsNode) Titles() []Title {
	return l.parent.Titles()
}

func (l chartsNode) Len() int {
	return l.parent.Len()
}
