package slicer

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"

	"github.com/thushan/smash/internal/algorithms"
)

func BenchmarkSlice(b *testing.B) {
	// Create test data of various sizes
	sizes := []struct {
		name string
		size int
	}{
		{"1KB", 1024},
		{"10KB", 10 * 1024},
		{"100KB", 100 * 1024},
		{"1MB", 1024 * 1024},
		{"10MB", 10 * 1024 * 1024},
	}

	for _, size := range sizes {
		data := make([]byte, size.size)
		rand.Read(data)

		b.Run(size.name, func(b *testing.B) {
			b.Run("full_hash", func(b *testing.B) {
				benchmarkSliceFullHash(b, data)
			})

			b.Run("sliced", func(b *testing.B) {
				benchmarkSliceWithSlicing(b, data)
			})
		})
	}
}

func benchmarkSliceFullHash(b *testing.B, data []byte) {
	slicer := New(algorithms.Algorithm(algorithms.Xxhash))
	options := &Options{
		DisableSlicing: true,
	}

	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		sr := io.NewSectionReader(reader, 0, int64(len(data)))
		stats := SlicerStats{}

		if err := slicer.Slice(sr, options, &stats); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkSliceWithSlicing(b *testing.B, data []byte) {
	slicer := New(algorithms.Algorithm(algorithms.Xxhash))
	options := &Options{}

	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		sr := io.NewSectionReader(reader, 0, int64(len(data)))
		stats := SlicerStats{}

		if err := slicer.Slice(sr, options, &stats); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSliceAlgorithms(b *testing.B) {
	// Compare performance of different algorithms
	algorithms := []struct {
		name string
		algo algorithms.Algorithm
	}{
		{"xxhash", algorithms.Xxhash},
		{"murmur3", algorithms.Murmur3_128},
		{"sha256", algorithms.Sha256},
		{"md5", algorithms.Md5},
	}

	data := make([]byte, 1024*1024) // 1MB test file
	rand.Read(data)

	for _, alg := range algorithms {
		b.Run(alg.name, func(b *testing.B) {
			slicer := New(alg.algo)
			options := &Options{}

			b.SetBytes(int64(len(data)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				reader := bytes.NewReader(data)
				sr := io.NewSectionReader(reader, 0, int64(len(data)))
				stats := SlicerStats{}

				if err := slicer.Slice(sr, options, &stats); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSliceConcurrent(b *testing.B) {
	data := make([]byte, 100*1024) // 100KB files
	rand.Read(data)

	slicer := New(algorithms.Algorithm(algorithms.Xxhash))
	options := &Options{}

	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reader := bytes.NewReader(data)
			sr := io.NewSectionReader(reader, 0, int64(len(data)))
			stats := SlicerStats{}

			if err := slicer.Slice(sr, options, &stats); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSliceMemoryAllocation(b *testing.B) {
	sizes := []int{
		8 * 1024,    // 8KB
		64 * 1024,   // 64KB
		512 * 1024,  // 512KB
		1024 * 1024, // 1MB
	}

	for _, size := range sizes {
		data := make([]byte, size)
		rand.Read(data)

		b.Run(fmt.Sprintf("size_%dKB", size/1024), func(b *testing.B) {
			slicer := New(algorithms.Algorithm(algorithms.Xxhash))
			options := &Options{}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				reader := bytes.NewReader(data)
				sr := io.NewSectionReader(reader, 0, int64(len(data)))
				stats := SlicerStats{}

				if err := slicer.Slice(sr, options, &stats); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSliceTextVsBinary(b *testing.B) {
	// Compare performance between text and binary files
	textData := bytes.Repeat([]byte("Hello, World! This is a text file.\n"), 1000)
	binaryData := make([]byte, len(textData))
	rand.Read(binaryData)

	b.Run("text_file", func(b *testing.B) {
		benchmarkSliceData(b, textData)
	})

	b.Run("binary_file", func(b *testing.B) {
		benchmarkSliceData(b, binaryData)
	})
}

func benchmarkSliceData(b *testing.B, data []byte) {
	slicer := New(algorithms.Algorithm(algorithms.Xxhash))
	options := &Options{}

	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		sr := io.NewSectionReader(reader, 0, int64(len(data)))
		stats := SlicerStats{}

		if err := slicer.Slice(sr, options, &stats); err != nil {
			b.Fatal(err)
		}
	}
}
