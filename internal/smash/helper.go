package smash

import (
	"path/filepath"

	"github.com/thushan/smash/pkg/indexer"
)

func resolveFilename(file *indexer.FileFS) string {
	if file.Path == "." {
		return filepath.Base(file.FullName)
	} else {
		return file.Path
	}
}
