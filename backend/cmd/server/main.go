// Command server provides the HTTP API for the Swipe-Files backend.
// It wires the HTTP mux from internal/server and starts listening.
package main

import (
    "log"
    "net/http"

    "Swipe-Files/backend/internal/server"
)

func main() {
    mux := server.NewMux()
    addr := ":8787"
    log.Printf("Backend up at %s", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}
