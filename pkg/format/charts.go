package format

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/nilsbu/lastfm/pkg/charts"
)

type Charts struct {
	Charts     charts.Charts
	Numbered   bool
	Precision  int
	Percentage bool
}

func (f *Charts) CSV(w io.Writer, decimal string) error {
	var header string
	if f.Numbered {
		header = "\"#\";\"Name\";\"Value\"\n"
	} else {
		header = "\"Name\";\"Value\"\n"
	}

	return f.format(header, f.getCSVPattern(), decimal, w)
}

func (f *Charts) Plain(w io.Writer) error {
	return f.format("", f.getPlainPattern(), ".", w)
}

func (f *Charts) HTML(w io.Writer) error {
	io.WriteString(w, "<table>")
	defer io.WriteString(w, "</table>")
	return f.format("", f.getHTMLPattern(), ".", w)
}

func (f *Charts) format(
	header, pattern, decimal string, w io.Writer) error {
	if f.Charts.Len() == 0 {
		return nil
	}

	io.WriteString(w, header)

	var multi float64
	if f.Percentage {
		multi = 100
	} else {
		multi = 1
	}

	f.Charts = charts.Cache(charts.Column(f.Charts, -1))
	data, err := f.Charts.Data(f.Charts.Titles(), 0, 1)
	if err != nil {
		return err
	}

	scorepattern := f.getScorePattern()

	for i, title := range f.Charts.Titles() {
		sscore := fmt.Sprintf(scorepattern, multi*data[i][0])
		if decimal != "." {
			sscore = strings.Replace(sscore, ".", decimal, 1)
		}

		if f.Numbered {
			fmt.Fprintf(w, pattern, i+1, title, sscore)
		} else {
			fmt.Fprintf(w, pattern, title, sscore)
		}
		i++
	}

	return nil
}

func (f *Charts) getCSVPattern() (pattern string) {
	if f.Numbered {
		pattern = "%d;"
	}

	pattern += "\"%v\";%v\n"

	return pattern
}

func (f *Charts) getPlainPattern() (pattern string) {
	if f.Numbered {
		width := int(math.Log10(float64(len(f.Charts.Titles())))) + 1
		pattern = "%" + strconv.Itoa(width) + "d: "
	}

	maxNameLen := strconv.Itoa(f.getMaxNameLen())
	pattern += "%-" + maxNameLen + "v - %v\n"

	return pattern
}

func (f *Charts) getHTMLPattern() (pattern string) {

	pattern = "<tr>"
	if f.Numbered {
		pattern += "<td>%d</td>"
	}

	pattern += "<td>%v</td><td>%v</td></tr>"

	return pattern
}

func (f *Charts) getScorePattern() (pattern string) {
	titles := f.Charts.Titles()
	var topValue float64
	if len(titles) > 0 {
		// error can be ignored because data has already been requested earlier
		data, _ := f.Charts.Data(titles[:1], 0, 1)
		topValue = data[0][0]
	}

	var maxValueLen int
	if len(f.Charts.Titles()) == 0 || topValue == 0 {
		maxValueLen = 1
	} else {
		maxValueLen = int(math.Log10(topValue)) + 1
	}
	if f.Precision > 0 {
		maxValueLen += 1 + f.Precision
	}

	strLen := strconv.Itoa(maxValueLen)
	pattern += "%" + strLen + "." + strconv.Itoa(f.Precision) + "f"

	if f.Percentage {
		pattern += "%%"
	}

	return pattern
}

func (f *Charts) getMaxNameLen() int {
	maxLen := 0
	for _, title := range f.Charts.Titles() {
		runeCnt := utf8.RuneCountInString(title.String())
		if maxLen < runeCnt {
			maxLen = runeCnt
		}
	}
	return maxLen
}
