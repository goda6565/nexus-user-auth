package errs

type InterfaceError struct {
	message string
}

func (e *InterfaceError) Error() string {
	return e.message
}

func NewInterfaceError(message string) error {
	return &InterfaceError{message: message}
}
