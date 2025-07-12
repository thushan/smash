package smash

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/thushan/smash/internal/algorithms"
	"github.com/thushan/smash/pkg/indexer"
	"golang.org/x/sync/errgroup"
)

func init() {
	// Set test environment variable so the app knows it's in test mode
	os.Setenv("GO_TEST", "1")
}

func TestAppConcurrentFileProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create temporary test directory structure
	tempDir := t.TempDir()

	// Create test files with duplicate content
	testFiles := []struct {
		path    string
		content []byte
	}{
		{"file1.txt", []byte("duplicate content 1")},
		{"file2.txt", []byte("duplicate content 1")},
		{"dir1/file3.txt", []byte("duplicate content 2")},
		{"dir1/file4.txt", []byte("duplicate content 2")},
		{"dir2/file5.bin", []byte{0xFF, 0x00, 0xFF, 0x00}},
		{"dir2/file6.bin", []byte{0xFF, 0x00, 0xFF, 0x00}},
		{"unique.txt", []byte("unique content")},
		{"empty.txt", []byte{}},
	}

	// Create directories first
	dirs := make(map[string]bool)
	for _, tf := range testFiles {
		dir := filepath.Dir(filepath.Join(tempDir, tf.path))
		dirs[dir] = true
	}
	for dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create files in parallel
	g, _ := errgroup.WithContext(context.Background())
	for _, tf := range testFiles {
		tf := tf // capture loop variable
		g.Go(func() error {
			fullPath := filepath.Join(tempDir, tf.path)
			return os.WriteFile(fullPath, tf.content, 0644)
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatalf("failed to create test files: %v", err)
	}

	// Test with various worker counts
	workerCounts := []int{1, 4, 16, 32}

	for _, workers := range workerCounts {
		t.Run(fmt.Sprintf("workers_%d", workers), func(t *testing.T) {
			app := &App{
				Flags: &Flags{
					Algorithm:       int(algorithms.Xxhash),
					MaxWorkers:      workers,
					MaxThreads:      workers,
					Slices:          4,
					SliceSize:       8192,
					SliceThreshold:  102400,
					Recurse:         true,
					Silent:          true,
					HideProgress:    true,
					ShowTop:         10,
					ProgressUpdate:  5,
					DisableMeta:     true,
					DisableAutoText: true,
				},
				Args: []string{tempDir},
				Locations: []indexer.LocationFS{
					{Name: tempDir, FS: os.DirFS(tempDir)},
				},
			}

			err := app.Run()
			if err != nil {
				t.Errorf("app.Run() failed with %d workers: %v", workers, err)
			}

			// Verify results
			duplicateCount := 0
			app.Session.Dupes.Range(func(key string, value *DuplicateFiles) bool {
				duplicateCount++
				return true
			})

			expectedDuplicates := 3 // 3 sets of duplicates
			if duplicateCount != expectedDuplicates {
				t.Errorf("expected %d duplicate sets, got %d with %d workers", expectedDuplicates, duplicateCount, workers)
			}

			emptyCount := len(app.Session.Empty.Files)
			if emptyCount != 1 {
				t.Errorf("expected 1 empty file, got %d with %d workers", emptyCount, workers)
			}
		})
	}
}

func TestAppWithInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		flags       Flags
		wantError   bool
		errorString string
	}{
		{
			name: "negative slice size",
			flags: Flags{
				SliceSize: -1,
			},
			wantError:   true,
			errorString: "slice size and threshold must be non-negative",
		},
		{
			name: "negative min size",
			flags: Flags{
				MinSize: -1024,
			},
			wantError:   true,
			errorString: "min size and max size must be non-negative",
		},
		{
			name: "negative max workers",
			flags: Flags{
				MaxWorkers: -1,
			},
			wantError:   true,
			errorString: "maxworkers cannot be below zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				Flags: &tt.flags,
				Args:  []string{"."},
			}

			err := app.Run()
			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tt.errorString != "" && err.Error() != tt.errorString {
					t.Errorf("expected error %q, got %q", tt.errorString, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAppMemoryUsageUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping memory test in short mode")
	}

	tempDir := t.TempDir()

	// Create many small files to stress memory usage
	fileCount := 1000
	g, _ := errgroup.WithContext(context.Background())

	// Limit concurrency to avoid file descriptor exhaustion
	g.SetLimit(100)

	for i := 0; i < fileCount; i++ {
		i := i // capture loop variable
		g.Go(func() error {
			content := []byte(fmt.Sprintf("test content %d", i))
			path := filepath.Join(tempDir, fmt.Sprintf("file_%d.txt", i))
			return os.WriteFile(path, content, 0644)
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatalf("failed to create test files: %v", err)
	}

	app := &App{
		Flags: &Flags{
			Algorithm:      int(algorithms.Xxhash),
			MaxWorkers:     16,
			MaxThreads:     16,
			Slices:         4,
			SliceSize:      8192,
			SliceThreshold: 102400,
			Recurse:        true,
			Silent:         true,
			HideProgress:   true,
			ShowTop:        10,
			ProgressUpdate: 5,
		},
		Args: []string{tempDir},
		Locations: []indexer.LocationFS{
			{Name: tempDir, FS: os.DirFS(tempDir)},
		},
	}

	// Measure memory before
	var memBefore runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	start := time.Now()
	err := app.Run()
	duration := time.Since(start)

	if err != nil {
		t.Errorf("app.Run() failed: %v", err)
	}

	// Measure memory after
	runtime.GC()
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Handle potential wraparound
	var allocatedMB float64
	if memAfter.Alloc >= memBefore.Alloc {
		allocatedMB = float64(memAfter.Alloc-memBefore.Alloc) / 1024 / 1024
	} else {
		allocatedMB = 0 // Memory was freed, not allocated
	}
	t.Logf("Processed %d files in %v, allocated %.2f MB", fileCount, duration, allocatedMB)

	// Rough heuristic: shouldn't allocate more than 100KB per file
	maxExpectedMB := float64(fileCount) * 100 / 1024
	if allocatedMB > maxExpectedMB {
		t.Errorf("excessive memory usage: %.2f MB allocated, expected less than %.2f MB", allocatedMB, maxExpectedMB)
	}
}

func TestAppRaceConditions(t *testing.T) {
	// This test runs multiple app instances concurrently to detect race conditions
	tempDir := t.TempDir()

	// Create test files in parallel
	fileGroup, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 10; i++ {
		i := i // capture loop variable
		fileGroup.Go(func() error {
			content := []byte(fmt.Sprintf("content %d", i))
			path := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
			return os.WriteFile(path, content, 0644)
		})
	}
	if err := fileGroup.Wait(); err != nil {
		t.Fatalf("failed to create test files: %v", err)
	}

	const instances = 5
	g, _ := errgroup.WithContext(context.Background())

	for i := 0; i < instances; i++ {
		g.Go(func() error {
			app := &App{
				Flags: &Flags{
					Algorithm:       int(algorithms.Xxhash),
					MaxWorkers:      4,
					MaxThreads:      4,
					Slices:          4,
					SliceSize:       8192,
					SliceThreshold:  102400,
					Silent:          true,
					HideProgress:    true,
					ShowTop:         10,
					ProgressUpdate:  5,
					DisableMeta:     true,
					DisableAutoText: true,
				},
				Args: []string{tempDir},
				Locations: []indexer.LocationFS{
					{Name: tempDir, FS: os.DirFS(tempDir)},
				},
			}

			return app.Run()
		})
	}

	if err := g.Wait(); err != nil {
		t.Errorf("concurrent app execution failed: %v", err)
	}
}
