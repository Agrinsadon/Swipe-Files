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

// Convert: muuntaa toimistoasiakirjan PDF:ksi LibreOfficella (soffice), jos saatavilla.
// Vain to=pdf tuettu. Kysely: ?path=...&to=pdf. Palauttaa 501 jos muunnostyökalu puuttuu.
func Convert(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }
    if r.Method != http.MethodGet {
        http.Error(w, "vain GET", http.StatusMethodNotAllowed)
        return
    }

    to := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("to")))
    if to == "" {
        to = "pdf"
    }
    if to != "pdf" {
        http.Error(w, "kohdeformaattia ei tueta", http.StatusBadRequest)
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
    st, err := os.Stat(abs)
    if err != nil || st.IsDir() {
        if err != nil {
            http.Error(w, "ei löytynyt: "+abs, http.StatusNotFound)
            return
        }
        http.Error(w, "polku on hakemisto: "+abs, http.StatusBadRequest)
        return
    }

    // Etsi soffice/libreoffice-binääri
    bin, err := exec.LookPath("soffice")
    if err != nil {
        if b2, err2 := exec.LookPath("libreoffice"); err2 == nil {
            bin = b2
        } else {
            // Kokeile tyypillisiä asennuspolkuja
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
                // Palauta selkeä 501-virhe ylläpidolle
                msg := "palvelinmuunnos ei käytettävissä: asenna LibreOffice (soffice) palvelimelle"
                if runtime.GOOS == "darwin" {
                    msg += "; esim. brew install --cask libreoffice"
                } else if runtime.GOOS == "linux" {
                    msg += "; esim. sudo apt-get install -y libreoffice"
                } else if runtime.GOOS == "windows" {
                    msg += "; lataa: https://www.libreoffice.org/download/"
                }
                respond.JSON(w, map[string]string{"error": msg}, http.StatusNotImplemented)
                return
            }
        }
    }

    // Välimuistiavain polusta ja muokkausajasta
    h := sha1.New()
    _, _ = h.Write([]byte(abs))
    _, _ = h.Write([]byte(st.ModTime().UTC().Format(time.RFC3339Nano)))
    key := hex.EncodeToString(h.Sum(nil))
    outDir := filepath.Join(os.TempDir(), "swipe-files-cache")
    _ = os.MkdirAll(outDir, 0o755)
    outPath := filepath.Join(outDir, key+".pdf")

    if _, err := os.Stat(outPath); err != nil {
        // Ei välimuistissa; muunna
        cmd := exec.Command(bin, "--headless", "--convert-to", "pdf", "--outdir", outDir, abs)
        // Joissain ympäristöissä HOME tarvitaan
        cmd.Env = append(os.Environ(), "HOME="+os.TempDir())
        if out, err := cmd.CombinedOutput(); err != nil {
            http.Error(w, fmt.Sprintf("muunnos epäonnistui: %v: %s", err, string(out)), http.StatusInternalServerError)
            return
        }
        // LibreOffice nimeää tulosteen <nimi>.pdf — uudelleennimeä välimuistipolkuun
        produced := filepath.Join(outDir, strings.TrimSuffix(filepath.Base(abs), filepath.Ext(abs))+".pdf")
        if produced != outPath {
            _ = os.Rename(produced, outPath)
        }
    }

    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(abs)+".pdf\"")
    http.ServeFile(w, r, outPath)
}
