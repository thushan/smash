// nolint
package slicer

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"github.com/thushan/smash/internal/algorithms"
	"io"
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSlice_New_OffsetMapWith1MbBlob(t *testing.T) {
	fsSize := 1024000 // 1Mb
	binary := randomBytes(fsSize)
	reader := bytes.NewReader(binary)

	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}

	sr := io.NewSectionReader(reader, 0, int64(fsSize))

	stats := SlicerStats{}

	slicer := New(algorithms.Xxhash)

	if err := slicer.Slice(sr, &options, &stats); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	}
	// For a 1024000 byte blob with 4 segments
	expected := make(map[int]int64)
	expected[0] = 0 // head
	expected[1] = 251904
	expected[2] = 503808
	expected[3] = 755712
	expected[4] = 1007616
	expected[5] = 1015808 // tail
	actual := stats.SliceOffsets

	if len(expected) != len(actual) {
		t.Errorf("offset total expected %d, got %d", len(expected), len(actual))
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v offset boundaries", expected, actual)
	}
}

func TestSlice_New_With0KbBlob(t *testing.T) {
	binary := []byte{}
	reader := bytes.NewReader(binary)

	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}

	sr := io.NewSectionReader(reader, 0, 0)

	stats := SlicerStats{}

	slicer := New(algorithms.Xxhash)

	if err := slicer.Slice(sr, &options, &stats); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	}

	if stats.EmptyFile != true {
		t.Errorf("expected Empty to be %v, got %v", true, stats.EmptyFile)
	}

	if stats.FileSize != 0 {
		t.Errorf("expected FileSize to be %d, got %d", 0, stats.FileSize)
	}

	if stats.HashedFullFile != false {
		t.Errorf("expected HashedFullFile to be %v, got %v", false, stats.HashedFullFile)
	}
}
func TestSlice_New_NoOffsetMapWith1KbBlob(t *testing.T) {
	fsSize := 1024 // 1Kb
	binary := randomBytes(fsSize)
	reader := bytes.NewReader(binary)

	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}

	sr := io.NewSectionReader(reader, 0, int64(fsSize))

	stats := SlicerStats{}

	slicer := New(algorithms.Xxhash)

	if err := slicer.Slice(sr, &options, &stats); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	}

	actual := stats.SliceOffsets

	if actual != nil {
		t.Errorf("offset not expected for small file, got %v", len(actual))
	}
}
func TestSlice_New_WithTextBinaryBlob(t *testing.T) {
	binary := []byte("OMG THIS IS TEXT!")
	fsSize := len(binary)
	reader := bytes.NewReader(binary)

	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}

	sr := io.NewSectionReader(reader, 0, int64(fsSize))

	stats := SlicerStats{}

	slicer := New(algorithms.Xxhash)

	if err := slicer.Slice(sr, &options, &stats); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	}

	expected := "67938b74b221486b"
	actual := hex.EncodeToString(stats.Hash)

	if !strings.EqualFold(actual, expected) {
		t.Errorf("expected hash %s, got %s", expected, actual)
	}
}
func TestSlice_New_WithTextBlob(t *testing.T) {
	stexty := "OMG THIS IS TEXT!"
	fsSize := len(stexty)
	reader := strings.NewReader(stexty)

	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}

	sr := io.NewSectionReader(reader, 0, int64(fsSize))

	stats := SlicerStats{}

	slicer := New(algorithms.Xxhash)

	if err := slicer.Slice(sr, &options, &stats); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	}

	expected := "67938b74b221486b"
	actual := hex.EncodeToString(stats.Hash)

	if !strings.EqualFold(actual, expected) {
		t.Errorf("expected hash %s, got %s", expected, actual)
	}
}
func TestSliceFS_New_FileSystemTestFile_TestWordPdf_WithSlicing(t *testing.T) {
	fsys := os.DirFS("./artefacts")
	filename := "test.pdf"
	algorithm := algorithms.Xxhash
	expected := "bedd0999e968547e"
	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile_WithSliceFS(fsys, filename, algorithm, &options, expected, t)
}
func TestSliceFS_New_FileSystemTestFile_TestAdobePdf_WithSlicing(t *testing.T) {
	fsys := os.DirFS("./artefacts")
	filename := "test-adobe.pdf"
	algorithm := algorithms.Xxhash
	expected := "41d4d0a83d10e962"
	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile_WithSliceFS(fsys, filename, algorithm, &options, expected, t)
}
func TestSliceFS_New_FileSystemTestFile_Test1mb_WithSlicing(t *testing.T) {
	fsys := os.DirFS("./artefacts")
	filename := "test.1mb"
	algorithm := algorithms.Xxhash
	expected := "bb83f43630ee546f"
	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile_WithSliceFS(fsys, filename, algorithm, &options, expected, t)
}
func TestSliceFS_New_FileSystemTestFile_Test1mb_WithoutMeta(t *testing.T) {
	fsys := os.DirFS("./artefacts")
	filename := "test.1mb"
	algorithm := algorithms.Xxhash
	expected := "913c30543271faaf"
	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          true,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile_WithSliceFS(fsys, filename, algorithm, &options, expected, t)
}
func TestSliceFS_New_FileSystemTestFile_TestManipulated1mb_WithSlicing(t *testing.T) {
	fsys := os.DirFS("./artefacts")
	filename := "test-manipulated.1mb"
	algorithm := algorithms.Xxhash
	expected := "4f595576799edcd9"
	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile_WithSliceFS(fsys, filename, algorithm, &options, expected, t)
}
func TestSlice_New_FileSystemTestFile_Test1mb_WithSlicing(t *testing.T) {
	algorithm := algorithms.Xxhash
	expected := "bb83f43630ee546f"
	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile("./artefacts/test.1mb", algorithm, &options, expected, t)
}
func TestSlice_New_FileSystemTestFile_TestManipulated1mb_WithSlicing(t *testing.T) {
	algorithm := algorithms.Xxhash
	expected := "4f595576799edcd9"
	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile("./artefacts/test-manipulated.1mb", algorithm, &options, expected, t)
}
func TestSlice_New_FileSystemTestFile_Test1mb_WithoutSlicing(t *testing.T) {
	algorithm := algorithms.Xxhash
	expected := "6b6255ee515dcc04"
	options := SlicerOptions{
		DisableSlicing:       true,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile("./artefacts/test.1mb", algorithm, &options, expected, t)
}
func TestSlice_New_FileSystemTestFile_TestManipulated1mb_WithoutSlicing(t *testing.T) {
	algorithm := algorithms.Xxhash
	expected := "4a1960f16a88960c"
	options := SlicerOptions{
		DisableSlicing:       true,
		DisableMeta:          false,
		DisableFileDetection: false,
	}
	runHashCheckTestsForFileSystemFile("./artefacts/test-manipulated.1mb", algorithm, &options, expected, t)
}
func runHashCheckTestsForFileSystemFile_WithSliceFS(fs fs.FS, filename string, algorithm algorithms.Algorithm, options *SlicerOptions, expected string, t *testing.T) {

	slicer := New(algorithm)

	if stats, err := slicer.SliceFS(fs, filename, options); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	} else {

		actual := hex.EncodeToString(stats.Hash)

		if len(expected) != len(actual) {
			t.Errorf("hash length expected %d, got %d", len(expected), len(actual))
		}

		if !strings.EqualFold(actual, expected) {
			t.Errorf("expected hash %s, got %s for %s", expected, actual, filename)
		}
	}

}
func runHashCheckTestsForFileSystemFile(filename string, algorithm algorithms.Algorithm, options *SlicerOptions, expected string, t *testing.T) {
	if binary, err := os.ReadFile(filename); err != nil {
		t.Errorf("Unexpected io error %v", err)
	} else {

		fsSize := len(binary)
		reader := bytes.NewReader(binary)
		sr := io.NewSectionReader(reader, 0, int64(fsSize))

		stats := SlicerStats{}

		slicer := New(algorithm)

		if err := slicer.Slice(sr, options, &stats); err != nil {
			t.Errorf("Unexpected Slicer error %v", err)
		}

		actual := hex.EncodeToString(stats.Hash)

		if len(expected) != len(actual) {
			t.Errorf("hash length expected %d, got %d", len(expected), len(actual))
		}

		if !strings.EqualFold(actual, expected) {
			t.Errorf("expected hash %s, got %s", expected, actual)
		}
	}
}
func TestSlice_New_Hash_xxHash_With1KbBlob(t *testing.T) {
	runHashAlgorithmTest(algorithms.Xxhash, t)
}
func TestSlice_New_Hash_Fnv128_With1KbBlob(t *testing.T) {
	runHashAlgorithmTest(algorithms.Fnv128, t)
}
func TestSlice_New_Hash_Fnv128a_With1KbBlob(t *testing.T) {
	runHashAlgorithmTest(algorithms.Fnv128a, t)
}

func runHashAlgorithmTest(algorithm algorithms.Algorithm, t *testing.T) {
	fsSize := 1024 // 1kb
	binary := randomBytes(fsSize)
	reader := bytes.NewReader(binary)

	sr := io.NewSectionReader(reader, 0, int64(fsSize))

	options := SlicerOptions{
		DisableSlicing:       false,
		DisableMeta:          false,
		DisableFileDetection: false,
	}

	stats := SlicerStats{}

	slicer := New(algorithm)

	if err := slicer.Slice(sr, &options, &stats); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	}

	hasher := algorithm.New()
	hasher.Write(binary)

	expected := hasher.Sum(nil)
	actual := stats.Hash

	if stats.HashedFullFile != true {
		t.Errorf("expected full hashing of file but wasn't")
	}

	if len(expected) != len(actual) {
		t.Errorf("hash size expected %d, got %d", len(expected), len(actual))
	}

	if !bytes.Equal(actual, expected) {
		t.Errorf("expected %x, got %x hash", expected, actual)
	}
}
func randomBytes(length int) []byte {
	buffer := make([]byte, length)
	_, _ = rand.Read(buffer)
	return buffer
}
