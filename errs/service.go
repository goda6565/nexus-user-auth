package errs

type ServiceError struct {
	message string
}

func (e *ServiceError) Error() string {
	return e.message
}

func NewServiceError(message string) error {
	return &ServiceError{message: message}
}
