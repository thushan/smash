package smash

import (
	"io/fs"
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

	for _, location := range locations {
		app.printVerbose("Indexing location ", aurora.Cyan(location))
		files, _ := walker.WalkDirectory(buildLocations, done)
	}

	totalFiles := 0

	for file := range list {
		totalFiles++
		app.printVerbose("Indexed file ", aurora.Blue(file.Name))
	}

	return nil
}
func buildLocations(locations []string) []fs.FS {
	paths := make([]fs.FS, len(locations))

	for _, location := range locations {
		// we support local for now
		paths = append(paths, os.DirFS(location))
	}
	return paths
}
