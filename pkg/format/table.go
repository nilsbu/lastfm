package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/nilsbu/lastfm/pkg/charts"
)

type Table struct {
	Charts charts.Charts
	Ranges charts.Ranges
}

func (f *Table) CSV(w io.Writer, decimal string) error {
	if f.Charts.Len() == 0 {
		return nil
	}

	io.WriteString(w, "\"name\";")

	for i := 0; i < f.Charts.Len(); i++ {
		if i > 0 {
			io.WriteString(w, ";")
		}

		io.WriteString(w, f.Ranges.Delims[i].String())
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

	titles := f.Charts.Titles()
	data := f.Charts.Data(titles, 0, f.Charts.Len())

	for t, x := range titles {
		io.WriteString(w, start)

		fmt.Fprintf(w, pattern, x.String(), delim0)

		for i := 0; i < len(data[t]); i++ {
			if i > 0 {
				io.WriteString(w, delim)
			}

			s := fmt.Sprintf("%.08g", data[t][i])
			if decimal != "." {
				s = strings.Replace(s, ".", decimal, 1)
			}
			io.WriteString(w, s)
		}
		io.WriteString(w, end)
	}

	return nil
}
