// Package respond contains helpers for writing HTTP responses.
package respond

import (
    "encoding/json"
    "net/http"
)

// JSON writes the given value with content-type application/json and status code.
func JSON(w http.ResponseWriter, v any, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}
