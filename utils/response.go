package utils

import (
	"fmt"
	"net/http"
)

type Response struct {
	Message string  `json:"message"`
	Status  int     `json:"status"`
	Error   *string `json:"error,omitempty"`
}

func NewResponse(message string, status int) *Response {
	return &Response{Message: message, Status: status}
}

func GetPutSuccessfullResponse(key string) *Response {
	return NewResponse(fmt.Sprintf("Successfully put key %s", key), http.StatusOK)
}

func GetDeleteSuccessfullResponse(key string) *Response {
	return NewResponse(fmt.Sprintf("Successfully deleted key %s", key), http.StatusOK)
}

// ----

func NewErrResponse(message string, status int, errType string) *Response {
	return &Response{Message: message, Status: status, Error: &errType}
}

func GetBadRequestError(message string) *Response {
	return NewErrResponse(message, http.StatusBadRequest, "bad_request")
}

func GetNotFoundError(key string) *Response {
	return NewErrResponse(fmt.Sprintf("Key %s not found!", key), http.StatusNotFound, "not_found")
}

func GetInternalServerError(message string) *Response {
	return NewErrResponse(message, http.StatusInternalServerError, "internal_server_error")
}
