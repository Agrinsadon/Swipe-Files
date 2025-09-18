package handlers

import (
    "net/http"
    "os"
    "path/filepath"

    "Swipe-Files/backend/internal/util"
)

// Open serves the file bytes for a given path so the frontend can preview
// images (and other inline-viewable types). Query: ?path=... (absolute or ~).
// Uses util.ResolvePath to expand ~ and make absolute, then serves inline.
func Open(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        // CORS preflight handled by middleware, but keep explicit return.
        w.WriteHeader(http.StatusNoContent)
        return
    }
    if r.Method != http.MethodGet {
        http.Error(w, "GET only", http.StatusMethodNotAllowed)
        return
    }
    p := r.URL.Query().Get("path")
    if p == "" {
        http.Error(w, "missing path", http.StatusBadRequest)
        return
    }
    abs, err := util.ResolvePath(p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Provide clearer 404s instead of default text from ServeFile.
    if st, err := os.Stat(abs); err != nil || st.IsDir() {
        if err != nil {
            http.Error(w, "not found: "+abs, http.StatusNotFound)
            return
        }
        http.Error(w, "path is a directory: "+abs, http.StatusBadRequest)
        return
    }

    // Force inline display when possible (images, pdf, etc.).
    w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(abs)+"\"")
    // CORS headers are set by middleware.WithCORS in the router.
    http.ServeFile(w, r, abs)
}
