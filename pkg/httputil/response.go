package httputil

import (
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func NewResponse(code int, data interface{}) Response {
	return Response{Code: code, Data: data, Status: http.StatusText(code)}
}

func NewErrorResponse(err error) RestError {
	return *parseErrors(err)
}
