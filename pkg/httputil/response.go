package httputil

import (
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Error  interface{} `json:"error,omitempty"`
}

func NewResponse(code int, data interface{}) Response {
	return Response{Code: code, Data: data, Status: http.StatusText(code)}
}

func NewErrorResponse(code int, data interface{}) ErrorResponse {
	return ErrorResponse{Code: code, Error: data, Status: http.StatusText(code)}
}

func BuildAPIResponse(code int, data interface{}) any {
	if code >= 400 {
		return NewErrorResponse(code, data)
	}
	return NewResponse(code, data)
}
