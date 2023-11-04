package smash

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/thushan/smash/internal/app"
	"github.com/thushan/smash/pkg/indexer"
	"os"
)

var FileQueueSize = 1000

type App struct {
	Args      []string
	Locations []string
	Flags     *app.Flags
}

func (app *App) Run() error {

	var locations = app.Locations
	var excludeDirs = app.Flags.ExcludeDir
	var excludeFiles = app.Flags.ExcludeFile

	if !app.Flags.Silent {
		app.printConfiguration()
	}

	app.setMaxThreads()

	fsq := make(chan string, FileQueueSize)

	go func() {
		for _, location := range locations {
			app.printVerbose("Indexing location ", aurora.Cyan(location))
			indexer := indexer.NewConfigured(excludeDirs, excludeFiles)
			indexer.WalkDirectory(os.DirFS(location), fsq)
		}
		close(fsq)
	}()

	var totalFiles = 0

	for filename := range fsq {
		totalFiles++
		app.printVerbose("Indexed file ", aurora.Blue(filename))
	}

	return nil
}
