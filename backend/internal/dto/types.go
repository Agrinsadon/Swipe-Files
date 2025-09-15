// Package dto defines transport data structures for JSON responses.
package dto

import "time"

// FileInfoDTO mirrors frontend expectations for file listing items.
type FileInfoDTO struct {
    Name    string    `json:"name"`
    Path    string    `json:"path"`
    Ext     string    `json:"ext"`
    Size    int64     `json:"size"`
    ModTime time.Time `json:"modTime"`
}
