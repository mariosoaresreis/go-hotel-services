package errors

import "net/http"

type ApplicationError struct {
	Code    int    `json: ",omitempty"`
	Message string `json: message`
}

func (e ApplicationError) Text() *ApplicationError {
	return &ApplicationError{
		Message: e.Message,
	}
}

func NewNotFoundError(message string) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewBadRequestError(message string) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}
func NewUnexpectedError(message string) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
