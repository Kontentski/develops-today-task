package errs

// Err implements the Error interface with error marshaling.
type Err struct {
	Message string `json:"message"`
}

func New(message string) error {
	return &Err{
		Message: message,
	}
}

func (e *Err) Error() string {
	return e.Message
}

// IsExpected finds Err{} inside passed error.
func IsExpected(e error) bool {
	_, ok := e.(*Err)
	return ok
}
