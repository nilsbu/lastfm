package charts

import "fmt"

// Title identifies a line of charts.
// It is a Stringer, meaning String() can be used to print a meaningful name of
// the Title. It has a function Key() which provides a unique string which
// identifies the title uniquely in the charts. The methods Artist()
// optionally provides an artist name.
type Title interface {
	fmt.Stringer
	Key() string
	Artist() string
}

type keyTitle string

// KeyTitle is a Title which has the same key as the printed string. Artist()
// is an empty string.
func KeyTitle(s string) Title {
	return keyTitle(s)
}

// String returns the keyTitle's string.
func (t keyTitle) String() string {
	return string(t)
}

// Key returns the keyTitle's string.
func (t keyTitle) Key() string {
	return string(t)
}

// Artist returns ''.
func (t keyTitle) Artist() string {
	return ""
}

type artistTitle string

// ArtistTitle is a Title which uses an artist's name as key and printed string.
func ArtistTitle(s string) Title {
	return artistTitle(s)
}

// String returns the artist's name.
func (k artistTitle) String() string {
	return string(k)
}

// Artist returns the artist's name.
func (k artistTitle) Artist() string {
	return string(k)
}

// Key returns the artist's name.
func (k artistTitle) Key() string {
	return string(k)
}

type songTitle struct {
	artist, title string
}

// SongTitle returns a Title that prints the title in the format
// "<artist> - <song>" and has a unique key for each artist-song combination
// assuming the artist's name does not contain a line break.
func SongTitle(s Song) Title {
	return songTitle{s.Artist, s.Title}
}

func (t songTitle) String() string {
	return fmt.Sprintf("%v - %v", t.artist, t.title)
}

func (t songTitle) Artist() string {
	return t.artist
}

func (t songTitle) Key() string {
	return fmt.Sprintf("%v\n%v", t.artist, t.title)
}

type stringTitle string

// StringTitle is a Title which is non-empty only for String().
func StringTitle(s string) Title {
	return stringTitle(s)
}

// String returns the the string.
func (k stringTitle) String() string {
	return string(k)
}

// Artist returns ""
func (k stringTitle) Artist() string {
	return ""
}

// Key returns ""
func (k stringTitle) Key() string {
	return string(k)
}