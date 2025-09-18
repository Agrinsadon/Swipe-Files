// Package middleware: jaetut HTTP-väliohjelmat (CORS ym.).
package middleware

import "net/http"

// WithCORS: lisää sallitut CORS-otsikot ja käsittelee OPTIONS-esipyynnöt.
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
