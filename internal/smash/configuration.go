package smash

import (
	"runtime"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/slicer"

	"github.com/thushan/smash/internal/theme"

	"github.com/thushan/smash/internal/algorithms"
)

func (app *App) printConfiguration() {
	var config any
	f := app.Flags
	b := theme.StyleBold

	theme.StyleHeading.Println("---| Configuration")

	if app.Flags.Verbose {
		theme.Println(b.Sprint("Concurrency: "), theme.ColourConfig(f.MaxWorkers), "workers |", theme.ColourConfig(f.MaxThreads), "threads")
		config = "(Slices: " + theme.ColourConfig(slicer.DefaultSlices) + " | Size: " + theme.ColourConfig(humanize.Bytes(slicer.DefaultSliceSize)) + " | Threshold: " + theme.ColourConfig(humanize.Bytes(slicer.DefaultThreshold)) + ")"
	} else {
		config = ""
	}

	theme.Println(b.Sprint("Slicing:     "), theme.ColourConfig(enabledOrDisabled(!f.DisableSlicing)), config)
	theme.Println(b.Sprint("Algorithm:   "), theme.ColourConfig(algorithms.Algorithm(f.Algorithm)))
	theme.Println(b.Sprint("Locations:   "), theme.ColourConfig(strings.Join(app.Locations, ", ")))

	if len(f.ExcludeDir) > 0 || len(f.ExcludeFile) > 0 {
		theme.StyleBold.Println(b.Sprint("Excluded"))
		if len(f.ExcludeDir) > 0 {
			theme.Println(b.Sprint("       Dirs: "), theme.ColourConfigA(strings.Join(f.ExcludeDir, ", ")))
		}
		if len(f.ExcludeFile) > 0 {
			theme.Println(b.Sprint("      Files: "), theme.ColourConfigA(strings.Join(f.ExcludeFile, ", ")))
		}
	}
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
