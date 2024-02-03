package smash

import (
	"testing"
)

func TestApp_ValidateArgs(t *testing.T) {
	tests := []struct {
		flags   *Flags
		name    string
		wantErr bool
	}{
		{
			name: "Should fail when verbose and silent are both true",
			flags: &Flags{
				Silent:  true,
				Verbose: true,
			},
			wantErr: true,
		},
		{
			name: "Should fail when maxThreads is below zero",
			flags: &Flags{
				MaxThreads: -1,
			},
			wantErr: true,
		},
		{
			name: "Should fail when maxWorkers is below zero",
			flags: &Flags{
				MaxWorkers: -1,
			},
			wantErr: true,
		},
		{
			name: "Should fail when showTop is equal or below 1",
			flags: &Flags{
				ShowTop: 1,
			},
			wantErr: true,
		},
		{
			name: "Should fail when showTop is not 10 and hidetop is set",
			flags: &Flags{
				ShowTop:     5,
				HideTopList: true,
			},
			wantErr: true,
		},
		{
			name: "Should fail when progressUpdate is below 1",
			flags: &Flags{
				ProgressUpdate: 0,
			},
			wantErr: true,
		},
		{
			name: "Should succeed when valid arguments are provided",
			flags: &Flags{
				Verbose:        true,
				MaxThreads:     5,
				MaxWorkers:     5,
				ShowTop:        10,
				ProgressUpdate: 2,
			},
			wantErr: false,
		},
		{
			name: "Should succeed when valid arguments are provided for min and max size",
			flags: &Flags{
				MinSize:        100,
				Verbose:        true,
				MaxThreads:     5,
				MaxWorkers:     5,
				ShowTop:        10,
				ProgressUpdate: 2,
			},
			wantErr: false,
		},
		{
			name: "Should succeed when valid arguments are provided for max size",
			flags: &Flags{
				MaxSize:        200,
				Verbose:        true,
				MaxThreads:     5,
				MaxWorkers:     5,
				ShowTop:        10,
				ProgressUpdate: 2,
			},
			wantErr: false,
		},
		{
			name: "Should succeed when valid arguments are provided for min size",
			flags: &Flags{
				MinSize:        200,
				MaxSize:        0,
				Verbose:        true,
				MaxThreads:     5,
				MaxWorkers:     5,
				ShowTop:        10,
				ProgressUpdate: 2,
			},
			wantErr: false,
		},
		{
			name: "Should fail when minSize is below zero",
			flags: &Flags{
				MinSize: -100,
			},
			wantErr: true,
		},
		{
			name: "Should fail when maxSize is below zero",
			flags: &Flags{
				MaxSize: -100,
			},
			wantErr: true,
		},
		{
			name: "Should fail when maxSize is below minSize",
			flags: &Flags{
				MaxSize: 20,
				MinSize: 200,
			},
			wantErr: true,
		},
		{
			name: "Should fail when maxSize is below minSize",
			flags: &Flags{
				MaxSize: 20,
				MinSize: 200,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				Flags: tt.flags,
			}

			if err := app.validateArgs(); (err != nil) != tt.wantErr {
				t.Errorf("App.validateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
