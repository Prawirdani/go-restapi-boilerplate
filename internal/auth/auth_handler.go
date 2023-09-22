package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(us AuthService) *AuthHandler {
	return &AuthHandler{
		authService: us,
	}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Post("/auth/login", h.Login)
}

//	@Summary		Login
//	@Description	Login
//	@Accept			json
//	@Param			User	body	auth.LoginRequest	true	"Login Payload"
//	@Produce		json
//	@Tags			Auth
//	@Success		200		{object}	httputil.Response
//	@Failure		default	{object}	httputil.ErrorResponse	"400 & 500 status, error field can be string or object"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody LoginRequest
	if err := httputil.BindJson(r, &reqBody); err != nil {
		httputil.SendError(w, err)
		return
	}

	token, err := h.authService.Login(r.Context(), reqBody)
	if err != nil || token == "" {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJson(w, 200, token)
}
