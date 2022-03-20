package charts

// Song contains basic information about a song.
type Song struct {
	Artist, Title, Album string
	Duration             float64
}

// Tag contains information about a tag.
type Tag struct {
	Name   string
	Total  int64
	Reach  int64
	Weight int
}
