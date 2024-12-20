package middleware

import (
	"context"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mileusna/useragent"
)

func IdentifyDevice(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		os := useragent.Parse(r.UserAgent()).OS

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.UserAgent{}, os)))
	})
}
