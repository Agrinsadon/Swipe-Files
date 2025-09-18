package handlers

import (
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "time"

    "Swipe-Files/backend/internal/http/respond"
    "Swipe-Files/backend/internal/util"
)

// Convert converts office-like documents to PDF using LibreOffice (soffice) if available.
// Only supports to=pdf. Query: ?path=...&to=pdf
// Returns 501 if conversion tool not present.
func Convert(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }
    if r.Method != http.MethodGet {
        http.Error(w, "GET only", http.StatusMethodNotAllowed)
        return
    }

    to := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("to")))
    if to == "" {
        to = "pdf"
    }
    if to != "pdf" {
        http.Error(w, "unsupported target format", http.StatusBadRequest)
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
    st, err := os.Stat(abs)
    if err != nil || st.IsDir() {
        if err != nil {
            http.Error(w, "not found: "+abs, http.StatusNotFound)
            return
        }
        http.Error(w, "path is a directory: "+abs, http.StatusBadRequest)
        return
    }

    // Check for soffice/libreoffice
    bin, err := exec.LookPath("soffice")
    if err != nil {
        if b2, err2 := exec.LookPath("libreoffice"); err2 == nil {
            bin = b2
        } else {
            // Try common install locations
            candidates := []string{
                "/opt/homebrew/bin/soffice",
                "/usr/local/bin/soffice",
                "/usr/bin/soffice",
                "/snap/bin/libreoffice",
                "/usr/lib/libreoffice/program/soffice",
            }
            if runtime.GOOS == "darwin" {
                candidates = append([]string{"/Applications/LibreOffice.app/Contents/MacOS/soffice"}, candidates...)
            }
            if runtime.GOOS == "windows" {
                pf := os.Getenv("ProgramFiles")
                pfx := os.Getenv("ProgramFiles(x86)")
                if pf != "" {
                    candidates = append(candidates, filepath.Join(pf, "LibreOffice", "program", "soffice.exe"))
                }
                if pfx != "" {
                    candidates = append(candidates, filepath.Join(pfx, "LibreOffice", "program", "soffice.exe"))
                }
            }
            for _, c := range candidates {
                if st, err := os.Stat(c); err == nil && !st.IsDir() {
                    bin = c
                    break
                }
            }
            if bin == "" {
                // Return a helpful 501 JSON error for admins
                msg := "server conversion not available: install LibreOffice (soffice) on the server"
                if runtime.GOOS == "darwin" {
                    msg += "; e.g. brew install --cask libreoffice"
                } else if runtime.GOOS == "linux" {
                    msg += "; e.g. sudo apt-get install -y libreoffice"
                } else if runtime.GOOS == "windows" {
                    msg += "; download from https://www.libreoffice.org/download/"
                }
                respond.JSON(w, map[string]string{"error": msg}, http.StatusNotImplemented)
                return
            }
        }
    }

    // Cache key based on path and modtime
    h := sha1.New()
    _, _ = h.Write([]byte(abs))
    _, _ = h.Write([]byte(st.ModTime().UTC().Format(time.RFC3339Nano)))
    key := hex.EncodeToString(h.Sum(nil))
    outDir := filepath.Join(os.TempDir(), "swipe-files-cache")
    _ = os.MkdirAll(outDir, 0o755)
    outPath := filepath.Join(outDir, key+".pdf")

    if _, err := os.Stat(outPath); err != nil {
        // Not cached; convert
        cmd := exec.Command(bin, "--headless", "--convert-to", "pdf", "--outdir", outDir, abs)
        // Some environments require HOME set
        cmd.Env = append(os.Environ(), "HOME="+os.TempDir())
        if out, err := cmd.CombinedOutput(); err != nil {
            http.Error(w, fmt.Sprintf("conversion failed: %v: %s", err, string(out)), http.StatusInternalServerError)
            return
        }
        // LibreOffice names output as <basename>.pdf — rename to our cache path
        produced := filepath.Join(outDir, strings.TrimSuffix(filepath.Base(abs), filepath.Ext(abs))+".pdf")
        if produced != outPath {
            _ = os.Rename(produced, outPath)
        }
    }

    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(abs)+".pdf\"")
    http.ServeFile(w, r, outPath)
}
