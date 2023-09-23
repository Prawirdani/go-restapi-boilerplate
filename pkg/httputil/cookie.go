package httputil

import (
	"net/http"
	"time"
)

const (
	// Refresh token cookie name
	RFT_COOKIE_NAME = "refreshToken"
	// Access token cookie name
	ACT_COOKIE_NAME = "accessToken"
)

func GetCookieValue(r *http.Request, tokenName string) string {
	if cookie, err := r.Cookie(tokenName); err == nil {
		return cookie.Value
	}
	return ""
}

func RemoveCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func SetCookies(name string, value string, expires time.Time, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}
