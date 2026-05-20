package utils

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func SuccessResponse(message string, data interface{}, meta *Meta) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func ErrorResponse(message string, errorCode string, details interface{}) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   errorCode,
		Details: details,
	}
}