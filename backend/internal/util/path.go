// Package util contains internal helpers not specific to HTTP.
package util

import (
    "os"
    "path/filepath"
    "strings"
)

// ResolvePath expands ~ and returns an absolute path. Relative paths are resolved from HOME.
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
