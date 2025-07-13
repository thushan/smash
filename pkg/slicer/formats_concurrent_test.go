package slicer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestSlicingSupportedConcurrent(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		size     uint64
		expected bool
	}{
		{
			name:     "binary data",
			data:     []byte{0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF},
			size:     6,
			expected: true,
		},
		{
			name:     "text data",
			data:     []byte("Hello, World! This is a text file."),
			size:     34,
			expected: true, // slicingSupported returns true when read succeeds
		},
		{
			name:     "empty data",
			data:     []byte{},
			size:     0,
			expected: false,
		},
		{
			name:     "mixed content",
			data:     append([]byte("Text"), []byte{0x00, 0xFF, 0x00}...),
			size:     7,
			expected: true,
		},
	}

	const goroutines = 100
	const iterations = 100

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, _ := errgroup.WithContext(context.Background())

			for i := 0; i < goroutines; i++ {
				id := i // capture loop variable
				g.Go(func() error {
					for j := 0; j < iterations; j++ {
						reader := bytes.NewReader(tc.data)
						sr := io.NewSectionReader(reader, 0, int64(tc.size))

						result := slicingSupported(sr, tc.size)
						if result != tc.expected {
							return fmt.Errorf("goroutine %d iteration %d: expected %v, got %v", id, j, tc.expected, result)
						}
					}
					return nil
				})
			}

			if err := g.Wait(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDetectBufferPoolRaceCondition(t *testing.T) {
	// This test specifically targets the race condition we fixed in 0.9.5a & expands on it
	const workers = 50
	const filesPerWorker = 100

	// Create diverse test data to increase chance of detecting corruption
	testData := [][]byte{
		bytes.Repeat([]byte("A"), MaxReadBytes),
		bytes.Repeat([]byte{0xFF}, MaxReadBytes),
		bytes.Repeat([]byte{0x00}, MaxReadBytes),
		[]byte("Mixed content with text and binary\x00\xFF\x00"),
	}

	results := make([][]bool, workers)

	g, _ := errgroup.WithContext(context.Background())

	for w := 0; w < workers; w++ {
		results[w] = make([]bool, filesPerWorker)
		workerID := w // capture loop variable
		g.Go(func() error {
			for f := 0; f < filesPerWorker; f++ {
				data := testData[(workerID+f)%len(testData)]
				reader := bytes.NewReader(data)
				sr := io.NewSectionReader(reader, 0, int64(len(data)))

				results[workerID][f] = slicingSupported(sr, uint64(len(data)))
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	// Verify results are consistent
	for w := 0; w < workers; w++ {
		for f := 0; f < filesPerWorker; f++ {
			// slicingSupported returns true for most cases except tiny files
			// The actual logic doesn't depend on text vs binary for normal reads
			expected := true
			if results[w][f] != expected {
				t.Errorf("worker %d file %d: inconsistent result, possible race condition", w, f)
			}
		}
	}
}

func BenchmarkSlicingSupportedConcurrent(b *testing.B) {
	testData := []byte("This is a sample text file with enough content to test detection.")
	size := uint64(len(testData))

	b.Run("sequential", func(b *testing.B) {
		reader := bytes.NewReader(testData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sr := io.NewSectionReader(reader, 0, int64(size))
			_ = slicingSupported(sr, size)
			reader.Seek(0, io.SeekStart)
		}
	})

	b.Run("concurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			reader := bytes.NewReader(testData)
			for pb.Next() {
				sr := io.NewSectionReader(reader, 0, int64(size))
				_ = slicingSupported(sr, size)
				reader.Seek(0, io.SeekStart)
			}
		})
	})
}
