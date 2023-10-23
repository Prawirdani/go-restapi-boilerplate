package jwt

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prawirdani/go-restapi-boilerplate/internal/user"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
)

var jwtSecret = []byte(
	"akBfd4k+enMwQ61edGpfsu3uLvxXa9aIlM0MIGm6BobvIGA/r3xUY0CqCyGpl65cp8ytxr1gg8Ssp9SEmDOEGQ==",
)

type TokenType int8

const (
	// Type for refresh token
	REFRESH_TOKEN TokenType = iota

	// Type for access token
	ACCESS_TOKEN
)

const (
	// Constant key for extract id value from JWT Map Claims
	CLAIMS_KEY_ID = "id"
	// Constant key for extract username value from JWT Map Claims
	CLAIMS_KEY_USERNAME = "username"
)

// JWT Payload Claims
type Claims struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

type Token struct {
	Type  TokenType
	Value string
}

type TokenPair struct {
	RefreshToken *Token
	AccessToken  *Token
}

// Sign Token as Refresh Token
func (t *Token) SignAsRefreshToken(user user.User) error {
	t.Type = REFRESH_TOKEN
	tokenString, err := generateToken(user, REFRESH_TOKEN)
	if err == nil {
		t.Value = tokenString
	}
	return err
}

// Sign Token as Access Token
func (t *Token) SignAsAccessToken(user user.User) error {
	t.Type = ACCESS_TOKEN
	tokenString, err := generateToken(user, ACCESS_TOKEN)
	if err == nil {
		t.Value = tokenString
	}
	return err
}

// Set token value into cookie
func (t *Token) SetToCookie(w http.ResponseWriter) {
	timeNow := time.Now()
	switch t.Type {
	case ACCESS_TOKEN:
		httputil.SetCookies(httputil.ACT_COOKIE_NAME, t.Value, timeNow.Add(1*time.Minute), w)
	default:
		httputil.SetCookies(httputil.RFT_COOKIE_NAME, t.Value, timeNow.Add(60*time.Minute), w)
	}
}

// Sign a pair of jwt tokens.
func (tp *TokenPair) SignPair(user user.User) error {
	tp.RefreshToken = new(Token)
	tp.AccessToken = new(Token)

	if err := tp.RefreshToken.SignAsRefreshToken(user); err != nil {
		return err
	}
	if err := tp.AccessToken.SignAsAccessToken(user); err != nil {
		return err
	}

	return nil
}

// Sign the values of both tokens to cookies
func (tp *TokenPair) SetToCookies(w http.ResponseWriter) {
	tp.RefreshToken.SetToCookie(w)
	tp.AccessToken.SetToCookie(w)
}

// Token generator, based on user data
func generateToken(u user.User, tokenType TokenType) (string, error) {
	timeNow := time.Now()
	claims := &Claims{
		Id: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(timeNow),
		},
	}
	switch tokenType {
	case REFRESH_TOKEN:
		claims.ExpiresAt = jwt.NewNumericDate(timeNow.Add(60 * time.Minute))
	case ACCESS_TOKEN:
		claims.Username = u.Username
		claims.ExpiresAt = jwt.NewNumericDate(timeNow.Add(30 * time.Second))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, err := token.SignedString(jwtSecret)
	if err != nil {
		slog.Error("JWT Signing error", slog.Any("cause", err))
		return "", err
	}
	return sign, nil
}

// Validate token and return map claims
func ValidateFromRequest(r *http.Request, tokenCookieName string) (map[string]interface{}, error) {
	tokenString := ""

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = authHeader[len("Bearer "):]
	} else {
		tokenString = httputil.GetCookieValue(r, tokenCookieName)
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, httputil.ErrInternalServer(
				fmt.Errorf("unexpected signing method: %v", t.Header["alg"]),
			)
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, httputil.ErrUnauthorized("invalid or expired token")
	}

	return claims, nil
}
