package display

import "github.com/nilsbu/lastfm/pkg/format"

type Display interface {
	Display(f format.Formatter) error
}
