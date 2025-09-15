// Package handlers contains HTTP handler functions for API endpoints.
package handlers

import (
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "strconv"

    "Swipe-Files/backend/internal/dto"
    "Swipe-Files/backend/internal/util"
    "Swipe-Files/backend/internal/http/respond"
)

// ListFiles responds with a JSON array of file metadata for the requested
// directory (query param `dir`). It sorts results by modification time descending
// and optionally limits results via `?limit=N`.
func ListFiles(w http.ResponseWriter, r *http.Request) {
    dir := r.URL.Query().Get("dir")
    abs, err := util.ResolvePath(dir)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    entries, err := os.ReadDir(abs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    limit := 0
    if l := r.URL.Query().Get("limit"); l != "" {
        if n, err := strconv.Atoi(l); err == nil && n > 0 {
            limit = n
        }
    }

    var out []dto.FileInfoDTO
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        info, err := e.Info()
        if err != nil {
            continue
        }
        p := filepath.Join(abs, e.Name())
        out = append(out, dto.FileInfoDTO{
            Name:    e.Name(),
            Path:    p,
            Ext:     filepath.Ext(e.Name()),
            Size:    info.Size(),
            ModTime: info.ModTime(),
        })
    }
    sort.Slice(out, func(i, j int) bool { return out[i].ModTime.After(out[j].ModTime) })
    if limit > 0 && len(out) > limit {
        out = out[:limit]
    }
    respond.JSON(w, out, http.StatusOK)
}
