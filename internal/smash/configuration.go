package smash

import (
	"runtime"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/internal/algorithms"
	"github.com/thushan/smash/internal/theme"
	"github.com/thushan/smash/pkg/indexer"
)

func (app *App) printConfiguration() {
	var config any
	f := app.Flags
	b := theme.StyleBold

	theme.StyleHeading.Println("---| Configuration")

	if app.Flags.Verbose {
		slices := theme.ColourConfig(f.Slices)
		size := theme.ColourConfig(humanize.Bytes(uint64(f.SliceSize)))
		threshold := theme.ColourConfig(humanize.Bytes(uint64(f.SliceThreshold)))

		config = "(Slices: " + slices + " | Size: " + size + " | Threshold: " + threshold + ")"

		maxThreads := theme.ColourConfig(f.MaxThreads)
		maxWorkers := theme.ColourConfig(f.MaxWorkers)

		theme.Println(b.Sprint("Concurrency: "), maxWorkers, "workers |", maxThreads, "threads")

	} else {
		config = ""
	}

	theme.Println(b.Sprint("Slicing:     "), theme.ColourConfig(enabledOrDisabled(!f.DisableSlicing)), config)
	theme.Println(b.Sprint("Algorithm:   "), theme.ColourConfig(algorithms.Algorithm(f.Algorithm)))
	theme.Println(b.Sprint("Locations:   "), theme.ColourConfig(buildLocations(app.Locations)))
	theme.Println(b.Sprint("Recursive:   "), theme.ColourConfig(enabledOrDisabled(f.Recurse)))

	if !f.HideOutput && f.OutputFile != "" {
		theme.Println(b.Sprint("Output:      "), theme.ColourConfig(f.OutputFile), "(json)")
	}

	if len(f.ExcludeDir) > 0 || len(f.ExcludeFile) > 0 {
		theme.StyleBold.Println("Excluded")
		if len(f.ExcludeDir) > 0 {
			theme.Println(b.Sprint("       Dirs: "), theme.ColourConfigA(strings.Join(f.ExcludeDir, ", ")))
		}
		if len(f.ExcludeFile) > 0 {
			theme.Println(b.Sprint("      Files: "), theme.ColourConfigA(strings.Join(f.ExcludeFile, ", ")))
		}
	}
}

func buildLocations(locations []indexer.LocationFS) string {
	locs := make([]string, len(locations))
	for i, location := range locations {
		locs[i] = location.Name
	}
	return strings.Join(locs, ", ")
}

func enabledOrDisabled(value bool) string {
	if value {
		return "Enabled"
	} else {
		return "Disabled"
	}
}

func (app *App) setMaxThreads() {
	maxThreads := app.Flags.MaxThreads
	if maxThreads < 1 || maxThreads > runtime.NumCPU() {
		maxThreads = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(maxThreads)
}
