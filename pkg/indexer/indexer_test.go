package indexer

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestIndexDirectoryWithFilesInRoot(t *testing.T) {
	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, nil, nil, true, t)

	expected := mockFiles
	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestIndexDirectoryWithFilesAcrossFolders(t *testing.T) {
	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
		"subfolder-1/DSC19845.ARW",
		"subfolder-1/DSC19846.ARW",
		"subfolder-2/DSC19847.ARW",
		"subfolder-2/DSC19848.ARW",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, nil, nil, true, t)

	expected := mockFiles
	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestIndexDirectoryWithDirExclusions(t *testing.T) {
	exclude_dir := []string{"subfolder-1", "subfolder-2", "subfolder-not-found"}
	exclude_file := []string{}

	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
		"subfolder-1/DSC19845.ARW",
		"subfolder-1/DSC19846.ARW",
		"subfolder-2/DSC19847.ARW",
		"subfolder-2/DSC19848.ARW",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, exclude_dir, exclude_file, true, t)

	expected := []string{
		mockFiles[0],
		mockFiles[1],
	}

	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestIndexDirectoryWithFileExclusions(t *testing.T) {
	exclude_dir := []string{}
	exclude_file := []string{"exclude.me"}

	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
		"exclude.me",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, exclude_dir, exclude_file, true, t)

	expected := []string{
		mockFiles[0],
		mockFiles[1],
	}

	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestIndexDirectoryWithFileAndDirExclusions(t *testing.T) {

	exclude_dir := []string{"exclude-dir"}
	exclude_file := []string{"exclude.me"}

	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
		"exclude.me",
		"exclude-dir/random.file",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, exclude_dir, exclude_file, true, t)

	expected := []string{
		mockFiles[0],
		mockFiles[1],
	}

	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestIndexDirectoryWithHiddenFilesThatShouldBeIndexed(t *testing.T) {
	exclude_dir := []string{}
	exclude_file := []string{}

	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
		".tmux",
		".config/smash/config.json",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, exclude_dir, exclude_file, false, t)

	expected := []string{
		mockFiles[3],
		mockFiles[2],
		mockFiles[0],
		mockFiles[1],
	}

	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}

func TestIndexDirectoryWithHiddenFiles(t *testing.T) {

	exclude_dir := []string{"exclude-dir"}
	exclude_file := []string{"exclude.me"}

	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
		".tmux",
		".config/smash/config.json",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, exclude_dir, exclude_file, true, t)

	expected := []string{
		mockFiles[0],
		mockFiles[1],
	}

	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}
func TestIndexDirectoryWhichContainsSystemFiles(t *testing.T) {
	exclude_dir := []string{}
	exclude_file := []string{}

	mockFiles := []string{
		"DSC19841.ARW",
		"THUMBS.DB",
		"desktop.ini",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, exclude_dir, exclude_file, true, t)

	expected := []string{
		mockFiles[0],
	}

	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}
func TestIndexDirectoryWhichContainsWindowsSystemFiles(t *testing.T) {
	exclude_dir := []string{}
	exclude_file := []string{}

	mockFiles := []string{
		"DSC19841.ARW",
		"DSC19842.ARW",
		"$RECYCLE.BIN/test.txt",
		"$MFT/random.file",
	}

	walkedFiles := walkDirectoryTestRunner(mockFiles, exclude_dir, exclude_file, true, t)

	expected := []string{
		mockFiles[0],
		mockFiles[1],
	}

	actual := walkedFiles

	if len(actual) != len(expected) {
		t.Errorf("expected %d, got %d files", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v files", expected, actual)
	}
}
func channelFileToSliceOfFiles(ch <-chan FileFS) []string {
	var result []string
	for f := range ch {
		result = append(result, f.Path)
	}
	return result
}

func walkDirectoryTestRunner(files []string, excludeDir []string, excludeFiles []string, ignoreHiddenItems bool, t *testing.T) []string {
	fr := "mock://"
	fs := createMockFS(files)
	ch := make(chan FileFS)

	go func() {
		defer close(ch)
		indexer := NewConfigured(excludeDir, excludeFiles, ignoreHiddenItems)
		err := indexer.WalkDirectory(fs, fr, ch)
		if err != nil {
			t.Errorf("WalkDirectory returned an error: %v", err)
		}
	}()

	return channelFileToSliceOfFiles(ch)
}
func createMockFS(files []string) fstest.MapFS {
	var fs fstest.MapFS = make(map[string]*fstest.MapFile)
	for _, file := range files {
		fs[file] = &fstest.MapFile{}
	}
	return fs
}
