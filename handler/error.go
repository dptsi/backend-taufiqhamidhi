package handler

type ErrorResponse struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}

func NewErrorResponse(code string, message string) *ErrorResponse {
	return &ErrorResponse{
		ResponseCode:    code,
		ResponseMessage: message,
	}
}
