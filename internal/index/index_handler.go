package index

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/json"
)

type IndexHandler struct{}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) Routes(r chi.Router) {
	r.Get("/", h.hello)
}

func (h *IndexHandler) hello(w http.ResponseWriter, r *http.Request) {
	json.Send(w, http.StatusOK, "Hello World")
}
