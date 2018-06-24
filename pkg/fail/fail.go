package fail

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

// Threat indicates severity levels of an error.
type Threat interface {
	Severity() Severity
}

// AssessedError is an error with a severity level. It implements Threat.
type AssessedError struct {
	Sev Severity
	Err error
}

func (err *AssessedError) Error() string {
	raw := err.Err.Error()
	switch err.Sev {
	case Control:
		return "[control] " + raw
	case Suspicious:
		return "[suspicious] " + raw
	case Critical:
		return "[critical] " + raw
	default:
		return raw
	}
}

func (err *AssessedError) Severity() Severity {
	return err.Sev
}
