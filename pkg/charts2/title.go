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
