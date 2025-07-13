package algorithms

import (
	"bytes"
	"fmt"
	"hash"
	"testing"
)

func TestAlgorithmNew(t *testing.T) {
	tests := []struct {
		name     string
		algo     Algorithm
		wantType string
	}{
		{
			name:     "xxhash",
			algo:     Xxhash,
			wantType: "*xxhash.digest",
		},
		{
			name:     "murmur3_32",
			algo:     Murmur3_32,
			wantType: "*murmur3.digest32",
		},
		{
			name:     "murmur3_128",
			algo:     Murmur3_128,
			wantType: "*murmur3.digest128",
		},
		{
			name:     "sha256",
			algo:     Sha256,
			wantType: "*sha256.digest",
		},
		{
			name:     "md5",
			algo:     Md5,
			wantType: "*md5.digest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.algo.New()
			if h == nil {
				t.Fatal("New() returned nil")
			}

			// Verify it implements hash.Hash
			if _, ok := h.(hash.Hash); !ok {
				t.Error("returned value does not implement hash.Hash")
			}
		})
	}
}

func TestAlgorithmConsistency(t *testing.T) {
	// Test that algorithms produce consistent results
	testData := [][]byte{
		[]byte(""),
		[]byte("a"),
		[]byte("abc"),
		[]byte("The quick brown fox jumps over the lazy dog"),
		bytes.Repeat([]byte("x"), 1024),
	}

	algorithms := []Algorithm{
		Xxhash,
		Murmur3_32,
		Murmur3_128,
		Sha256,
		Md5,
	}

	for _, algo := range algorithms {
		t.Run(algo.String(), func(t *testing.T) {
			for i, data := range testData {
				// Hash the data twice and ensure results match
				h1 := algo.New()
				h1.Write(data)
				sum1 := h1.Sum(nil)

				h2 := algo.New()
				h2.Write(data)
				sum2 := h2.Sum(nil)

				if !bytes.Equal(sum1, sum2) {
					t.Errorf("inconsistent hash for data[%d]: %x != %x", i, sum1, sum2)
				}
			}
		})
	}
}

func TestAlgorithmReset(t *testing.T) {
	algorithms := []Algorithm{
		Xxhash,
		Murmur3_32,
		Murmur3_128,
		Sha256,
		Md5,
	}

	data1 := []byte("first data")
	data2 := []byte("second data")

	for _, algo := range algorithms {
		t.Run(algo.String(), func(t *testing.T) {
			h := algo.New()

			// Hash first data
			h.Write(data1)
			sum1 := h.Sum(nil)

			// Reset and hash second data
			h.Reset()
			h.Write(data2)
			sum2 := h.Sum(nil)

			// Hash second data with fresh hasher
			h2 := algo.New()
			h2.Write(data2)
			sum2Fresh := h2.Sum(nil)

			if bytes.Equal(sum1, sum2) {
				t.Error("reset did not clear state")
			}

			if !bytes.Equal(sum2, sum2Fresh) {
				t.Error("reset hasher produced different result than fresh hasher")
			}
		})
	}
}

func TestAlgorithmIncrementalHashing(t *testing.T) {
	// Test that incremental writes produce same result as single write
	algorithms := []Algorithm{
		Xxhash,
		Murmur3_32,
		Murmur3_128,
		Sha256,
		Md5,
	}

	fullData := []byte("The quick brown fox jumps over the lazy dog")

	for _, algo := range algorithms {
		t.Run(algo.String(), func(t *testing.T) {
			// Single write
			h1 := algo.New()
			h1.Write(fullData)
			sum1 := h1.Sum(nil)

			// Incremental writes
			h2 := algo.New()
			h2.Write(fullData[:10])
			h2.Write(fullData[10:20])
			h2.Write(fullData[20:])
			sum2 := h2.Sum(nil)

			if !bytes.Equal(sum1, sum2) {
				t.Errorf("incremental hashing produced different result: %x != %x", sum1, sum2)
			}
		})
	}
}

func TestAlgorithmKnownValues(t *testing.T) {
	// Test against known hash values
	data := []byte("test")

	tests := []struct {
		algo     Algorithm
		expected string
	}{
		{Md5, "098f6bcd4621d373cade4e832627b4f6"},
		{Sha256, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
		// Note: xxhash and murmur3 values depend on seed/implementation
	}

	for _, tt := range tests {
		t.Run(tt.algo.String(), func(t *testing.T) {
			h := tt.algo.New()
			h.Write(data)
			sum := h.Sum(nil)

			// Skip non-deterministic algorithms
			if tt.expected == "" {
				return
			}

			got := fmt.Sprintf("%x", sum)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func BenchmarkAlgorithms(b *testing.B) {
	sizes := []int{
		64,
		1024,
		16 * 1024,
		1024 * 1024,
	}

	algorithms := []Algorithm{
		Xxhash,
		Murmur3_32,
		Murmur3_128,
		Sha256,
		Md5,
	}

	for _, size := range sizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		for _, algo := range algorithms {
			b.Run(algo.String()+"_"+humanizeBytes(size), func(b *testing.B) {
				b.SetBytes(int64(size))
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					h := algo.New()
					h.Write(data)
					_ = h.Sum(nil)
				}
			})
		}
	}
}

func BenchmarkAlgorithmsParallel(b *testing.B) {
	data := make([]byte, 1024*1024) // 1MB
	for i := range data {
		data[i] = byte(i % 256)
	}

	algorithms := []Algorithm{
		Xxhash,
		Murmur3_32,
		Murmur3_128,
		Sha256,
		Md5,
	}

	for _, algo := range algorithms {
		b.Run(algo.String(), func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					h := algo.New()
					h.Write(data)
					_ = h.Sum(nil)
				}
			})
		})
	}
}

func humanizeBytes(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%d%cB", b/div, "KMGTPE"[exp])
}
