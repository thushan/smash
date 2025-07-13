// nolint
package slicer

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/thushan/smash/internal/algorithms"
)

func TestSlice_New_OffsetMapWith1MbBlob(t *testing.T) {
	fsSize := 1024000 // 1Mb
	binary := randomBytes(fsSize)
	reader := bytes.NewReader(binary)

	options := Options{
		DisableSlicing:  false,
		DisableMeta:     false,
		DisableAutoText: false,
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

	options := Options{
		DisableSlicing:  false,
		DisableMeta:     false,
		DisableAutoText: false,
	}

	sr := io.NewSectionReader(reader, 0, 0)

	stats := SlicerStats{}

	slicer := New(algorithms.Xxhash)

	if err := slicer.Slice(sr, &options, &stats); err != nil {
		t.Errorf("Unexpected Slicer error %v", err)
	}

	if stats.EmptyFile != true {
		t.Errorf("expected Files to be %v, got %v", true, stats.EmptyFile)
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

	options := Options{
		DisableSlicing:  false,
		DisableMeta:     false,
		DisableAutoText: false,
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

	options := Options{
		DisableSlicing:  false,
		DisableMeta:     false,
		DisableAutoText: false,
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
func TestSlice_Slice_CheckSizeThresholds(t *testing.T) {
	stexty := "OMG THIS IS TEXT!"
	fsSize := len(stexty)
	reader := strings.NewReader(stexty)
	sr := io.NewSectionReader(reader, 0, int64(fsSize))
	slicer := New(algorithms.Xxhash)

	tests := []struct {
		name     string
		options  Options
		expected bool
	}{
		{
			name: "ShouldIgnoreWhenSizeBelowMinSize",
			options: Options{
				MinSize:         uint64(fsSize + 10),
				MaxSize:         DefaultMaxSize,
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: true,
		},
		{
			name: "ShouldIgnoreWhenSizeAboveMaxSize",
			options: Options{
				MinSize:         DefaultMinSize,
				MaxSize:         uint64(fsSize - 10),
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: true,
		},
		{
			name: "ShouldNotIgnoreWhenMinSizeWithinRange",
			options: Options{
				MinSize:         uint64(fsSize - 10),
				MaxSize:         DefaultMaxSize,
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: false,
		},
		{
			name: "ShouldNotIgnoreWhenMaxSizeWithinRange",
			options: Options{
				MinSize:         DefaultMinSize,
				MaxSize:         uint64(fsSize + 10),
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: false,
		},
		{
			name: "ShouldNotIgnoreWhenDefaultSize",
			options: Options{
				MinSize:         DefaultMinSize,
				MaxSize:         DefaultMaxSize,
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: false,
		},
		{
			name: "ShouldNotIgnoreWhenSizeWithinMinMaxThreshold",
			options: Options{
				MinSize:         uint64(fsSize - 10),
				MaxSize:         uint64(fsSize + 10),
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := SlicerStats{}
			if err := slicer.Slice(sr, &tt.options, &stats); err != nil {
				t.Errorf("Unexpected Slicer error %v", err)
			}
			if stats.IgnoredFile != tt.expected {
				t.Errorf("expected ignore %t, got %t", tt.expected, stats.IgnoredFile)
			}
		})
	}
}
func TestSliceFS_SliceFS_WithDisabledOptions_ReturnsValidHash(t *testing.T) {
	fsys := os.DirFS("./artefacts")
	algorithm := algorithms.Xxhash

	tests := []struct {
		name     string
		filename string
		options  Options
		expected string
	}{
		{
			name:     "FileSystemTestFile_TestWordPdf_WithSlicing",
			filename: "test.pdf",
			options: Options{
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: "bedd0999e968547e",
		},
		{
			name:     "FileSystemTestFile_TestAdobePdf_WithSlicing",
			filename: "test-adobe.pdf",
			options: Options{
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: "41d4d0a83d10e962",
		},
		{
			name:     "FileSystemTestFile_Test1mb_WithSlicing",
			filename: "test.1mb",
			options: Options{
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: "bb83f43630ee546f",
		},
		{
			name:     "FileSystemTestFile_Test1mb_WithoutMeta",
			filename: "test.1mb",
			options: Options{
				DisableSlicing:  false,
				DisableMeta:     true,
				DisableAutoText: false,
			},
			expected: "913c30543271faaf",
		},
		{
			name:     "FileSystemTestFile_TestManipulated1mb_WithSlicing",
			filename: "test-manipulated.1mb",
			options: Options{
				DisableSlicing:  false,
				DisableMeta:     false,
				DisableAutoText: false,
			},
			expected: "4f595576799edcd9",
		},
		// Add other test cases here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			slicer := New(algorithm)

			if stats, err := slicer.SliceFS(fsys, tt.filename, &tt.options); err != nil {
				t.Errorf("Unexpected Slicer error %v", err)
			} else {

				actual := hex.EncodeToString(stats.Hash)

				if len(tt.expected) != len(actual) {
					t.Errorf("hash length expected %d, got %d", len(tt.expected), len(actual))
				}

				if !strings.EqualFold(actual, tt.expected) {
					t.Errorf("expected hash %s, got %s for %s", tt.expected, actual, tt.filename)
				}
			}
		})
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

	options := Options{
		DisableSlicing:  false,
		DisableMeta:     false,
		DisableAutoText: false,
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
