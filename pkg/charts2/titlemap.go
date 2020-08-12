package charts2

// TitleValue contains a Title and a single value.
type TitleValue struct {
	Title Title
	Value float64
}

// TitleValueMap maps from string to TitleValue.
type TitleValueMap map[string]TitleValue

// TitleLine contains a Title and a line of values.
type TitleLine struct {
	Title Title
	Line  []float64
}

// TitleLineMap maps from string to TitleLine.
type TitleLineMap map[string]TitleLine
