package smash

import (
	"log"
	"runtime"
	"strings"

	. "github.com/logrusorgru/aurora/v3"
)

func (app *App) printConfiguration() {
	f := app.Flags
	log.Println(Bold(Cyan("Configuration")))
	log.Println(Bold("Locations:   "), Magenta(strings.Join(app.Locations, ", ")))
	log.Println(Bold("Max Threads: "), Magenta(f.MaxThreads))
	log.Println(Bold("Max Workers: "), Magenta(f.MaxWorkers))
	if len(f.ExcludeDir) > 0 || len(f.ExcludeFile) > 0 {
		log.Println(Bold("Excluded"))
		if len(f.ExcludeDir) > 0 {
			log.Println(Bold("       Dirs: "), Yellow(strings.Join(f.ExcludeDir, ", ")))
		}
		if len(f.ExcludeFile) > 0 {
			log.Println(Bold("      Files: "), Yellow(strings.Join(f.ExcludeFile, ", ")))
		}
	}
}

func (app *App) setMaxThreads() {
	maxThreads := app.Flags.MaxThreads
	if maxThreads < 1 || maxThreads > runtime.NumCPU() {
		return
	}
	runtime.GOMAXPROCS(maxThreads)
	app.printVerbose("Max Threads set to ", Magenta(runtime.GOMAXPROCS(0)))
}
