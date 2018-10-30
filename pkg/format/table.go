package format

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Table struct {
	Charts charts.Charts
	First  rsrc.Day
	Step   int
	Count  int
}

func (f *Table) CSV(w io.Writer, decimal string) error {
	if f.Charts.Len() == 0 {
		return nil
	}

	io.WriteString(w, "\"name\";")

	date, ok := f.First.Midnight()
	if !ok {
		return errors.New("'First' date is invalid")
	}

	for i := 0; i < f.Charts.Len(); i += f.Step {
		if i > 0 {
			io.WriteString(w, ";")
		}

		t := time.Unix(date+int64(i*86400), 0).UTC()
		fmt.Fprintf(w, "%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
	}
	io.WriteString(w, "\n")

	return f.formatBody(w, ";", ";", decimal, true)
}

func (f *Table) Plain(w io.Writer) error {
	if f.Charts.Len() == 0 {
		return nil
	}

	return f.formatBody(w, ": ", ", ", ".", false)
}

func (f *Table) formatBody(
	w io.Writer,
	delim0, delim, decimal string,
	quoteName bool) error {
	var pattern string
	if quoteName {
		pattern = "\"%v\"%v"
	} else {
		pattern = "%v%v"
	}

	// f.Charts is asserted to non-empty without check
	col, _ := f.Charts.Column(-1)

	for _, x := range col.Top(f.Count) {
		line := f.Charts[x.Name]

		fmt.Fprintf(w, pattern, x.Name, delim0)

		for i := 0; i < len(line); i += f.Step {
			if i > 0 {
				io.WriteString(w, delim)
			}

			s := fmt.Sprintf("%.08g", line[i])
			if decimal != "." {
				s = strings.Replace(s, ".", decimal, 1)
			}
			io.WriteString(w, s)
		}
		io.WriteString(w, "\n")
	}

	return nil
}
