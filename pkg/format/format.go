package format

import "io"

// Formatter contains functions to format something in various formats.
type Formatter interface {
	Plain(io.Writer)
}
