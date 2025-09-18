// Package server exposes HTTP router construction for the backend API.
package server

import (
    "net/http"

    "Swipe-Files/backend/internal/http/handlers"
    "Swipe-Files/backend/internal/http/middleware"
)

// NewMux wires routes and middleware.
// Routes:
//   GET  /api/files  -> list files in a directory
//   POST /api/trash  -> move a file to trash (platform specific)
//   GET  /healthz    -> health probe
func NewMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/files", middleware.WithCORS(handlers.ListFiles))
    mux.HandleFunc("/api/trash", middleware.WithCORS(handlers.Trash))
    mux.HandleFunc("/api/open", middleware.WithCORS(handlers.Open))
    mux.HandleFunc("/api/convert", middleware.WithCORS(handlers.Convert))
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })
    return mux
}
