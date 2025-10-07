package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request logger : %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Response time: %v", time.Since(start))
	})
}