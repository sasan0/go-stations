package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type UserAgent struct{}

func IdentifyDevice(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserAgent{}, ua)))
	})
}
