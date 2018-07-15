package format

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"unicode/utf8"

	"github.com/nilsbu/lastfm/pkg/charts"
)

type Charts struct {
	Charts    charts.Charts
	Column    int
	Count     int
	Numbered  bool
	Precision int
}

func (formatter *Charts) Plain(w io.Writer) {
	col, err := formatter.Charts.Column(formatter.Column)
	if err != nil {
		return
	}

	n := formatter.Count
	if n == 0 {
		n = 10
	}
	top := col.Top(n)

	colFormatter := &Column{
		Column:    top,
		Numbered:  formatter.Numbered,
		Precision: formatter.Precision,
	}

	colFormatter.Plain(w)
}

type Column struct {
	Column    charts.Column
	Numbered  bool
	Precision int
}

func (formatter *Column) Plain(w io.Writer) {
	if len(formatter.Column) == 0 {
		return
	}

	var pattern string

	if formatter.Numbered {
		width := int(math.Log10(float64(len(formatter.Column)))) + 1
		pattern += "%" + strconv.Itoa(width) + "d: "
	}

	maxNameLen := strconv.Itoa(formatter.getMaxNameLen())
	pattern += "%-" + maxNameLen + "v - "

	maxValueLen := int(math.Log10(formatter.Column[0].Score)) + 1
	if formatter.Precision > 0 {
		maxValueLen += 1 + formatter.Precision
	}
	strLen := strconv.Itoa(maxValueLen)
	pattern += "%" + strLen + "." + strconv.Itoa(formatter.Precision) + "f\n"

	if formatter.Numbered {
		for i, score := range formatter.Column {
			str := fmt.Sprintf(pattern, i+1, score.Name, score.Score)
			io.WriteString(w, str)
		}
	} else {
		for _, score := range formatter.Column {
			str := fmt.Sprintf(pattern, score.Name, score.Score)
			io.WriteString(w, str)
		}
	}
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
