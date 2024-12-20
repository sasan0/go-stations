package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nowTimeDate := time.Now()

		h.ServeHTTP(w, r)

		latency := time.Since(nowTimeDate).Milliseconds()

		path := r.URL.Path
		os := r.Context().Value(model.UserAgent{}).(string)
		logStr := model.Time{Timestamp: nowTimeDate, Latency: latency, Path: path, OS: os}

		log, err := json.Marshal(logStr)

		if err != nil {
			fmt.Println("Log Error!")
		}
		fmt.Println(string(log))
	})
}
