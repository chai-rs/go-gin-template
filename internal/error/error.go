package error

// AppError represents an application error.
type AppError struct {
	Code        int    `json:"code"`
	SystemError error  `json:"-"`
	Message     string `json:"message"`
}

// Error returns the error message.
func (e *AppError) Error() string {
	return e.SystemError.Error()
}

// New creates a new application error.
func New(code int, systemError error, messages ...string) *AppError {
	message := ""
	if len(messages) > 0 {
		message = messages[0]
	}

	return &AppError{
		Code:        code,
		SystemError: systemError,
		Message:     message,
	}
}
