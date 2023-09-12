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
	res := &HttpResponse{Code: code, Data: data}
	switch code {
	case http.StatusOK:
		res.Status = "OK"
	case http.StatusCreated:
		res.Status = "CREATED"
	case http.StatusBadRequest:
		res.Status = "BAD_REQUEST"
	case http.StatusUnauthorized:
		res.Status = "UNAUTHORIZED"
	case http.StatusNotFound:
		res.Status = "NOT_FOUND"
	case http.StatusInternalServerError:
		res.Status = "INTERNAL_SERVER_ERROR"
	default:
		res.Status = "INTERNAL_SERVER_ERROR"
	}
	return *res
}
