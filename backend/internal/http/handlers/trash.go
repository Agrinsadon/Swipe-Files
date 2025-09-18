package handlers

import (
    "encoding/json"
    "net/http"

    "Swipe-Files/backend/internal/http/respond"
    "Swipe-Files/backend/internal/util"
    "Swipe-Files/backend/platform"
)

type trashRequest struct {
    Path string `json:"path"`
}

// Trash: siirtää annetun polun tiedoston käyttöjärjestelmän roskakoriin.
// Odottaa JSON POST -rungon: {"path": "/absoluuttinen/tai/~/polku"}
func Trash(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "vain POST", http.StatusMethodNotAllowed)
        return
    }
    var req trashRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
        http.Error(w, "virheellinen JSON (tarvitaan {\"path\": \"...\"})", http.StatusBadRequest)
        return
    }
    abs, err := util.ResolvePath(req.Path)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := platform.MoveToTrash(abs); err != nil {
        http.Error(w, "roskikseen siirto epäonnistui: "+err.Error(), http.StatusInternalServerError)
        return
    }
    respond.JSON(w, map[string]string{"status": "ok", "path": abs}, http.StatusOK)
}
