package response

type AppResponse struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message"`
}

func NewAppResponse(data interface{}, message string, statusCode ...int) AppResponse {
	status := 200
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	return AppResponse{
		StatusCode: status,
		Data:       data,
		Message:    message,
	}
}
