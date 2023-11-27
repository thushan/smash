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

	excludeSysFileFilter []string
	excludeSysDirFilter  []string

	ExcludeDirFilter  []string
	ExcludeFileFilter []string

	IgnoreHiddenItems bool
}

func New() *IndexerConfig {
	return &IndexerConfig{
		IgnoreHiddenItems: true,
		ExcludeFileFilter: nil,
		ExcludeDirFilter:  nil,
		dirMatcher:        nil,
		fileMatcher:       nil,
		excludeSysDirFilter: []string{
			"System Volume Information", "$RECYCLE.BIN", "$MFT", /* Windows */
			".Trash", ".Trash-1000", /* Linux */
			".Trashes", /* macOS */
		},
		excludeSysFileFilter: []string{
			"thumbs.db", "desktop.ini", /* Windows */
			".ds_store", /* macOS */
		},
	}
}
func NewConfigured(excludeDirFilter []string, excludeFileFilter []string, ignoreHiddenItems bool) *IndexerConfig {
	indexer := New()
	if len(excludeFileFilter) > 0 {
		indexer.ExcludeFileFilter = excludeFileFilter
		indexer.fileMatcher = regexp.MustCompile(strings.Join(excludeFileFilter, "|"))
	}
	if len(excludeDirFilter) > 0 {
		indexer.ExcludeDirFilter = excludeDirFilter
		indexer.dirMatcher = regexp.MustCompile(strings.Join(excludeDirFilter, "|"))
	}
	indexer.IgnoreHiddenItems = ignoreHiddenItems
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

		isSystemObj := config.IgnoreHiddenItems && config.isHidden(name)

		if d.IsDir() {

			isIgnoreDir := config.isIgnored(name, config.excludeSysDirFilter)
			isExludeDir := len(config.ExcludeDirFilter) > 0 && config.dirMatcher.MatchString(path)

			if isSystemObj || isIgnoreDir || isExludeDir {
				return fs.SkipDir
			}

		} else {

			isIgnoreFile := config.isIgnored(name, config.excludeSysFileFilter)
			isExludeFile := len(config.ExcludeFileFilter) > 0 && config.fileMatcher.MatchString(name)

			if isSystemObj || isIgnoreFile || isExludeFile {
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

func (config *IndexerConfig) isIgnored(item string, collection []string) bool {
	for _, v := range collection {
		if strings.EqualFold(v, item) {
			return true
		}
	}
	return false
}
func (config *IndexerConfig) isHidden(name string) bool {
	return len(name) > 1 && name[0] == '.'
}
