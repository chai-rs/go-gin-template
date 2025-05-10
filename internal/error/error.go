package error

type AppError struct {
	Code        int    `json:"code"`
	SystemError error  `json:"-"`
	Message     string `json:"message"`
}

func (e *AppError) Error() string {
	return e.SystemError.Error()
}

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
