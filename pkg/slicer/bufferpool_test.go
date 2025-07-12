package slicer

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestGetSliceBuffer(t *testing.T) {
	tests := []struct {
		name     string
		size     uint64
		wantSize int
	}{
		{
			name:     "default size returns pooled buffer",
			size:     DefaultSliceSize,
			wantSize: DefaultSliceSize,
		},
		{
			name:     "non-default size creates new buffer",
			size:     1024,
			wantSize: 1024,
		},
		{
			name:     "zero size creates empty buffer",
			size:     0,
			wantSize: 0,
		},
		{
			name:     "large size creates new buffer",
			size:     1024 * 1024,
			wantSize: 1024 * 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := getSliceBuffer(tt.size)
			if len(buf) != tt.wantSize {
				t.Errorf("getSliceBuffer() returned buffer of size %d, want %d", len(buf), tt.wantSize)
			}
			if cap(buf) < tt.wantSize {
				t.Errorf("getSliceBuffer() returned buffer with capacity %d, want at least %d", cap(buf), tt.wantSize)
			}
		})
	}
}

func TestPutSliceBuffer(t *testing.T) {
	t.Run("returns default size buffer to pool", func(t *testing.T) {
		buf := make([]byte, DefaultSliceSize)
		putSliceBuffer(buf) // Should not panic
	})

	t.Run("does not return non-default size buffer", func(t *testing.T) {
		buf := make([]byte, 1024)
		putSliceBuffer(buf) // Should not panic, but won't pool
	})

	t.Run("handles nil buffer", func(t *testing.T) {
		putSliceBuffer(nil) // Should not panic
	})
}

func TestBufferPoolConcurrency(t *testing.T) {
	const goroutines = 100
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				buf := getSliceBuffer(DefaultSliceSize)
				if len(buf) != DefaultSliceSize {
					t.Errorf("concurrent getSliceBuffer returned wrong size: %d", len(buf))
				}
				// Simulate some work
				buf[0] = byte(j % 256)
				putSliceBuffer(buf)
			}
		}()
	}

	wg.Wait()
}

func TestBufferPoolStress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	const workers = 50
	const duration = 100 // milliseconds

	stop := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(workers)

	allocations := make([]int64, workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			var count int64
			for {
				select {
				case <-stop:
					allocations[id] = count
					return
				default:
					buf := getSliceBuffer(DefaultSliceSize)
					// Write pattern to detect corruption
					for j := 0; j < len(buf); j += 4 {
						buf[j] = byte(id)
					}
					// Verify no corruption from other goroutines
					for j := 0; j < len(buf); j += 4 {
						if buf[j] != byte(id) {
							t.Errorf("buffer corruption detected: expected %d, got %d", id, buf[j])
							return
						}
					}
					putSliceBuffer(buf)
					count++
				}
			}
		}(i)
	}

	// Let the test run for the specified duration
	time.Sleep(time.Duration(duration) * time.Millisecond)
	close(stop)
	wg.Wait()

	var total int64
	for _, count := range allocations {
		total += count
	}
	t.Logf("Processed %d buffer allocations across %d workers", total, workers)
}

func BenchmarkBufferPool(b *testing.B) {
	b.Run("with pool", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := getSliceBuffer(DefaultSliceSize)
			putSliceBuffer(buf)
		}
	})

	b.Run("without pool", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = make([]byte, DefaultSliceSize)
		}
	})

	b.Run("concurrent with pool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				buf := getSliceBuffer(DefaultSliceSize)
				putSliceBuffer(buf)
			}
		})
	})

	b.Run("concurrent without pool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = make([]byte, DefaultSliceSize)
			}
		})
	})
}

func BenchmarkBufferPoolMemory(b *testing.B) {
	b.Run("memory usage with pool", func(b *testing.B) {
		var m runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m)
		startAlloc := m.Alloc

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := getSliceBuffer(DefaultSliceSize)
			// Simulate sume work
			buf[0] = byte(i)
			buf[len(buf)-1] = byte(i)
			putSliceBuffer(buf)
		}

		runtime.GC()
		runtime.ReadMemStats(&m)
		endAlloc := m.Alloc
		b.ReportMetric(float64(endAlloc-startAlloc)/float64(b.N), "bytes/op")
	})
}
