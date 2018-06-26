package fail

import "fmt"

// Severity indicates the amount of damage an error can do.
type Severity int

// Control denotes a failure of a function that impacts the flow of the program
// but on its own has no consequence for the user.
//
// Suspicious denotes a failure that ought to be investigated since it might
// hint towards a bug or data corruption but that might turn out to be innocent.
//
// Critical failures make the continuation of the program impossible. They may
// imply data corruption.
//
// TODO split critical into cases where data corruption has been safely avoided
// and those where uncontrolled behaviour has occurred.
const (
	Control Severity = iota
	Suspicious
	Critical
)

// Threat is an error with a severity level.
type Threat interface {
	error
	Severity() Severity
}

// AssessedError is an error with a severity level. It implements Threat.
type AssessedError struct {
	Sev Severity
	Err error
}

func (err *AssessedError) Error() string {
	raw := err.Err.Error()

	if err.Sev < 0 || err.Sev > Critical {
		return raw
	}
	return fmt.Sprintf("[%v] %v", GetSeverityString(err.Sev), err.Err)
}

func (err *AssessedError) Severity() Severity {
	return err.Sev
}

// GetSeverityString returns the severity level as a string.
func GetSeverityString(sev Severity) string {
	switch sev {
	case Control:
		return "control"
	case Suspicious:
		return "suspicious"
	case Critical:
		return "critical"
	default:
		return ""
	}
}
