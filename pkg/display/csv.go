package display

import (
	"os"

	"github.com/nilsbu/lastfm/pkg/format"
)

type CSV struct {
	path    string
	decimal string
}

func NewCSV(path, decimal string) *CSV {
	return &CSV{path: path, decimal: decimal}
}

func (d *CSV) Display(f format.Formatter) error {
	file, err := os.Create(d.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return f.CSV(file, d.decimal)
}
