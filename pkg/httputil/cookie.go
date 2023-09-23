package httputil

import (
	"net/http"
	"time"
)

const (
	REFRESH_TOKEN_COOKIE_NAME = "refreshToken"
	ACCESS_TOKEN_COOKIE_NAME  = "accessToken"
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

func SetCookieAccessToken(tokenValue string, rw http.ResponseWriter) {
	SetCookies(ACCESS_TOKEN_COOKIE_NAME, tokenValue, time.Now().Add(1*time.Minute), rw)
}

func SetCookieRefreshToken(tokenValue string, rw http.ResponseWriter) {
	SetCookies(REFRESH_TOKEN_COOKIE_NAME, tokenValue, time.Now().Add(1*time.Hour), rw)
}
