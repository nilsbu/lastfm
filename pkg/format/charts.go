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
	Numbered   bool
	Precision  int
	Percentage bool
}

func (f *Charts) CSV(w io.Writer, decimal string) error {
	var header string
	if f.Numbered {
		header = "\"#\";\"Name\";\"Value\""
	} else {
		header = "\"Name\";\"Value\""
	}

	return f.format(params{
		lineEnd:       "\n",
		header:        header,
		numberpattern: "%d;",
		pattern:       func(i int) string { return "\"%v\";%v" },
		chartsDelim:   ";",
		decimal:       decimal,
	}, w)
}

type params struct {
	lineStart, lineEnd string
	header             string
	numberpattern      string
	pattern            func(int) string
	chartsDelim        string
	decimal            string
}

func (f *Charts) Plain(w io.Writer) error {
	numberpattern := "%d: "

	if f.Numbered && len(f.Charts) > 0 {
		width := int(math.Log10(float64(len(f.Charts[0].Titles())))) + 1
		numberpattern = "%" + strconv.Itoa(width) + "d: "
	}

	return f.format(params{
		lineEnd:       "\n",
		numberpattern: numberpattern,
		pattern:       f.getPlainPattern,
		chartsDelim:   "\t",
		decimal:       ".",
	}, w)
}

func (f *Charts) HTML(w io.Writer) error {
	io.WriteString(w, "<table>")
	defer io.WriteString(w, "</table>")

	return f.format(params{
		lineStart:     "<tr>",
		lineEnd:       "</tr>",
		numberpattern: "<td>%d</td>",
		pattern:       func(i int) string { return "<td>%v</td><td>%v</td>" },
		decimal:       ".",
	}, w)
}

func (f *Charts) format(p params, w io.Writer) error {
	if len(f.Charts) == 0 {
		return nil
	}

	if p.header != "" {
		fmt.Fprint(w, p.lineStart)
		for j := range f.Charts {
			io.WriteString(w, p.header)
			if j+1 < len(f.Charts) {
				fmt.Fprint(w, p.chartsDelim)
			}
		}
		fmt.Fprint(w, p.lineEnd)
	}

	var multi float64
	if f.Percentage {
		multi = 100
	} else {
		multi = 1
	}

	data := make([][][]float64, len(f.Charts))
	err := async.Pie(len(f.Charts), func(i int) error {
		f.Charts[i] = charts.Cache(charts.Column(f.Charts[i], -1))

		var err error
		data[i], err = f.Charts[i].Data(f.Charts[i].Titles(), 0, 1)
		return err
	})
	if err != nil {
		return err
	}

	scorepatterns := make([]string, len(f.Charts))
	for i := range f.Charts {
		scorepatterns[i] = f.getScorePattern(i)
	}

	// first charts determines the number of lines
	n := len(f.Charts[0].Titles())
	for i := 0; i < n; i++ {
		fmt.Fprint(w, p.lineStart)
		if f.Numbered {
			fmt.Fprintf(w, p.numberpattern, i+1)
		}

		for j := range f.Charts {
			sscore := fmt.Sprintf(scorepatterns[j], multi*data[j][i][0])
			if p.decimal != "." {
				sscore = strings.Replace(sscore, ".", p.decimal, 1)
			}

			fmt.Fprintf(w, p.pattern(j), f.Charts[j].Titles()[i], sscore)
			if j+1 < len(f.Charts) {
				fmt.Fprint(w, p.chartsDelim)
			}
		}
		fmt.Fprint(w, p.lineEnd)
	}

	return nil
}

func (f *Charts) getPlainPattern(i int) (pattern string) {
	maxNameLen := strconv.Itoa(f.getMaxNameLen(i))
	pattern += "%-" + maxNameLen + "v - %v"

	return pattern
}

func (f *Charts) getScorePattern(j int) (pattern string) {
	titles := f.Charts[j].Titles()
	var topValue float64
	if len(titles) > 0 {
		// error can be ignored because data has already been requested earlier
		data, _ := f.Charts[j].Data(titles[:1], 0, 1)
		topValue = data[0][0]
	}

	var maxValueLen int
	if len(f.Charts[j].Titles()) == 0 || topValue == 0 {
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

func (f *Charts) getMaxNameLen(j int) int {
	maxLen := 0
	for _, title := range f.Charts[j].Titles() {
		runeCnt := utf8.RuneCountInString(title.String())
		if maxLen < runeCnt {
			maxLen = runeCnt
		}
	}
	return maxLen
}
