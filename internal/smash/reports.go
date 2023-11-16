package smash

import (
	"github.com/thushan/smash/pkg/indexer"
	"path/filepath"
)

func resolveFilename(file indexer.FileFS) string {
	if file.Path == "." {
		return filepath.Base(file.FullName)
	} else {
		return file.Path
	}
}
