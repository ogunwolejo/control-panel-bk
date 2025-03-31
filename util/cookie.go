package util

import "net/http"

// Allows me to set httpOnly cookie

func SetHttpOnlyCookie(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
}


// Example of a http.Cookie setup
/* http.Cookie {
	Name: "",
	Value: "",
	HttpOnly: true,
	Secure: true,
	path: "/",
	SameSite: http.SameSiteStrictMode,
	Expires: time.Now().Add()
} */