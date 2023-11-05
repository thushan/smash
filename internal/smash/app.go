package smash

import (
	"os"

	"github.com/logrusorgru/aurora/v3"
	"github.com/thushan/smash/internal/app"
	"github.com/thushan/smash/pkg/indexer"
)

var FileQueueSize = 1000

type App struct {
	Flags     *app.Flags
	Args      []string
	Locations []string
}

func (app *App) Run() error {
	var locations = app.Locations
	var excludeDirs = app.Flags.ExcludeDir
	var excludeFiles = app.Flags.ExcludeFile
	var walker = indexer.NewConfigured(excludeDirs, excludeFiles)

	if !app.Flags.Silent {
		app.printConfiguration()
	}

	app.setMaxThreads()

	fsq := make(chan string, FileQueueSize)

	go func() {
		for _, location := range locations {
			app.printVerbose("Indexing location ", aurora.Cyan(location))
			walker.WalkDirectory(os.DirFS(location), fsq)
		}
		close(fsq)
	}()

	totalFiles := 0

	for filename := range fsq {
		totalFiles++
		app.printVerbose("Indexed file ", aurora.Blue(filename))
	}

	return nil
}
