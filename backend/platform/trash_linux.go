package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func MoveToTrash(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); err != nil {
		return err
	}
	home, _ := os.UserHomeDir()
	root := filepath.Join(home, ".local", "share", "Trash")
	filesDir := filepath.Join(root, "files")
	infoDir := filepath.Join(root, "info")
	if err := os.MkdirAll(filesDir, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(infoDir, 0o700); err != nil {
		return err
	}
	dst, err := uniqueDest(filesDir, filepath.Base(abs))
	if err != nil {
		return err
	}
	if err := renameOrCopy(abs, dst); err != nil {
		return err
	}
	infoPath := filepath.Join(infoDir, filepath.Base(dst)+".trashinfo")
	info := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		abs, time.Now().Format("2006-01-02T15:04:05"))
	return os.WriteFile(infoPath, []byte(info), 0o600)
}
