package indexer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type IndexerConfig struct {
	dirMatcher  *regexp.Regexp
	fileMatcher *regexp.Regexp

	ExcludeDirFilter  []string
	ExcludeFileFilter []string
}

func New() *IndexerConfig {
	return &IndexerConfig{
		ExcludeFileFilter: nil,
		ExcludeDirFilter:  nil,
		dirMatcher:        nil,
		fileMatcher:       nil,
	}
}

func NewConfigured(excludeDirFilter []string, excludeFileFilter []string) *IndexerConfig {
	return &IndexerConfig{
		ExcludeDirFilter:  excludeDirFilter,
		ExcludeFileFilter: excludeFileFilter,
		dirMatcher:        regexp.MustCompile(strings.Join(excludeDirFilter, "|")),
		fileMatcher:       regexp.MustCompile(strings.Join(excludeFileFilter, "|")),
	}
}

func (config *IndexerConfig) WalkDirectory(fsys fs.FS, files chan string) {
	walkErr := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Index just the files
		if d.IsDir() {
			if isSystemFolder(d.Name()) || (len(config.ExcludeDirFilter) > 0 && config.dirMatcher.MatchString(path)) {
				return filepath.SkipDir
			}
		} else {
			filename := filepath.Base(path)
			if len(config.ExcludeFileFilter) > 0 && config.fileMatcher.MatchString(filename) {
				return nil
			}
			files <- path
		}

		return nil
	})
	if walkErr != nil {
		fmt.Fprintln(os.Stderr, "Walk Failed: ", walkErr)
	}
}

func isSystemFolder(path string) bool {
	folder := filepath.Clean(path)
	skipDirs := []string{
		"System Volume Information", "$RECYCLE.BIN", "$MFT", /* Windows */
		".Trash", ".Trash-1000", /* Linux */
		".Trashes", /* macOS */
	}
	for _, v := range skipDirs {
		if folder == v {
			return true
		}
	}
	return false
}
