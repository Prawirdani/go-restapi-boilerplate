package utils

import (
	"net/http"
)

type HttpResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func BuildHttpResponse(code int, data interface{}) HttpResponse {
	return HttpResponse{Code: code, Data: data, Status: http.StatusText(code)}
}
