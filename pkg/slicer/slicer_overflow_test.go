package slicer

import (
	"bytes"
	"io"
	"io/fs"
	"math"
	"os"
	"testing"
	"time"

	"github.com/thushan/smash/internal/algorithms"
)

type mockFileInfo struct {
	name string
	size int64
	mode os.FileMode
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() interface{}   { return nil }

type mockFS struct {
	files map[string]mockFile
}

type mockFile struct {
	data []byte
	info mockFileInfo
	err  error
}

func (m mockFS) Open(name string) (fs.File, error) {
	if f, ok := m.files[name]; ok {
		if f.err != nil {
			return nil, f.err
		}
		return &mockOpenFile{Reader: bytes.NewReader(f.data), info: f.info}, nil
	}
	return nil, os.ErrNotExist
}

func (m mockFS) Stat(name string) (fs.FileInfo, error) {
	if f, ok := m.files[name]; ok {
		return f.info, nil
	}
	return nil, os.ErrNotExist
}

type mockOpenFile struct {
	*bytes.Reader
	info mockFileInfo
}

func (m *mockOpenFile) Stat() (fs.FileInfo, error) { return m.info, nil }
func (m *mockOpenFile) Close() error               { return nil }

func TestSliceFSWithNegativeFileSize(t *testing.T) {
	slicer := New(algorithms.Algorithm(algorithms.Xxhash))

	mockFS := mockFS{
		files: map[string]mockFile{
			"negative.bin": {
				data: []byte("test data"),
				info: mockFileInfo{
					name: "negative.bin",
					size: -1, // Negative size
					mode: 0644,
				},
			},
		},
	}

	stats, err := slicer.SliceFS(mockFS, "negative.bin", &Options{})
	if err == nil {
		t.Error("expected error for negative file size, got nil")
	}
	if err != nil && err.Error() != "file size cannot be negative" {
		t.Errorf("unexpected error message: %v", err)
	}
	if stats.Hash == nil || len(stats.Hash) != 0 {
		t.Error("expected empty hash slice for error case")
	}
}

func TestSliceWithBoundaryConditions(t *testing.T) {
	tests := []struct {
		name      string
		fileSize  uint64
		options   Options
		wantError bool
	}{
		{
			name:      "zero byte file",
			fileSize:  0,
			options:   Options{},
			wantError: false,
		},
		{
			name:      "one byte file",
			fileSize:  1,
			options:   Options{},
			wantError: false,
		},
		{
			name:     "max uint64 file size",
			fileSize: math.MaxUint64,
			options: Options{
				MaxSize: math.MaxUint64 - 1,
			},
			wantError: false, // Should be ignored due to MaxSize
		},
		{
			name:     "file at min size boundary",
			fileSize: 1024,
			options: Options{
				MinSize: 1024,
			},
			wantError: false,
		},
		{
			name:     "file below min size",
			fileSize: 1023,
			options: Options{
				MinSize: 1024,
			},
			wantError: false, // Ignored, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slicer := New(algorithms.Algorithm(algorithms.Xxhash))

			data := make([]byte, minUint64(tt.fileSize, 1024*1024)) // Cap at 1MB for tests
			reader := bytes.NewReader(data)
			sr := io.NewSectionReader(reader, 0, int64(len(data)))

			stats := SlicerStats{}
			err := slicer.Slice(sr, &tt.options, &stats)

			if tt.wantError && err == nil {
				t.Error("expected error but got nil")
			} else if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestSlicesPlus2Overflow(t *testing.T) {
	// Test the specific case where slices + 2 could overflow
	slicer := NewConfigured(algorithms.Algorithm(algorithms.Xxhash), math.MaxInt-1, 8192, DefaultThreshold)

	data := make([]byte, 1024)
	reader := bytes.NewReader(data)
	sr := io.NewSectionReader(reader, 0, int64(len(data)))

	stats := SlicerStats{}
	err := slicer.Slice(sr, &Options{}, &stats)

	if err == nil {
		t.Error("expected error for slices overflow")
	}
	if err != nil && !bytes.Contains([]byte(err.Error()), []byte("slices overflow")) {
		t.Errorf("unexpected error: %v", err)
	}
}

func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
