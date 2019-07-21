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

func (f *Charts) column() *Column {
	col, err := f.Charts.FullTitleColumn(f.Column)
	if err != nil {
		return nil
	}

	sumTotal := col.Sum()

	n := f.Count
	if n == 0 {
		n = 10
	}
	top := col.Top(n)

	return &Column{
		Column:     top,
		Numbered:   f.Numbered,
		Precision:  f.Precision,
		Percentage: f.Percentage,
		SumTotal:   sumTotal,
	}
}

type Column struct {
	Column     charts.Column
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

func (f *Column) format(
	header, pattern, decimal string, w io.Writer) error {
	if len(f.Column) == 0 {
		return nil
	}

	io.WriteString(w, header)

	var outCol charts.Column
	if f.Percentage {
		outCol = f.getPercentageColumn()
	} else {
		outCol = f.Column
	}

	for i, score := range outCol {
		sscore := fmt.Sprintf(f.getScorePattern(), score.Score)
		if decimal != "." {
			sscore = strings.Replace(sscore, ".", decimal, 1)
		}

		if f.Numbered {
			fmt.Fprintf(w, pattern, i+1, score.Name, sscore)
		} else {
			fmt.Fprintf(w, pattern, score.Name, sscore)
		}
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
		width := int(math.Log10(float64(len(f.Column)))) + 1
		pattern = "%" + strconv.Itoa(width) + "d: "
	}

	maxNameLen := strconv.Itoa(f.getMaxNameLen())
	pattern += "%-" + maxNameLen + "v - %v\n"

	return pattern
}

func (f *Column) getScorePattern() (pattern string) {
	var maxValueLen int
	if len(f.Column) == 0 || f.Column[0].Score == 0 {
		maxValueLen = 1
	} else {
		maxValueLen = int(math.Log10(f.Column[0].Score)) + 1
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

func (f *Column) getPercentageColumn() charts.Column {
	var sum float64
	if f.SumTotal == 0.0 {
		sum = f.Column.Sum()
	} else {
		sum = f.SumTotal
	}

	result := charts.Column{}
	if sum > 0 {
		for _, line := range f.Column {
			result = append(result, charts.Score{
				Name:  line.Name,
				Score: 100 * line.Score / sum})
		}
	} else {
		for _, line := range f.Column {
			result = append(result, charts.Score{Name: line.Name, Score: 0})
		}
	}

	return result
}

func (f *Column) getMaxNameLen() int {
	maxLen := 0
	for _, score := range f.Column {
		runeCnt := utf8.RuneCountInString(score.Name)
		if maxLen < runeCnt {
			maxLen = runeCnt
		}
	}
	return maxLen
}
