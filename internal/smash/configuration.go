package smash

import (
	"runtime"
	"strings"

	"github.com/thushan/smash/internal/theme"

	"github.com/thushan/smash/internal/algorithms"
)

func (app *App) printConfiguration() {
	f := app.Flags
	b := theme.StyleBold

	theme.StyleHeading.Println("---| Configuration")
	theme.Println(b.Sprint("Concurrency: "), theme.ColourConfig(f.MaxWorkers), "workers |", theme.ColourConfig(f.MaxThreads), "threads")
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

func (app *App) setMaxThreads() {
	maxThreads := app.Flags.MaxThreads
	if maxThreads < 1 || maxThreads > runtime.NumCPU() {
		return
	}
	runtime.GOMAXPROCS(maxThreads)
	app.printVerbose("Max Threads set to ", theme.ColourConfig(runtime.GOMAXPROCS(0)))
}
