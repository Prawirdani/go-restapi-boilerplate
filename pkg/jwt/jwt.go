package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/go-restapi-boilerplate/internal/user"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
)

var jwtSecret = []byte("akBfd4k+enMwQ61edGpfsu3uLvxXa9aIlM0MIGm6BobvIGA/r3xUY0CqCyGpl65cp8ytxr1gg8Ssp9SEmDOEGQ==")

type Claims struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

type TokenPairs struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (tp *TokenPairs) SetToCookies(rw http.ResponseWriter) {
	httputil.SetCookieAccessToken(tp.AccessToken, rw)
	httputil.SetCookieRefreshToken(tp.RefreshToken, rw)
}

// Todo it should contains UUID Store to Redis. When refresh new access token it will look up the redis with this UUID and return the payload as long as the refresh token is not yet expired.
/* Signing Refresh Token takes User Id as payload */
func SignRefreshToken(userID string) string {
	claims := &Claims{
		Id: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rfToken",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, err := token.SignedString(jwtSecret)
	if err != nil {
		panic(err)
	}
	return sign
}

/* Signing Access Token takes User Id and Username  as payload */
func SignAccessToken(user *user.User) string {
	claims := &Claims{
		Id:       user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "acToken",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, err := token.SignedString(jwtSecret)
	if err != nil {
		panic(err)
	}
	return sign
}

/* Sign Access and Refresh token at the same time, used for login */
func SignPairs(u *user.User) *TokenPairs {
	tokenPairs := &TokenPairs{
		RefreshToken: SignRefreshToken(u.Id),
		AccessToken:  SignAccessToken(u),
	}

	return tokenPairs
}


/* Returning Token Map Claims straighly from incoming request */
func ValidateFromRequest(r *http.Request, tokenCookieName string) (map[string]interface{}, error) {
	tokenString := httputil.GetCookieValue(r, tokenCookieName)
	if tokenString == "" {
		return nil, httputil.ErrUnauthorized("please login to proceed!")
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, httputil.ErrUnauthorized("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, httputil.ErrUnauthorized("invalid/expired token")
	}

	return claims, nil
}
