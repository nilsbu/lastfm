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

func (formatter *Charts) CSV(w io.Writer, decimal string) error {
	colFormatter := formatter.column()
	if colFormatter == nil {
		return nil
	}
	return colFormatter.CSV(w, decimal)
}

func (formatter *Charts) Plain(w io.Writer) error {
	colFormatter := formatter.column()
	if colFormatter == nil {
		return nil
	}
	return colFormatter.Plain(w)
}

func (formatter *Charts) column() *Column {
	col, err := formatter.Charts.Column(formatter.Column)
	if err != nil {
		return nil
	}

	sumTotal := col.Sum()

	n := formatter.Count
	if n == 0 {
		n = 10
	}
	top := col.Top(n)

	return &Column{
		Column:     top,
		Numbered:   formatter.Numbered,
		Precision:  formatter.Precision,
		Percentage: formatter.Percentage,
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

func (formatter *Column) CSV(w io.Writer, decimal string) error {
	return formatter.format(formatter.getCSVPattern(), decimal, w)
}

func (formatter *Column) Plain(w io.Writer) error {
	return formatter.format(formatter.getPlainPattern(), ".", w)
}

func (formatter *Column) format(
	pattern, decimal string, w io.Writer) error {
	if len(formatter.Column) == 0 {
		return nil
	}

	var outCol charts.Column
	if formatter.Percentage {
		outCol = formatter.getPercentageColumn()
	} else {
		outCol = formatter.Column
	}

	for i, score := range outCol {
		sscore := fmt.Sprintf(formatter.getScorePattern(), score.Score)
		if decimal != "." {
			sscore = strings.Replace(sscore, ".", decimal, 1)
		}

		var str string
		if formatter.Numbered {
			str = fmt.Sprintf(pattern, i+1, score.Name, sscore)
		} else {
			str = fmt.Sprintf(pattern, score.Name, sscore)
		}
		if _, err := io.WriteString(w, str); err != nil {
			return err
		}
	}

	return nil
}

func (formatter *Column) getCSVPattern() (pattern string) {
	if formatter.Numbered {
		pattern = "%d;"
	}
	if formatter.Percentage {
		pattern += "\"%v\";%v;\n"
	} else {
		pattern += "\"%v\";%v;\n"
	}

	return pattern
}

func (formatter *Column) getPlainPattern() (pattern string) {
	if formatter.Numbered {
		width := int(math.Log10(float64(len(formatter.Column)))) + 1
		pattern = "%" + strconv.Itoa(width) + "d: "
	}

	maxNameLen := strconv.Itoa(formatter.getMaxNameLen())
	pattern += "%-" + maxNameLen + "v - %v\n"

	return pattern
}

func (formatter *Column) getScorePattern() (pattern string) {
	var maxValueLen int
	if len(formatter.Column) == 0 || formatter.Column[0].Score == 0 {
		maxValueLen = 1
	} else {
		maxValueLen = int(math.Log10(formatter.Column[0].Score)) + 1
	}
	if formatter.Precision > 0 {
		maxValueLen += 1 + formatter.Precision
	}

	strLen := strconv.Itoa(maxValueLen)
	pattern += "%" + strLen + "." + strconv.Itoa(formatter.Precision) + "f"

	if formatter.Percentage {
		pattern += "%%"
	}

	return pattern
}

func (formatter *Column) getPercentageColumn() charts.Column {
	var sum float64
	if formatter.SumTotal == 0.0 {
		sum = formatter.Column.Sum()
	} else {
		sum = formatter.SumTotal
	}

	result := charts.Column{}
	if sum > 0 {
		for _, line := range formatter.Column {
			result = append(result, charts.Score{
				Name:  line.Name,
				Score: 100 * line.Score / sum})
		}
	} else {
		for _, line := range formatter.Column {
			result = append(result, charts.Score{Name: line.Name, Score: 0})
		}
	}

	return result
}

func (formatter *Column) getMaxNameLen() int {
	maxLen := 0
	for _, score := range formatter.Column {
		runeCnt := utf8.RuneCountInString(score.Name)
		if maxLen < runeCnt {
			maxLen = runeCnt
		}
	}
	return maxLen
}
