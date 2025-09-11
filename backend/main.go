package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"../backend/platform"
)

type FileInfoDTO struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Ext     string    `json:"ext"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/files", withCORS(listFilesHandler))
	mux.HandleFunc("/api/trash", withCORS(trashHandler))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":8787"
	log.Printf("Backend up at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func resolvePath(p string) (string, error) {
	home, _ := os.UserHomeDir()
	if p == "" {
		return home, nil
	}
	// ~ ja ~/...
	if strings.HasPrefix(p, "~") {
		if p == "~" {
			p = home
		} else if strings.HasPrefix(p, "~/") {
			p = filepath.Join(home, p[2:])
		}
	}
	// jos edelleen suhteellinen → tulkitaan HOME-relativeksi
	if !filepath.IsAbs(p) {
		p = filepath.Join(home, p)
	}
	return filepath.Abs(p)
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	abs, err := resolvePath(dir)
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

	var out []FileInfoDTO
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		p := filepath.Join(abs, e.Name())
		out = append(out, FileInfoDTO{
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
	writeJSON(w, out, http.StatusOK)
}

func trashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		http.Error(w, "invalid JSON (need {\"path\": \"...\"})", http.StatusBadRequest)
		return
	}
	abs, err := resolvePath(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := platform.MoveToTrash(abs); err != nil {
		http.Error(w, "trash failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok", "path": abs}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
