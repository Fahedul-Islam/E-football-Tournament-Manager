package middleware

import (
	"log"
	"net/http"
)

func CorsWithPreflight(mux http.Handler) http.Handler {
	handleAllReq := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS middleware executed for %s %s", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return // Stop further processing
	}
		mux.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handleAllReq)
}