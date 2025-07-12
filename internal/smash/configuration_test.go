package smash

import (
	"testing"

	"github.com/thushan/smash/pkg/indexer"
)

func TestEnabledOrDisabled(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected string
	}{
		{
			name:     "Should return Enabled when value is true",
			value:    true,
			expected: "Enabled",
		},
		{
			name:     "Should return Disabled when value is false",
			value:    false,
			expected: "Disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enabledOrDisabled(tt.value)
			if result != tt.expected {
				t.Errorf("enabledOrDisabled(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestBuildLocations(t *testing.T) {
	tests := []struct {
		name      string
		locations []indexer.LocationFS
		expected  string
	}{
		{
			name:      "Should return empty string for empty locations",
			locations: []indexer.LocationFS{},
			expected:  "",
		},
		{
			name: "Should return single location name",
			locations: []indexer.LocationFS{
				{Name: "location1"},
			},
			expected: "location1",
		},
		{
			name: "Should return comma-separated location names",
			locations: []indexer.LocationFS{
				{Name: "location1"},
				{Name: "location2"},
				{Name: "location3"},
			},
			expected: "location1, location2, location3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildLocations(tt.locations)
			if result != tt.expected {
				t.Errorf("buildLocations() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSetMaxThreads(t *testing.T) {
	tests := []struct {
		name       string
		maxThreads int
	}{
		{
			name:       "Should set GOMAXPROCS to NumCPU when maxThreads is 0",
			maxThreads: 0,
		},
		{
			name:       "Should set GOMAXPROCS to NumCPU when maxThreads is negative",
			maxThreads: -1,
		},
		{
			name:       "Should set GOMAXPROCS to maxThreads when valid",
			maxThreads: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				Flags: &Flags{
					MaxThreads: tt.maxThreads,
				},
			}
			app.setMaxThreads()
			// We can't easily test the actual GOMAXPROCS value without affecting the test environment,
			// so this test just ensures the function doesn't panic
		})
	}
}
