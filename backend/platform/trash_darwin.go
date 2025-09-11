package platform

import (
	"os"
	"path/filepath"
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
	dstDir := filepath.Join(home, ".Trash")
	if err := os.MkdirAll(dstDir, 0o700); err != nil {
		return err
	}
	dst, err := uniqueDest(dstDir, filepath.Base(abs))
	if err != nil {
		return err
	}
	return renameOrCopy(abs, dst)
}
