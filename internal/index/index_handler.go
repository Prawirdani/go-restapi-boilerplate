package index

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/jason"
)

type IndexHandler struct{}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) Routes(r chi.Router) {
	r.Get("/", h.hello)
}

func (h *IndexHandler) hello(w http.ResponseWriter, r *http.Request) {
	jason.Send(w, http.StatusOK, "Hello World")
}
