package errs

type PkgError struct {
	message string
}

func (e *PkgError) Error() string {
	return e.message
}

func NewPkgError(message string) error {
	return &PkgError{message: message}
}
