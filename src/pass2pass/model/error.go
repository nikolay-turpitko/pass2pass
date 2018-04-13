package model

// FatalError is a marker interface of fatal errors.
type FatalError interface {
	IsFatal() bool
}

// NewFatalError wraps error with FatalError interface.
func NewFatalError(err error) error {
	return fatalError{err}
}

type fatalError struct {
	error
}

func (fatalError) IsFatal() bool {
	return true
}
