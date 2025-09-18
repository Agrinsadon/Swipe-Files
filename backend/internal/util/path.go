// Package util: sisäiset apurit (ei HTTP-spesifejä).
package util

import (
    "os"
    "path/filepath"
    "strings"
)

// ResolvePath: laajentaa ~ ja palauttaa absoluuttisen polun. Suhteelliset ratkaistaan HOME:sta.
func ResolvePath(p string) (string, error) {
    home, _ := os.UserHomeDir()
    if p == "" {
        return home, nil
    }
    if strings.HasPrefix(p, "~") {
        if p == "~" {
            p = home
        } else if strings.HasPrefix(p, "~/") {
            p = filepath.Join(home, p[2:])
        }
    }
    if !filepath.IsAbs(p) {
        p = filepath.Join(home, p)
    }
    return filepath.Abs(p)
}
