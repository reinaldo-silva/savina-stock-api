package error

type AppError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func NewAppError(message string, statusCode ...int) AppError {
	status := 400
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	return AppError{
		StatusCode: status,
		Message:    message,
	}
}
