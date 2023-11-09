package smash

import (
	"encoding/hex"
	slicer2 "github.com/thushan/smash/pkg/slicer"
	"hash/fnv"
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

	if !app.Flags.Silent {
		app.printConfiguration()
	}

	app.setMaxThreads()

	files := make(chan indexer.FileFS)

	slicer := slicer2.New(fnv.New128a())
	walker := indexer.NewConfigured(excludeDirs, excludeFiles)

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
		app.printVerbose("Smashing file ", aurora.Blue(file.Path))
		hash, _ := slicer.SliceFS(file.FileSystem, file.Path)
		hashText := hex.EncodeToString(hash[:])
		app.printVerbose(" Hash: ", aurora.Magenta(hashText))
	}

	app.printVerbose("Total Files: ", aurora.Blue(totalFiles))
	return nil
}
