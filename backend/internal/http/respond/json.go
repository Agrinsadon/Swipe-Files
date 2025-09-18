// Package respond: apurit HTTP-vastauksiin.
package respond

import (
    "encoding/json"
    "net/http"
)

// JSON: kirjoita v JSON:na annetulla statuskoodilla.
func JSON(w http.ResponseWriter, v any, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}
