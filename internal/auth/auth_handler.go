package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/internal/middleware"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/jwt"
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
	r.Post("/login", h.Login)
	r.Get("/refresh", h.RefreshToken)
	r.With(middleware.ValidateAccessToken).Get("/me", h.Me)
	r.With(middleware.ValidateAccessToken).Delete("/logout", h.Logout)
}

// @Summary		Login
// @Description	Login
// @Accept			json
// @Param			User	body	auth.LoginRequest	true	"Login Payload"
// @Produce		json
// @Tags			Auth
// @Success		200		{object}	httputil.Response
// @Failure		default	{object}	httputil.ErrorResponse	"400 & 500 status, error field can be string or object"
// @Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody LoginRequest
	if err := httputil.BindJson(r, &reqBody); err != nil {
		httputil.SendError(w, err)
		return
	}

	tokenPairs, err := h.authService.Login(r.Context(), reqBody)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	tokenPairs.SetToCookies(w)
	httputil.SendJson(w, 200, tokenPairs)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accessTokenPayload := r.Context().Value(middleware.AccessTokenContextKey)
	httputil.SendJson(w, 200, accessTokenPayload)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenClaims, err := jwt.ValidateFromRequest(r, httputil.REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	newAccessToken, err := h.authService.RefreshToken(r.Context(), refreshTokenClaims["id"].(string))
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SetCookieAccessToken(newAccessToken, w)
	httputil.SendJson(w, 200, map[string]string{"access_token": newAccessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
}
