package httputil

import (
	"log/slog"
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Error  interface{} `json:"error"`
}

func NewResponse(code int, data interface{}) Response {
	return Response{Code: code, Data: data, Status: http.StatusText(code)}
}

func NewErrorResponse(err error) ErrorResponse {
	parsedError := parseErrors(err)
	slog.Error("API_ERROR", "cause", parsedError.Message)
	return ErrorResponse{
		Code:   parsedError.Code,
		Error:  parsedError.Message,
		Status: http.StatusText(parsedError.Code),
	}
}
