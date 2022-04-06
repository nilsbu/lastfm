package format

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/nilsbu/async"
	"github.com/nilsbu/lastfm/pkg/charts"
)

type Charts struct {
	Charts     []charts.Charts
	Ranges     charts.Ranges
	Numbered   bool
	Precision  int
	Percentage bool
}

type data struct {
	titles [][]charts.Title
	values [][][]float64
}

func prep(f *Charts) (*data, error) {
	titles := make([][]charts.Title, len(f.Charts))
	values := make([][][]float64, len(f.Charts))
	err := async.Pie(len(f.Charts), func(i int) error {
		f.Charts[i] = charts.Cache(charts.Column(f.Charts[i], -1))

		titles[i] = f.Charts[i].Titles()

		var err error
		values[i], err = f.Charts[i].Data(titles[i], 0, 1)
		if err != nil {
			return err
		}
		return nil
	})
	return &data{titles: titles, values: values}, err
}

func (f *Charts) CSV(w io.Writer, decimal string) error {
	d, err := prep(f)
	if err != nil {
		return err
	}
	return format(d, f, initCsvWriter(f, w, decimal))
}

func (f *Charts) Plain(w io.Writer) error {
	d, err := prep(f)
	if err != nil {
		return err
	}
	return format(d, f, initPlainWriter(d, f, w))
}

func (f *Charts) HTML(w io.Writer) error {
	fmt.Fprint(w, "<table>")
	defer fmt.Fprint(w, "</table>")

	d, err := prep(f)
	if err != nil {
		return err
	}
	return format(d, f, initHtmlWriter(f, w))
}

func format(d *data, f *Charts, p chartsWriter) error {
	if len(d.values) == 0 {
		return nil
	}

	if p.hasHeader() {
		p.lineStart()
		if f.Numbered {
			p.headerNumbers()
		}
		for c := range d.values {
			p.headerTitle(c)
			p.headerValue(c)
			if c+1 < len(d.values) {
				p.columnBreak()
			}
		}
		p.lineEnd()
	}

	m := 1.0
	if f.Percentage {
		m = 100.0
	}

	// // first charts determines the number of lines
	n := len(d.titles[0])
	for l := 0; l < n; l++ {
		p.lineStart()
		if f.Numbered {
			p.lineNumber(l)
		}
		for c, titles := range d.titles {
			p.lineTitle(titles[l].String(), l, c)

			p.lineValue(m*d.values[c][l][0], l, c)
			if c+1 < len(f.Charts) {
				p.columnBreak()
			}
		}
		p.lineEnd()
	}

	return nil
}

type chartsWriter interface {
	lineStart()
	lineEnd()
	columnBreak()
	hasHeader() bool
	headerNumbers()
	headerTitle(c int)
	headerValue(c int)
	lineNumber(l int)
	lineTitle(title string, l, c int)
	lineValue(value float64, l, c int)
}

type csvWriter struct {
	w            io.Writer
	decimal      string
	valuePattern string
}

func initCsvWriter(f *Charts, w io.Writer, decimal string) chartsWriter {
	return &csvWriter{
		w:            w,
		decimal:      decimal,
		valuePattern: numberPattern(0, f.Precision, f.Percentage, false),
	}
}

func (f *csvWriter) lineStart() {
}

func (f *csvWriter) lineEnd() {
	fmt.Fprint(f.w, "\n")
}

func (f *csvWriter) columnBreak() {
	fmt.Fprint(f.w, ";")
}

func (f *csvWriter) hasHeader() bool {
	return true
}

func (f *csvWriter) headerNumbers() {
	fmt.Fprint(f.w, `"#";`)
}

func (f *csvWriter) headerTitle(c int) {
	fmt.Fprint(f.w, `"Name";`)
}

func (f *csvWriter) headerValue(c int) {
	fmt.Fprint(f.w, `"Value"`)
}

func (f *csvWriter) lineNumber(l int) {
	fmt.Fprintf(f.w, "%d;", l+1)
}

func (f *csvWriter) lineTitle(title string, l, c int) {
	fmt.Fprintf(f.w, `"%v";`, title)
}

func (f *csvWriter) lineValue(value float64, l, c int) {
	if f.decimal != "." {
		s := fmt.Sprintf(f.valuePattern, value)
		fmt.Fprint(f.w, strings.Replace(s, ".", f.decimal, 1))
	} else {
		fmt.Fprintf(f.w, f.valuePattern, value)
	}
}

type plainWriter struct {
	c *Charts
	w io.Writer

	header        bool
	numPattern    string
	titleLens     []int
	valuePatterns []string
	valueLens     []int
}

func initPlainWriter(d *data, c *Charts, w io.Writer) chartsWriter {
	header := len(c.Charts) == len(c.Ranges.Delims)

	var numPattern string
	if c.Numbered && len(d.titles) > 0 {
		numPattern = "%" + strconv.Itoa(int(math.Log10(float64(len(d.titles[0]))))+1) + "d: "
	}

	titleLens := make([]int, len(d.titles))
	for i, titles := range d.titles {
		l := getMaxNameLen(titles)
		if header {
			ll := len(c.Ranges.Delims[i].String())
			if l < ll {
				l = ll
			}
		}
		titleLens[i] = l
	}

	valuePatterns := make([]string, len(d.values))
	valueLens := make([]int, len(d.titles))
	for i, values := range d.values {
		var m float64
		if len(values) > 0 {
			m = values[0][0]
		}
		valuePatterns[i] = numberPattern(m, c.Precision, c.Percentage, true)
		valueLens[i] = maxValueLen(m, c.Precision, c.Percentage)
		fmt.Println(valueLens[i])
	}

	p := &plainWriter{c: c, w: w,
		header:        header,
		numPattern:    numPattern,
		titleLens:     titleLens,
		valuePatterns: valuePatterns,
		valueLens:     valueLens,
	}

	return p
}

func (f *plainWriter) lineStart() {
}

func (f *plainWriter) lineEnd() {
	fmt.Fprint(f.w, "\n")
}

func (f *plainWriter) columnBreak() {
	fmt.Fprint(f.w, "\t")
}

func (f *plainWriter) hasHeader() bool {
	return f.header
}

func (f *plainWriter) headerNumbers() {
	fmt.Fprint(f.w, "#  ")
}

func (f *plainWriter) headerTitle(c int) {
	fmt.Fprint(f.w, f.c.Ranges.Delims[c], "   ")
}

func (f *plainWriter) headerValue(c int) {
	fmt.Fprint(f.w, strings.Repeat(" ", f.valueLens[c]))
}

func (f *plainWriter) lineNumber(l int) {
	fmt.Fprintf(f.w, f.numPattern, l+1)
}

func (f *plainWriter) lineTitle(title string, l, c int) {
	fmt.Fprintf(f.w, "%-"+strconv.Itoa(f.titleLens[c])+"v - ", title)
}

func (f *plainWriter) lineValue(value float64, l, c int) {
	fmt.Fprintf(f.w, f.valuePatterns[c], value)
}

type htmlWriter struct {
	w            io.Writer
	ranges       charts.Ranges
	header       bool
	valuePattern string
}

func initHtmlWriter(f *Charts, w io.Writer) chartsWriter {
	return &htmlWriter{
		w:            w,
		ranges:       f.Ranges,
		header:       len(f.Charts) == len(f.Ranges.Delims),
		valuePattern: "<td>" + numberPattern(0, f.Precision, f.Percentage, false) + "</td>",
	}
}

func (f *htmlWriter) lineStart() {
	fmt.Fprint(f.w, "<tr>")
}

func (f *htmlWriter) lineEnd() {
	fmt.Fprint(f.w, "</tr>")
}

func (f *htmlWriter) columnBreak() {
}

func (f *htmlWriter) hasHeader() bool {
	return f.header
}

func (f *htmlWriter) headerNumbers() {
	fmt.Fprint(f.w, "<td>#</td>")
}

func (f *htmlWriter) headerTitle(c int) {
	fmt.Fprintf(f.w, "<td>%v</td>", f.ranges.Delims[c])
}

func (f *htmlWriter) headerValue(c int) {
	fmt.Fprint(f.w, "<td></td>")
}

func (f *htmlWriter) lineNumber(l int) {
	fmt.Fprintf(f.w, "<td>%d</td>", l+1)
}

func (f *htmlWriter) lineTitle(title string, l, c int) {
	fmt.Fprintf(f.w, "<td>%v</td>", title)
}

func (f *htmlWriter) lineValue(value float64, l, c int) {
	fmt.Fprintf(f.w, f.valuePattern, value)
}

////// helpers

func numberPattern(topValue float64, precision int, percentage, align bool) string {
	var length int
	if align {
		length = maxValueLen(topValue, precision, percentage)
	}

	if percentage {
		length--
	}

	strLen := strconv.Itoa(length)
	pattern := "%" + strLen + "." + strconv.Itoa(precision) + "f"

	if percentage {
		pattern += "%%"
	}

	return pattern
}

func maxValueLen(topValue float64, precision int, percentage bool) int {
	var length int
	if topValue == 0 {
		length = 1
	} else {
		length = int(math.Log10(topValue)) + 1
	}
	if precision > 0 {
		length += 1 + precision
	}

	return length
}

func getMaxNameLen(titles []charts.Title) int {
	maxLen := 0
	for _, title := range titles {
		runeCnt := utf8.RuneCountInString(title.String())
		if maxLen < runeCnt {
			maxLen = runeCnt
		}
	}
	return maxLen
}
