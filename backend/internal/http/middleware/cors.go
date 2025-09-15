// Package middleware contains HTTP middleware used across endpoints.
package middleware

import "net/http"

// WithCORS adds permissive CORS headers and handles OPTIONS preflight.
func WithCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next(w, r)
    }
}
