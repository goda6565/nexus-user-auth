package errs

type DomainError struct {
	message string // エラーメッセージ
}

func (e *DomainError) Error() string {
	return e.message
}

func NewDomainError(message string) error {
	return &DomainError{message: message}
}
