package display

import (
	"os"

	"github.com/nilsbu/lastfm/pkg/format"
)

type File struct {
	Path string
}

func NewFile(path string) *File {
	return &File{Path: path}
}

func (d *File) Display(f format.Formatter) error {
	file, err := os.Create(d.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	return f.Plain(file)
}
