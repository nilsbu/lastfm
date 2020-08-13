package charts2

import "fmt"

// Title identifies a line of charts.
// It is a Stringer, meaning String() can be used to print a meaningful name of
// the Title. It has a function Key() which provides a unique string which
// identifies the title uniquely in the charts. The methods Artist() and Song()
// optionally provide artist and song name.
type Title interface {
	fmt.Stringer
	Key() string
	Artist() string
	Song() string
}

// KeyTitle is a Title which has the same key as the printed string. Artist()
// and Song() are empty strings.
type KeyTitle string

// String returns the KeyTitle's string.
func (t KeyTitle) String() string {
	return string(t)
}

// Key returns the KeyTitle's string.
func (t KeyTitle) Key() string {
	return string(t)
}

// Artist returns ''.
func (t KeyTitle) Artist() string {
	return ""
}

// Song returns ''.
func (t KeyTitle) Song() string {
	return ""
}

// ArtistTitle is a Title which uses an artit's name as key and printed string.
// Song() is "".
type ArtistTitle string

// String returns the artist's name.
func (k ArtistTitle) String() string {
	return string(k)
}

// Artist returns the artist's name.
func (k ArtistTitle) Artist() string {
	return string(k)
}

// Key returns the artist's name.
func (k ArtistTitle) Key() string {
	return string(k)
}

// Song returns "".
func (k ArtistTitle) Song() string {
	return ""
}
