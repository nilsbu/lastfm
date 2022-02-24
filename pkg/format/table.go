package format

import (
	"fmt"
	"io"
	"strings"

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

	for i := 0; i < f.Charts.Len(); i += f.Step {
		if i > 0 {
			io.WriteString(w, ";")
		}

		t := f.First.AddDate(0, 0, i).Time()
		fmt.Fprintf(w, "%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
	}
	io.WriteString(w, "\n")

	return f.formatBody(w, "", "\n", ";", ";", decimal, true)
}

func (f *Table) Plain(w io.Writer) error {
	if f.Charts.Len() == 0 {
		return nil
	}

	return f.formatBody(w, "", "\n", ": ", ", ", ".", false)
}

func (f *Table) HTML(w io.Writer) error {
	if f.Charts.Len() == 0 {
		return nil
	}

	io.WriteString(w, "<table>")
	defer io.WriteString(w, "</table>")

	return f.formatBody(w, "<tr><td>", "</td></tr>", "</td><td>", "</td><td>", ".", false)
}

func (f *Table) formatBody(
	w io.Writer,
	start, end, delim0, delim, decimal string,
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
		io.WriteString(w, start)

		var line []float64
		for i, key := range f.Charts.Keys {
			if key.String() == x.Name {
				line = f.Charts.Values[i]
				break
			}
		}

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
		io.WriteString(w, end)
	}

	return nil
}
