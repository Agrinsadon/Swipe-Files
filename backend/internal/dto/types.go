// Package dto: JSON-siirtorakenteet frontendin ja backendin välillä.
package dto

import "time"

// FileInfoDTO: vastaa frontendin odottamaa tiedostolistan rakennetta.
type FileInfoDTO struct {
    Name    string    `json:"name"`
    Path    string    `json:"path"`
    Ext     string    `json:"ext"`
    Size    int64     `json:"size"`
    ModTime time.Time `json:"modTime"`
}
