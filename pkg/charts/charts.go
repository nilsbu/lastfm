package charts

import "sort"

type Charts interface {
	// TODO can Charts.Column() be removed?
	Column(titles []Title, index int) []float64
	Data(titles []Title, begin, end int) [][]float64

	Titles() []Title
	Len() int
}

type charts struct {
	songs   [][]Song
	key     func(Song) Title
	value   func(Song) float64
	jobChan chan compileJob
	titles  []Title
	values  map[string][]float64
}
type compileJob chan<- interface{}

// Artists compiles LazyCharts in which all songs by an artist are grouped.
func Artists(songs [][]Song) Charts {
	return new(songs,
		func(s Song) Title { return ArtistTitle(s.Artist) },
		func(s Song) float64 { return 1.0 })
}

// ArtistsDuration compiles LazyCharts in which all songs by an artist are
// grouped. The songs are weighted by duration before they are summed up.
func ArtistsDuration(songs [][]Song) Charts {
	return new(songs,
		func(s Song) Title { return ArtistTitle(s.Artist) },
		fDuration)
}

func fDuration(s Song) float64 {
	if s.Duration == 0 {
		return 4
	} else {
		return s.Duration
	}
}

// Songs compiles LazyCharts in which all songs are listed separately.
func Songs(songs [][]Song) Charts {
	return new(songs,
		func(s Song) Title { return SongTitle(s) },
		func(s Song) float64 { return 1.0 })
}

// SongsDuration compiles LazyCharts in which all songs are listed separately.
// The songs are weighted by duration.
func SongsDuration(songs [][]Song) Charts {
	return new(songs,
		func(s Song) Title { return SongTitle(s) },
		fDuration)
}

func new(songs [][]Song, key func(Song) Title, value func(Song) float64) *charts {
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
	for {
		select {
		case back := <-c.jobChan:
			if c.songs != nil {
				c.compile()
			}
			back <- nil
		}
	}
}

func (c *charts) compile() {
	c.values = map[string][]float64{}
	c.titles = []Title{}

	// TODO can this be parallelized?
	for d, day := range c.songs {
		for _, song := range day {
			k := c.key(song)
			if line, ok := c.values[k.Key()]; ok {
				line[d] += c.value(song)
			} else {
				c.titles = append(c.titles, k)
				c.values[k.Key()] = make([]float64, len(c.songs))
				c.values[k.Key()][d] = c.value(song)
			}
		}
	}

	c.songs = nil
}

func (c *charts) await() {
	back := make(chan interface{})
	c.jobChan <- back
	<-back
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

func (l *charts) Column(titles []Title, index int) []float64 {
	if l.songs != nil {
		l.await()
	}

	col := make([]float64, len(titles))
	for i, t := range titles {
		col[i] = l.values[t.Key()][index]
	}
	return col
}

func (l *charts) Data(titles []Title, begin, end int) [][]float64 {
	if l.songs != nil {
		l.await()
	}

	data := make([][]float64, len(titles))
	for i, t := range titles {
		data[i] = l.values[t.Key()][begin:end]
	}
	return data
}

func (l *charts) Titles() []Title {
	if l.songs != nil {
		l.await()
	}

	// assumption: noone touches the return value
	return l.titles
}

func (l *charts) Len() int {
	if l.songs != nil {
		l.await()
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
