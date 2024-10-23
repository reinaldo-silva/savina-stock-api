package response

type AppResponse struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message"`
	Total      *int64      `json:"total,omitempty"`
}

func NewAppResponse(data interface{}, message string, total *int64, statusCode ...int) AppResponse {
	status := 200
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	return AppResponse{
		StatusCode: status,
		Data:       data,
		Message:    message,
		Total:      total,
	}
}
