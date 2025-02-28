package errs

type InfraError struct {
	message string
}

func (e *InfraError) Error() string {
	return e.message
}

func NewInfraError(message string) error {
	return &InfraError{message: message}
}
