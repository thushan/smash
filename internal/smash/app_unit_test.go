package smash

import (
	"testing"

	"github.com/thushan/smash/internal/algorithms"
)

func TestAppValidateArgs(t *testing.T) {
	tests := []struct {
		name      string
		flags     Flags
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid configuration",
			flags: Flags{
				MaxWorkers:     4,
				MaxThreads:     4,
				ShowTop:        10,
				ProgressUpdate: 5,
				SliceSize:      8192,
				SliceThreshold: 102400,
				Slices:         4,
			},
			wantError: false,
		},
		{
			name: "negative max workers",
			flags: Flags{
				MaxWorkers: -1,
				ShowTop:    10,
			},
			wantError: true,
			errorMsg:  "maxworkers cannot be below zero",
		},
		{
			name: "negative slice size",
			flags: Flags{
				SliceSize:      -1,
				ShowTop:        10,
				MaxWorkers:     1,
				MaxThreads:     1,
				ProgressUpdate: 1,
				Slices:         4,
				SliceThreshold: 102400,
			},
			wantError: true,
			errorMsg:  "slice size and threshold must be non-negative",
		},
		{
			name: "invalid slice count",
			flags: Flags{
				Slices:         0,
				ShowTop:        10,
				MaxWorkers:     1,
				MaxThreads:     1,
				ProgressUpdate: 1,
			},
			wantError: true,
			errorMsg:  `defaultSlices cannot be less than '\x04'`,
		},
		{
			name: "negative min size",
			flags: Flags{
				MinSize:        -1,
				ShowTop:        10,
				MaxWorkers:     1,
				MaxThreads:     1,
				ProgressUpdate: 1,
				Slices:         4,
				SliceSize:      8192,
				SliceThreshold: 102400,
			},
			wantError: true,
			errorMsg:  "min size and max size must be non-negative",
		},
		{
			name: "verbose and silent conflict",
			flags: Flags{
				Verbose: true,
				Silent:  true,
				ShowTop: 10,
			},
			wantError: true,
			errorMsg:  "cannot be verbose and silent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				Flags: &tt.flags,
			}

			err := app.validateArgs()

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAppSessionInitialisation(t *testing.T) {
	app := &App{
		Flags: &Flags{
			Algorithm:      int(algorithms.Xxhash),
			MaxWorkers:     4,
			MaxThreads:     4,
			SliceSize:      8192,
			SliceThreshold: 102400,
			Slices:         4,
		},
	}

	// Just test the session initialisation part
	err := app.Run()
	if err == nil {
		t.Error("expected error due to missing locations")
	}

	// Verify session was created
	if app.Session == nil {
		t.Error("expected session to be initialised")
	}

	if app.Session != nil {
		if app.Session.Dupes == nil {
			t.Error("expected Dupes map to be initialised")
		}
		if app.Session.Fails == nil {
			t.Error("expected Fails map to be initialised")
		}
		if app.Session.Empty == nil {
			t.Error("expected Empty to be initialised")
		}
	}
}
