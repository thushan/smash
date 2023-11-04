package io

import (
	"crypto/rand"
	"testing"
	"testing/fstest"
)

func randomBytes(length int) []byte {
	buffer := make([]byte, length)
	_, _ = rand.Read(buffer)
	return buffer
}
func TestIndexDirectoryWithFilesInRoot(t *testing.T) {
	fsq := make(chan string, 10)
	exclude_dir := []string{}
	exclude_file := []string{}

	fs := fstest.MapFS{
		"DSC19841.ARW": {Data: randomBytes(1024)},
		"DSC19842.ARW": {Data: randomBytes(2048)},
	}

	IndexDirectory(fs, exclude_dir, exclude_file, fsq)

	expected := len(fs)
	actual := len(fsq)

	if actual != expected {
		t.Errorf("expected %d, got %d files", expected, actual)
	}
}

func TestIndexDirectoryWithFilesAcrossFolders(t *testing.T) {
	fsq := make(chan string, 10)
	exclude_dir := []string{}
	exclude_file := []string{}

	fs := fstest.MapFS{
		"DSC19841.ARW":             {Data: randomBytes(1024)},
		"DSC19842.ARW":             {Data: randomBytes(2048)},
		"subfolder-1/DSC19845.ARW": {Data: randomBytes(1024)},
		"subfolder-1/DSC19846.ARW": {Data: randomBytes(1024)},
		"subfolder-2/DSC19847.ARW": {Data: randomBytes(1024)},
		"subfolder-2/DSC19848.ARW": {Data: randomBytes(1024)},
	}

	IndexDirectory(fs, exclude_dir, exclude_file, fsq)

	expected := len(fs)
	actual := len(fsq)

	if actual != expected {
		t.Errorf("expected %d, got %d files", expected, actual)
	}
}

func TestIndexDirectoryWithDirExclusions(t *testing.T) {
	fsq := make(chan string, 10)
	exclude_dir := []string{"subfolder-1", "subfolder-2", "subfolder-not-found"}
	exclude_file := []string{}

	fs := fstest.MapFS{
		"DSC19841.ARW":             {Data: randomBytes(1024)},
		"DSC19842.ARW":             {Data: randomBytes(2048)},
		"subfolder-1/DSC19845.ARW": {Data: randomBytes(1024)},
		"subfolder-1/DSC19846.ARW": {Data: randomBytes(1024)},
		"subfolder-2/DSC19847.ARW": {Data: randomBytes(1024)},
		"subfolder-2/DSC19848.ARW": {Data: randomBytes(1024)},
	}

	IndexDirectory(fs, exclude_dir, exclude_file, fsq)

	expected := len(fs) - 4
	actual := len(fsq)

	if actual != expected {
		t.Errorf("expected %d, got %d files", expected, actual)
	}
}

func TestIndexDirectoryWithFileExclusions(t *testing.T) {
	fsq := make(chan string, 10)
	exclude_dir := []string{}
	exclude_file := []string{"exclude.me"}

	fs := fstest.MapFS{
		"DSC19841.ARW": {Data: randomBytes(1024)},
		"DSC19842.ARW": {Data: randomBytes(2048)},
		"exclude.me":   {Data: randomBytes(1024)},
	}

	IndexDirectory(fs, exclude_dir, exclude_file, fsq)

	expected := len(fs) - 1
	actual := len(fsq)

	if actual != expected {
		t.Errorf("expected %d, got %d files", expected, actual)
	}
}

func TestIndexDirectoryWithFileAndDirExclusions(t *testing.T) {
	fsq := make(chan string, 10)
	exclude_dir := []string{"exclude-dir"}
	exclude_file := []string{"exclude.me"}

	fs := fstest.MapFS{
		"DSC19841.ARW":            {Data: randomBytes(1024)},
		"DSC19842.ARW":            {Data: randomBytes(2048)},
		"exclude.me":              {Data: randomBytes(1024)},
		"exclude-dir/random.file": {Data: randomBytes(1024)},
	}

	IndexDirectory(fs, exclude_dir, exclude_file, fsq)

	expected := len(fs) - 2
	actual := len(fsq)

	if actual != expected {
		t.Errorf("expected %d, got %d files", expected, actual)
	}
}
