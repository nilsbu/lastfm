package format

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/nilsbu/lastfm/pkg/charts2"
)

type Charts struct {
	Charts     charts2.LazyCharts
	Column     int
	Count      int
	Numbered   bool
	Precision  int
	Percentage bool
}

func (f *Charts) CSV(w io.Writer, decimal string) error {
	colFormatter := f.column()
	if colFormatter == nil {
		return nil
	}
	return colFormatter.CSV(w, decimal)
}

func (f *Charts) Plain(w io.Writer) error {
	colFormatter := f.column()
	if colFormatter == nil {
		return nil
	}
	return colFormatter.Plain(w)
}

func (f *Charts) HTML(w io.Writer) error {
	colFormatter := f.column()
	if colFormatter == nil {
		return nil
	}
	return colFormatter.HTML(w)
}

func (f *Charts) column() *Column {
	col := charts2.Column(f.Charts, -1)
	cache := charts2.Cache(col)
	sumTotal := 0.
	if totals := charts2.ColumnSum(cache).Row(charts2.KeyTitle("total"), 0, 1); len(totals) > 0 {
		sumTotal = totals[0]
	}

	n := f.Count
	if n == 0 {
		n = 10
	}
	top := charts2.Only(cache, charts2.Top(cache, n))

	return &Column{
		Column:     top,
		Numbered:   f.Numbered,
		Precision:  f.Precision,
		Percentage: f.Percentage,
		SumTotal:   sumTotal,
	}
}

type Column struct {
	Column     charts2.LazyCharts
	Numbered   bool
	Precision  int
	Percentage bool
	SumTotal   float64
}

func (f *Column) CSV(w io.Writer, decimal string) error {
	var header string
	if f.Numbered {
		header = "\"#\";\"Name\";\"Value\"\n"
	} else {
		header = "\"Name\";\"Value\"\n"
	}

	return f.format(header, f.getCSVPattern(), decimal, w)
}

func (f *Column) Plain(w io.Writer) error {
	return f.format("", f.getPlainPattern(), ".", w)
}

func (f *Column) HTML(w io.Writer) error {
	io.WriteString(w, "<table>")
	defer io.WriteString(w, "</table>")
	return f.format("", f.getHTMLPattern(), ".", w)
}

func (f *Column) format(
	header, pattern, decimal string, w io.Writer) error {
	if f.Column.Len() == 0 {
		return nil
	}

	io.WriteString(w, header)

	var outCol charts2.LazyCharts
	if f.Percentage {
		outCol = f.getPercentageColumn()
	} else {
		outCol = f.Column
	}

	data := outCol.Data(outCol.Titles(), 0, outCol.Len())
	for i, title := range f.Column.Titles() {
		sscore := fmt.Sprintf(f.getScorePattern(), data[title.Key()].Line[0])
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

func (f *Column) getCSVPattern() (pattern string) {
	if f.Numbered {
		pattern = "%d;"
	}

	pattern += "\"%v\";%v\n"

	return pattern
}

func (f *Column) getPlainPattern() (pattern string) {
	if f.Numbered {
		width := int(math.Log10(float64(len(f.Column.Titles())))) + 1
		pattern = "%" + strconv.Itoa(width) + "d: "
	}

	maxNameLen := strconv.Itoa(f.getMaxNameLen())
	pattern += "%-" + maxNameLen + "v - %v\n"

	return pattern
}

func (f *Column) getHTMLPattern() (pattern string) {

	pattern = "<tr>"
	if f.Numbered {
		pattern += "<td>%d</td>"
	}

	pattern += "<td>%v</td><td>%v</td></tr>"

	return pattern
}

func (f *Column) getScorePattern() (pattern string) {
	var maxValueLen int
	if len(f.Column.Titles()) == 0 || f.Column.Row(f.Column.Titles()[0], 0, 1)[0] == 0 {
		maxValueLen = 1
	} else {
		maxValueLen = int(math.Log10(f.Column.Row(f.Column.Titles()[0], 0, 1)[0])) + 1
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

func (f *Column) getPercentageColumn() charts2.LazyCharts {
	sum := f.SumTotal

	result := map[string][]float64{}
	if sum > 0 {
		for _, line := range f.Column.Data(f.Column.Titles(), 0, f.Column.Len()) {
			result[line.Title.Key()] = []float64{100 * line.Line[len(line.Line)-1] / sum}
		}
	} else {
		for _, line := range f.Column.Data(f.Column.Titles(), 0, f.Column.Len()) {
			result[line.Title.Key()] = []float64{0}
		}
	}

	return charts2.FromMap(result)
}

func (f *Column) getMaxNameLen() int {
	maxLen := 0
	for _, title := range f.Column.Titles() {
		runeCnt := utf8.RuneCountInString(title.String())
		if maxLen < runeCnt {
			maxLen = runeCnt
		}
	}
	return maxLen
}
