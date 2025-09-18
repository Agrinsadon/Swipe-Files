package handlers

import (
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"

    "Swipe-Files/backend/internal/dto"
    "Swipe-Files/backend/internal/http/respond"
    "Swipe-Files/backend/internal/util"
)

// defaultRoots: kotikansion tyypilliset juuret, joista etsitään "Äskeiset".
var defaultRoots = []string{"~/Downloads", "~/Desktop", "~/Documents", "~/Pictures"}

// Recents: listaa viimeksi muokattuja tiedostoja useista juurista, rekursiivisesti
// pieneen syvyyteen saakka. Kyselyparametrit:
//   dirs     = pilkulla eroteltu lista (oletus: defaultRoots)
//   limit    = max palautettavien tiedostojen määrä (oletus 200)
//   maxDepth = rekursion syvyys per juuri (oletus 2)
func Recents(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()

    // Juurikansiot
    rootsParam := strings.TrimSpace(q.Get("dirs"))
    roots := defaultRoots
    if rootsParam != "" {
        parts := strings.Split(rootsParam, ",")
        roots = nil
        for _, p := range parts {
            p = strings.TrimSpace(p)
            if p != "" {
                roots = append(roots, p)
            }
        }
        if len(roots) == 0 {
            roots = defaultRoots
        }
    }

    // Limit ja maxDepth
    limit := 200
    if v := q.Get("limit"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            limit = n
        }
    }
    maxDepth := 2
    if v := q.Get("maxDepth"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n >= 0 {
            maxDepth = n
        }
    }

    // Käydään läpi juuret, kerätään tiedostot, järjestetään uusin ensin.
    seen := make(map[string]struct{})
    // Ensure empty JSON array ([]) instead of null when no files
    out := make([]dto.FileInfoDTO, 0)

    for _, root := range roots {
        abs, err := util.ResolvePath(root)
        if err != nil {
            continue
        }
        // Syvyysseuranta: lasketaan rootin polun segmentit ja verrataan lapsiin.
        rootSegments := segments(abs)
        // WalkDir on tehokas; rajataan syvyys ja määrä maltilliseksi.
        _ = filepath.WalkDir(abs, func(p string, d os.DirEntry, err error) error {
            if err != nil {
                return nil
            }
            // Syvyyden rajoitus
            if depth(abs, p, rootSegments) > maxDepth {
                if d.IsDir() {
                    return filepath.SkipDir
                }
                return nil
            }
            if d.IsDir() {
                return nil
            }
            if _, ok := seen[p]; ok {
                return nil
            }
            info, err := d.Info()
            if err != nil {
                return nil
            }
            seen[p] = struct{}{}
            out = append(out, dto.FileInfoDTO{
                Name:    d.Name(),
                Path:    p,
                Ext:     filepath.Ext(d.Name()),
                Size:    info.Size(),
                ModTime: info.ModTime(),
            })
            return nil
        })
    }

    sort.Slice(out, func(i, j int) bool { return out[i].ModTime.After(out[j].ModTime) })
    if limit > 0 && len(out) > limit {
        out = out[:limit]
    }
    respond.JSON(w, out, http.StatusOK)
}

// segments: laskee polun segmenttien määrän.
func segments(p string) int { return len(splitPathClean(p)) }

// depth: palauttaa kuinka syvällä child on rootista laskettuna.
func depth(root, child string, rootSegs int) int {
    return len(splitPathClean(child)) - rootSegs
}

func splitPathClean(p string) []string {
    cp := filepath.Clean(p)
    if cp == string(filepath.Separator) {
        return []string{""}
    }
    parts := []string{}
    for cp != string(filepath.Separator) {
        dir, base := filepath.Split(cp)
        if base != "" {
            parts = append([]string{base}, parts...)
        }
        if dir == cp {
            break
        }
        cp = strings.TrimRight(dir, string(filepath.Separator))
        if cp == "" {
            break
        }
    }
    // lisää juuri
    if strings.HasPrefix(p, string(filepath.Separator)) {
        parts = append([]string{""}, parts...)
    }
    return parts
}
