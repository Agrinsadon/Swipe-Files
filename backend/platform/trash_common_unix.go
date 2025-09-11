package platform

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func uniqueDest(dir, name string) (string, error) {
	dst := filepath.Join(dir, name)
	if _, err := os.Stat(dst); err != nil {
		return dst, nil
	}
	ext := filepath.Ext(name)
	base := name[:len(name)-len(ext)]
	for i := 1; i < 1_000_000; i++ {
		alt := filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
		if _, err := os.Stat(alt); os.IsNotExist(err) {
			return alt, nil
		}
	}
	return "", fmt.Errorf("uniqueDest: too many duplicates")
}

func renameOrCopy(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Remove(src)
}
