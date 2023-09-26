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



//	@Summary		Login
//	@Description	Login
//	@Accept			json
//	@Param			User	body	auth.LoginRequest	true	"Login Payload"
//	@Produce		json
//	@Tags			Auth
//	@Success		200		{object}	httputil.Response{data=object{access_token=string,refresh_token=string}}
//	@Failure		default	{object}	httputil.RestError	"400 & 500 status, error field can be string or object"
//	@Header			200		{string}	refreshToken		"Refresh token automatically set up to user cookie header"
//	@Header			200		{string}	accessToken			"Access token automatically set up to user cookie header"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody LoginRequest
	if err := httputil.BindJson(r, &reqBody); err != nil {
		httputil.SendError(w, err)
		return
	}

	tokenPair, err := h.authService.Login(r.Context(), reqBody)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	tokenPair.SetToCookies(w)
	httputil.SendJson(w, 200, map[string]string{
		"refresh_token": tokenPair.RefreshToken.Value,
		"access_token":  tokenPair.AccessToken.Value,
	},
	)
}

//	@Summary		Me
//	@Description	Identify user login state 
//	@Produce		json
//	@Tags			Auth
//	@Success		200			{object}	httputil.Response{data=object{id=string,username=string}}
//	@Failure		default		{object}	httputil.RestError	"400 & 500 status, error field can be string or object"
//	@Param			accessToken	header		string				true	"accessToken cookie header"
//	@Security		JWT Cookies Auth
//	@Router			/auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accessTokenPayload := r.Context().Value(middleware.AC_TOKEN_PAYLOAD_CTX_KEY)
	httputil.SendJson(w, 200, accessTokenPayload)
}

//	@Summary		Refresh Token
//	@Description	Refresh new access token 
//	@Produce		json
//	@Tags			Auth
//	@Success		200				{object}	httputil.Response{data=object{access_token=string}}
//	@Failure		default			{object}	httputil.RestError	"400 & 500 status, error field can be string or object"
//	@Header			200				{string}	accessToken			"New Access token automatically set up to user cookie"
//	@Param			refreshToken	header		string				true	"refreshToken cookie header"
//	@Security		JWT Cookies Auth
//	@Router			/auth/refresh [get]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenClaims, err := jwt.ValidateFromRequest(r, httputil.RFT_COOKIE_NAME)
	if err != nil {
		httputil.SendError(w, err)
		return
	}
	
	newAccessToken, err := h.authService.RefreshToken(r.Context(), refreshTokenClaims["id"].(string))
	if err != nil {
		httputil.SendError(w, err)
		return
	}
	newAccessToken.SetToCookie(w)
	httputil.SendJson(w, 200, map[string]string{"access_token": newAccessToken.Value})
}

//	@Summary		Logout
//	@Description	Logout current user 
//	@Produce		json
//	@Tags			Auth
//	@Success		200			{object}	httputil.Response
//	@Failure		default		{object}	httputil.RestError	"400 & 500 status, error field can be string or object"
//	@Param			accessToken	header		string				true	"accessToken cookie header"
//	@Security		JWT Cookies Auth
//	@Router			/auth/logout [delete]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	httputil.RemoveCookie(w, httputil.ACT_COOKIE_NAME)
	httputil.RemoveCookie(w, httputil.RFT_COOKIE_NAME)
	
	httputil.SendJson(w, 200, "logout successfully!")
}
