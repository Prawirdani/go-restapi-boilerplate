package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func SendJson(w http.ResponseWriter, status_code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	response := NewResponse(status_code, data)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("json.Send", "cause", err)
	}
}

func BindJson(r *http.Request, request interface{}) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		slog.Error("json.Bind", "cause", err)
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, err error) {
	SendJson(w, parseErrors(err).Code, parseErrors(err).Message)
}
