package smash

import (
	"path/filepath"
	"testing"

	"github.com/thushan/smash/pkg/analysis"
)

func TestGetUsername(t *testing.T) {
	// This test just ensures the function doesn't panic
	// and returns a non-empty string
	username := getUsername()
	if username == "" {
		t.Error("getUsername() returned an empty string")
	}
}

func TestGetHostName(t *testing.T) {
	// This test just ensures the function doesn't panic
	// and returns a non-empty string
	hostname := getHostName()
	if hostname == "" {
		t.Error("getHostName() returned an empty string")
	}
}

func TestTransformTopFiles(t *testing.T) {
	tests := []struct {
		name     string
		files    []analysis.Item
		expected []ReportTopFilesSummary
	}{
		{
			name:     "Should return empty slice for empty input",
			files:    []analysis.Item{},
			expected: []ReportTopFilesSummary{},
		},
		{
			name: "Should transform single item",
			files: []analysis.Item{
				{Key: "hash1", Size: 100},
			},
			expected: []ReportTopFilesSummary{
				{Hash: "hash1", Size: 100},
			},
		},
		{
			name: "Should transform multiple items",
			files: []analysis.Item{
				{Key: "hash1", Size: 100},
				{Key: "hash2", Size: 200},
				{Key: "hash3", Size: 300},
			},
			expected: []ReportTopFilesSummary{
				{Hash: "hash1", Size: 100},
				{Hash: "hash2", Size: 200},
				{Hash: "hash3", Size: 300},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformTopFiles(tt.files)

			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("transformTopFiles() returned %d items, want %d", len(result), len(tt.expected))
				return
			}

			// Check each item
			for i, item := range result {
				if item.Hash != tt.expected[i].Hash || item.Size != tt.expected[i].Size {
					t.Errorf("transformTopFiles()[%d] = {Hash: %s, Size: %d}, want {Hash: %s, Size: %d}",
						i, item.Hash, item.Size, tt.expected[i].Hash, tt.expected[i].Size)
				}
			}
		})
	}
}

func TestSummariseSmashedFile(t *testing.T) {
	// Use filepath.Join to create platform-specific paths
	testPath := filepath.Join("path", "to", "test.txt")
	expectedDir := filepath.Dir(testPath)

	tests := []struct {
		name     string
		file     File
		expected ReportFileSummary
	}{
		{
			name: "Should summarize file correctly",
			file: File{
				Filename: "test.txt",
				Location: "location1",
				Path:     testPath,
				Hash:     "hash1",
				FullHash: true,
				FileSize: 100,
			},
			expected: ReportFileSummary{
				ReportFileBaseSummary: ReportFileBaseSummary{
					Filename: "test.txt",
					Location: "location1",
					Path:     expectedDir,
				},
				Hash:     "hash1",
				Size:     100,
				FullHash: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := summariseSmashedFile(tt.file)

			// Check fields
			if result.Filename != tt.expected.Filename {
				t.Errorf("summariseSmashedFile().Filename = %s, want %s", result.Filename, tt.expected.Filename)
			}
			if result.Location != tt.expected.Location {
				t.Errorf("summariseSmashedFile().Location = %s, want %s", result.Location, tt.expected.Location)
			}
			if result.Path != tt.expected.Path {
				t.Errorf("summariseSmashedFile().Path = %s, want %s", result.Path, tt.expected.Path)
			}
			if result.Hash != tt.expected.Hash {
				t.Errorf("summariseSmashedFile().Hash = %s, want %s", result.Hash, tt.expected.Hash)
			}
			if result.Size != tt.expected.Size {
				t.Errorf("summariseSmashedFile().Size = %d, want %d", result.Size, tt.expected.Size)
			}
			if result.FullHash != tt.expected.FullHash {
				t.Errorf("summariseSmashedFile().FullHash = %t, want %t", result.FullHash, tt.expected.FullHash)
			}
		})
	}
}
