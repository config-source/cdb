package cdb

type ErrorResponse struct {
	Message string
}

func (er ErrorResponse) Error() string {
	return er.Message
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Message: msg,
	}
}
