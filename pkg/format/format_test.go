package format

import "testing"

func TestAsserts(t *testing.T) {
	AssertFormatterInterface()
}

func AssertFormatterInterface() {
	var _ Formatter = (*Charts)(nil)
	var _ Formatter = (*Column)(nil)
	var _ Formatter = (*Error)(nil)
	var _ Formatter = (*Message)(nil)
	var _ Formatter = (*Table)(nil)
}
