package slicer

import (
	"bytes"
	"context"
	"io"
	"testing"

	"golang.org/x/sync/errgroup"
	"github.com/thushan/smash/internal/algorithms"
)

func TestSlicerBasicFunctionality(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expectErr bool
	}{
		{
			name:      "empty file",
			data:      []byte{},
			expectErr: false,
		},
		{
			name:      "small text file",
			data:      []byte("Hello, World!"),
			expectErr: false,
		},
		{
			name:      "binary data",
			data:      []byte{0xFF, 0x00, 0xFF, 0x00},
			expectErr: false,
		},
		{
			name:      "1KB file",
			data:      make([]byte, 1024),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slicer := New(algorithms.Algorithm(algorithms.Xxhash))
			reader := bytes.NewReader(tt.data)
			sr := io.NewSectionReader(reader, 0, int64(len(tt.data)))
			
			stats := SlicerStats{}
			err := slicer.Slice(sr, &Options{}, &stats)
			
			if tt.expectErr && err == nil {
				t.Error("expected error but got nil")
			} else if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			
			if !tt.expectErr && len(stats.Hash) == 0 && len(tt.data) > 0 {
				t.Error("expected non-empty hash for non-empty file")
			}
		})
	}
}

func TestSlicerConcurrentUsage(t *testing.T) {
	// Test that multiple goroutines can use the slicer concurrently
	slicer := New(algorithms.Algorithm(algorithms.Xxhash))
	data := []byte("test data for concurrent access")
	
	g, _ := errgroup.WithContext(context.Background())
	
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			reader := bytes.NewReader(data)
			sr := io.NewSectionReader(reader, 0, int64(len(data)))
			stats := SlicerStats{}
			
			if err := slicer.Slice(sr, &Options{}, &stats); err != nil {
				return err
			}
			return nil
		})
	}
	
	if err := g.Wait(); err != nil {
		t.Errorf("concurrent slice failed: %v", err)
	}
}

func TestBufferPoolDoesNotPanic(t *testing.T) {
	// Test various buffer sizes don't panic
	sizes := []uint64{0, 1, 1024, DefaultSliceSize, 1024 * 1024}
	
	for _, size := range sizes {
		buf := getSliceBuffer(size)
		if size <= 1<<30 && buf == nil {
			t.Errorf("getSliceBuffer(%d) returned nil unexpectedly", size)
		}
		if buf != nil {
			putSliceBuffer(buf)
		}
	}
}