package handlers

import (
    "net/http"
    "os"
    "path/filepath"

    "Swipe-Files/backend/internal/util"
)

// Open: palvelee tiedoston tavut esikatselua varten.
// Kysely: ?path=... (absoluuttinen tai ~). Käyttää ResolvePathia, palvelee inline.
func Open(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        // CORS-esipyyntö; palautetaan nopeasti.
        w.WriteHeader(http.StatusNoContent)
        return
    }
    if r.Method != http.MethodGet {
        http.Error(w, "vain GET", http.StatusMethodNotAllowed)
        return
    }
    p := r.URL.Query().Get("path")
    if p == "" {
        http.Error(w, "path puuttuu", http.StatusBadRequest)
        return
    }
    abs, err := util.ResolvePath(p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Selkeämmät 404-viestit kuin ServeFile:n oletus.
    if st, err := os.Stat(abs); err != nil || st.IsDir() {
        if err != nil {
            http.Error(w, "ei löytynyt: "+abs, http.StatusNotFound)
            return
        }
        http.Error(w, "polku on hakemisto: "+abs, http.StatusBadRequest)
        return
    }

    // Pyydä selain näyttämään inline (kuvat, PDF, jne.).
    w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(abs)+"\"")
    // CORS-otsikot lisätään routerissa (middleware.WithCORS).
    http.ServeFile(w, r, abs)
}
