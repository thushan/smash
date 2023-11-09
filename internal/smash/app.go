package smash

import (
	"log"
	"os"

	"github.com/logrusorgru/aurora/v3"
	"github.com/thushan/smash/internal/app"
	"github.com/thushan/smash/pkg/indexer"
)

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

	files := make(chan indexer.FileFS)

	go func() {
		defer close(files)
		for _, location := range locations {
			app.printVerbose("Indexing location ", aurora.Cyan(location))
			err := walker.WalkDirectory(os.DirFS(location), location, files)

			if err != nil {
				log.Println("Failed to walk location ", aurora.Magenta(location), " because ", aurora.Red(err))
			}
		}
	}()

	totalFiles := 0
	for file := range files {
		totalFiles++
		app.printVerbose("Indexed file ", aurora.Blue(file.Path))
	}

	app.printVerbose("Total Files: ", aurora.Blue(totalFiles))
	return nil
}
