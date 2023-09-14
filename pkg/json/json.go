package json

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/prawirdani/go-restapi-boilerplate/pkg/httpError"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/utils"
)

func Send(w http.ResponseWriter, status_code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	response := utils.BuildHttpResponse(status_code, data)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("json.Send", "cause", err)
	}
}

func Bind(r *http.Request, request interface{}) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		slog.Error("json.Bind", "cause", err)
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, err error) {
	Send(w, httpError.WrapError(err).Code, httpError.WrapError(err).Message)
}
