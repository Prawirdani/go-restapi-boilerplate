package middleware

import (
	"context"
	"net/http"

	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/jwt"
)

type CtxTokenPayloadKey string

const AccessTokenContextKey CtxTokenPayloadKey = "accessToken_payload"

func ValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenPayload, err := jwt.ValidateFromRequest(r, httputil.ACCESS_TOKEN_COOKIE_NAME)
		if err != nil {
			httputil.SendError(w, err)
			return
		}
		
		// Collect user payload from token then pass it to the next handler.
		ctx := context.WithValue(r.Context(), AccessTokenContextKey, map[string]string{
			"id":       accessTokenPayload["id"].(string),
			"username": accessTokenPayload["username"].(string),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
