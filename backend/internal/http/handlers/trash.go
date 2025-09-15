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

// Trash moves the given file path to the OS trash/recycle bin.
// Accepts JSON POST body: {"path": "/absolute/or/~/path"}
func Trash(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "POST only", http.StatusMethodNotAllowed)
        return
    }
    var req trashRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
        http.Error(w, "invalid JSON (need {\"path\": \"...\"})", http.StatusBadRequest)
        return
    }
    abs, err := util.ResolvePath(req.Path)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := platform.MoveToTrash(abs); err != nil {
        http.Error(w, "trash failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    respond.JSON(w, map[string]string{"status": "ok", "path": abs}, http.StatusOK)
}
