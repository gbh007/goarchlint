package render

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

func (r Render) CleanDocFolder() error {
	var files, dirs []string

	// Не используется os.RemoveAll специально, для последующих расширений
	err := filepath.WalkDir(r.BasePath, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}

		if d.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk docs dir: %w", err)
	}

	slices.Reverse(dirs)

	for _, s := range files {
		err := os.Remove(s)
		if err != nil {
			return fmt.Errorf("remove file: %w", err)
		}
	}

	for _, s := range dirs {
		err := os.Remove(s)
		if err != nil {
			return fmt.Errorf("remove dir: %w", err)
		}
	}

	return nil
}
