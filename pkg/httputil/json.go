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
		slog.Error("httputil.SendJson", "cause", err)
	}
}

func BindJson(r *http.Request, request interface{}) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, Err error) {
	w.Header().Set("Content-Type", "application/json")
	response := NewErrorResponse(Err)
	w.WriteHeader(response.Code)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("httputil.SendErr", "cause", err)
	}
}
