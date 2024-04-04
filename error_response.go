package cdb

type ErrorResponse struct {
	Message string `json:"message"`
}

func (er ErrorResponse) Error() string {
	return er.Message
}
