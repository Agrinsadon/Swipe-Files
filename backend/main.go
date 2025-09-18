// Pääohjelma: käynnistää HTTP-palvelimen.
package main

import (
	"log"
	"net/http"

	"Swipe-Files/backend/internal/server"
)

func main() {
	mux := server.NewMux()
	addr := ":8787"
	log.Printf("Backend käynnissä %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
