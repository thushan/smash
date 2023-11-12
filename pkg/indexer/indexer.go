package indexer

import (
	"errors"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type FileFS struct {
	FileSystem fs.FS
	Path       string
	Name       string
	FullName   string
}
type IndexerConfig struct {
	dirMatcher  *regexp.Regexp
	fileMatcher *regexp.Regexp

	excludeSysFilter  []string
	ExcludeDirFilter  []string
	ExcludeFileFilter []string
}

func New() *IndexerConfig {
	return &IndexerConfig{
		ExcludeFileFilter: nil,
		ExcludeDirFilter:  nil,
		dirMatcher:        nil,
		fileMatcher:       nil,
		excludeSysFilter: []string{
			"System Volume Information", "$RECYCLE.BIN", "$MFT", /* Windows */
			".Trash", ".Trash-1000", /* Linux */
			".Trashes", /* macOS */
		},
	}
}
func NewConfigured(excludeDirFilter []string, excludeFileFilter []string) *IndexerConfig {
	indexer := New()
	if len(excludeFileFilter) > 0 {
		indexer.ExcludeFileFilter = excludeFileFilter
		indexer.fileMatcher = regexp.MustCompile(strings.Join(excludeFileFilter, "|"))
	}
	if len(excludeDirFilter) > 0 {
		indexer.ExcludeDirFilter = excludeDirFilter
		indexer.dirMatcher = regexp.MustCompile(strings.Join(excludeDirFilter, "|"))
	}
	return indexer
}

func (config *IndexerConfig) WalkDirectory(f fs.FS, root string, files chan FileFS) error {
	walkErr := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				return fs.SkipDir
			}
			return err
		}
		name := filepath.Clean(d.Name())
		if d.IsDir() {
			if config.isSystemFolder(name) || (len(config.ExcludeDirFilter) > 0 && config.dirMatcher.MatchString(path)) {
				return fs.SkipDir
			}
		} else {
			if len(config.ExcludeFileFilter) > 0 && config.fileMatcher.MatchString(name) {
				return nil
			}
			files <- FileFS{
				FileSystem: f,
				Path:       path,
				Name:       name,
				FullName:   filepath.Join(root, path),
			}
		}
		return nil
	})
	return walkErr
}

func (config *IndexerConfig) isSystemFolder(folder string) bool {
	for _, v := range config.excludeSysFilter {
		if folder == v {
			return true
		}
	}
	return false
}
