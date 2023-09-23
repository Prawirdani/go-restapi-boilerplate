package middleware

import (
	"context"
	"net/http"

	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/jwt"
)

type ContextKey string

// Constant key for access token payload from request context
const AC_TOKEN_PAYLOAD_CTX_KEY ContextKey = "accessTokenPayload"

func ValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenPayload, err := jwt.ValidateFromRequest(r, httputil.ACT_COOKIE_NAME)
		if err != nil {
			httputil.SendError(w, err)
			return
		}

		// Collect user payload from token then pass it to the next handler.
		ctx := context.WithValue(r.Context(), AC_TOKEN_PAYLOAD_CTX_KEY, map[string]string{
			"id":       accessTokenPayload[jwt.CLAIMS_KEY_ID].(string),
			"username": accessTokenPayload[jwt.CLAIMS_KEY_USERNAME].(string),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
