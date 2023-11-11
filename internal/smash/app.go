package smash

import (
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/alphadose/haxmap"

	"github.com/thushan/smash/internal/algorithms"
	"github.com/thushan/smash/pkg/slicer"

	"github.com/logrusorgru/aurora/v3"
	"github.com/thushan/smash/pkg/indexer"
)

type Flags struct {
	Base           []string `yaml:"base"`
	ExcludeDir     []string `yaml:"exclude-dir"`
	ExcludeFile    []string `yaml:"exclude-file"`
	Algorithm      int      `yaml:"algorithm"`
	MaxThreads     int      `yaml:"max-threads"`
	MaxWorkers     int      `yaml:"max-workers"`
	DisableSlicing bool     `yaml:"disable-slicing"`
	Silent         bool     `yaml:"silent"`
	Verbose        bool     `yaml:"verbose"`
}

type App struct {
	Flags     *Flags
	Args      []string
	Locations []string
}
type SmashFile struct {
	Filename    string
	Hash        string
	FileSize    uint64
	ElapsedTime int64
	FullHash    bool
}

func (app *App) Run() error {

	var locations = app.Locations
	var excludeDirs = app.Flags.ExcludeDir
	var excludeFiles = app.Flags.ExcludeFile
	var disableSlicing = app.Flags.DisableSlicing

	if !app.Flags.Silent {
		PrintVersionInfo(false)
		app.printConfiguration()
	}

	app.setMaxThreads()

	files := make(chan indexer.FileFS)
	cache := haxmap.New[string, []SmashFile]()

	sl := slicer.New(algorithms.Algorithm(app.Flags.Algorithm))
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
				smashedFilename := resolveFilename(file)

				atomic.AddInt32(&totalFiles, 1)

				startTime := time.Now().UnixMilli()
				stats, err := sl.SliceFS(file.FileSystem, file.Path, disableSlicing)
				elapsedMs := time.Now().UnixMilli() - startTime

				app.printVerbose("Smashed: ", aurora.Magenta(smashedFilename), aurora.Green(strconv.FormatInt(elapsedMs, 10)+"ms"))

				if err != nil {
					app.printVerbose(" ERR: ", aurora.Red(err))
				} else {
					app.sumariseSmashedFile(cache, stats, smashedFilename, elapsedMs)
				}
			}
		}()
	}
	wg.Wait()
	app.printSmashHits(cache)
	log.Println("Total Files: ", aurora.Blue(totalFiles))
	log.Println("Total Unique: ", aurora.Blue(cache.Len()))
	return nil
}

func (app *App) sumariseSmashedFile(cache *haxmap.Map[string, []SmashFile], stats slicer.SlicerStats, filename string, ms int64) {
	smashedHash := hex.EncodeToString(stats.Hash)
	smashedFile := SmashFile{
		Filename:    filename,
		Hash:        smashedHash,
		FileSize:    stats.FileSize,
		FullHash:    false,
		ElapsedTime: ms,
	}
	if v, existing := cache.Get(smashedHash); existing {
		v = append(v, smashedFile)
		cache.Set(smashedHash, v)
	} else {
		cache.Set(smashedHash, []SmashFile{smashedFile})
	}
	app.printVerbose("   Size: ", aurora.Cyan(humanize.Bytes(stats.FileSize)))
	app.printVerbose("   Full: ", aurora.Blue(stats.HashedFullFile))
	app.printVerbose("   Hash: ", aurora.Blue(smashedHash))
}

func resolveFilename(file indexer.FileFS) string {
	if file.Path == "." {
		return filepath.Base(file.FullName)
	} else {
		return file.Path
	}
}
