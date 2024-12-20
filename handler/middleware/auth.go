package middleware

import (
	"net/http"
	"os"
)

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()

		if user == os.Getenv("BASIC_AUTH_USER_ID") && password == os.Getenv("BASIC_AUTH_PASSWORD") && ok {
			h.ServeHTTP(w, r)
		} else {
			// Status Code 401
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
