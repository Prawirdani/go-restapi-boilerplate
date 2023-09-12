package json

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prawirdani/go-restapi-boilerplate/pkg/utils"
)

func Send(w http.ResponseWriter, status_code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	response := utils.BuildHttpResponse(status_code, data)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func Bind(r *http.Request, request interface{}) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return err
	}
	return nil
}
