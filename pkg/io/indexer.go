package io

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func IndexDirectory(fsys fs.FS, excludeDirFilter []string, excludeFileFilter []string, files chan string) {

	var dirMatcher = regexp.MustCompile(strings.Join(excludeDirFilter, "|"))
	var fileMatcher = regexp.MustCompile(strings.Join(excludeFileFilter, "|"))

	walkErr := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Index just the files
		if d.IsDir() {
			if len(excludeDirFilter) > 0 && dirMatcher.MatchString(path) {
				return filepath.SkipDir
			}
		} else {
			var filename = filepath.Base(path)
			if len(excludeFileFilter) > 0 && fileMatcher.MatchString(filename) {
				return nil
			}
			files <- path
		}

		return nil
	})
	if walkErr != nil {
		fmt.Fprintln(os.Stderr, walkErr)
	}
}
