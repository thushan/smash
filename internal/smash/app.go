package smash

import (
	"encoding/hex"
	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/slicer"
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

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
	var disableSlicing = true

	if !app.Flags.Silent {
		app.printConfiguration()
	}

	app.setMaxThreads()

	files := make(chan indexer.FileFS)

	sl := slicer.New(fnv.New128a())
	wk := indexer.NewConfigured(excludeDirs, excludeFiles)

	go func() {
		defer close(files)
		for _, location := range locations {
			app.printVerbose("Indexing location ", aurora.Cyan(location))
			err := wk.WalkDirectory(os.DirFS(location), location, files)

			if err != nil {
				log.Println("Failed to walk location ", aurora.Magenta(location), " because ", aurora.Red(err))
			}
		}
	}()

	/**
	 * Good times: https://go-review.googlesource.com/c/go/+/293349
	 */

	totalFiles := int32(0)
	var wg sync.WaitGroup
	for i := 0; i < app.Flags.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range files {
				atomic.AddInt32(&totalFiles, 1)
				startTime := time.Now().UnixMilli()
				hash, size, _ := sl.SliceFS(file.FileSystem, file.Path, disableSlicing)
				elapsedMs := time.Now().UnixMilli() - startTime
				hashText := hex.EncodeToString(hash[:])
				app.printVerbose("Smashed: ", aurora.Magenta(file.Path), aurora.Green(strconv.FormatInt(elapsedMs, 10)+"ms"))
				app.printVerbose("   Size: ", aurora.Cyan(humanize.Bytes(size)))
				app.printVerbose("   Hash: ", aurora.Blue(hashText))
			}
		}()
	}

	wg.Wait()
	app.printVerbose("Total Files: ", aurora.Blue(totalFiles))
	return nil
}
