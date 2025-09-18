// Package server: kokoaa backendin HTTP-reitityksen.
package server

import (
	"net/http"

	"Swipe-Files/backend/internal/http/handlers"
	"Swipe-Files/backend/internal/http/middleware"
)

// NewMux: määrittää reitit ja väliohjelmat.
// Reitit:
//
//	GET  /api/files   -> listaa hakemiston tiedostot
//	POST /api/trash   -> siirtää tiedoston roskikseen
//	GET  /api/open    -> palvelee tiedoston tavut inline
//	GET  /api/convert -> muuntaa toimisto-PDF:ksi (jos mahdollista)
//	GET  /healthz     -> elinvoimatesti
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/files", middleware.WithCORS(handlers.ListFiles))
	mux.HandleFunc("/api/trash", middleware.WithCORS(handlers.Trash))
    mux.HandleFunc("/api/open", middleware.WithCORS(handlers.Open))
    mux.HandleFunc("/api/convert", middleware.WithCORS(handlers.Convert))
    mux.HandleFunc("/api/recents", middleware.WithCORS(handlers.Recents))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}
