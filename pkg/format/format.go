package format

import "io"

// Formatter contains functions to format something in various formats.
type Formatter interface {
	CSV(w io.Writer, decimal string) error
	Plain(io.Writer) error
	HTML(io.Writer) error
	JSON(io.Writer) error
}
