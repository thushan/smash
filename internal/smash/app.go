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

	list := make(chan indexer.FileFS)
	done := make(chan struct{})

	defer close(done)

	go func() {
		defer close(list)
		for _, location := range locations {
			app.printVerbose("Indexing location ", aurora.Cyan(location))
			errc := walker.WalkDirectory(os.DirFS(location), list, done)

			if err := <-errc; err != nil {
				log.Println("Failed to walk location ", aurora.Magenta(location), " because ", aurora.Red(errc))
			}
		}
	}()

	totalFiles := 0
	for file := range list {
		totalFiles++
		app.printVerbose("Indexed file ", aurora.Blue(file.Name))
	}

	app.printVerbose("Total Files: ", aurora.Blue(totalFiles))
	return nil
}
