package smash

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/term"

	"github.com/dustin/go-humanize"
	"github.com/pterm/pterm"
	"github.com/thushan/smash/internal/theme"

	"github.com/alphadose/haxmap"

	"github.com/thushan/smash/internal/algorithms"
	"github.com/thushan/smash/pkg/slicer"

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
type RunSummary struct {
	DuplicateFileSizeF string
	DuplicateFileSize  uint64
	TotalFiles         int64
	TotalFileErrors    int64
	ElapsedTime        int64
	UniqueFiles        int64
	DuplicateFiles     int64
}
type SmashFile struct {
	Filename    string
	Hash        string
	FileSizeF   string
	FileSize    uint64
	ElapsedTime int64
	FullHash    bool
}

func (app *App) Run() error {

	var locations = app.Locations
	var excludeDirs = app.Flags.ExcludeDir
	var excludeFiles = app.Flags.ExcludeFile
	var disableSlicing = app.Flags.DisableSlicing

	/**
	 * Good times: https://go-review.googlesource.com/c/go/+/293349
	 */
	appStartTime := time.Now().UnixMilli()
	updateTicker := int64(1000)

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		pterm.DisableColor()
		pterm.DisableStyling()
	}

	if !app.Flags.Silent {
		PrintVersionInfo(false)
		app.printConfiguration()
	}

	app.setMaxThreads()

	files := make(chan indexer.FileFS)
	cache := haxmap.New[string, []SmashFile]()
	fails := haxmap.New[string, error]()

	sl := slicer.New(algorithms.Algorithm(app.Flags.Algorithm))
	wk := indexer.NewConfigured(excludeDirs, excludeFiles)

	pap := theme.MultiWriter()
	psi, _ := theme.IndexingSpinner().WithWriter(pap.NewWriter()).Start("Indexing locations...")

	pap.Start()
	go func() {
		defer func() {
			close(files)
			psi.Success("Indexing locations...Done!")
		}()
		for _, location := range locations {
			psi.UpdateText("Indexing location: " + location)
			err := wk.WalkDirectory(os.DirFS(location), location, files)

			if err != nil {
				if app.Flags.Verbose {
					theme.WarnSkipWithContext(location, err)
				}
				_, _ = fails.GetOrSet(location, err)
			}
		}
	}()

	totalFiles := int64(0)

	pss, _ := theme.SmashingSpinner().WithWriter(pap.NewWriter()).Start("Finding duplicates...")

	var wg sync.WaitGroup
	for i := 0; i < app.Flags.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range files {
				sf := resolveFilename(file)

				currentFileCount := atomic.AddInt64(&totalFiles, 1)

				if currentFileCount%updateTicker == 0 {
					pss.UpdateText(fmt.Sprintf("Finding duplicates... (%s files smash'd)", pterm.Gray(currentFileCount)))
				}

				startTime := time.Now().UnixMilli()
				stats, err := sl.SliceFS(file.FileSystem, file.Path, disableSlicing)
				elapsedMs := time.Now().UnixMilli() - startTime

				if err != nil {
					if app.Flags.Verbose {
						theme.WarnSkipWithContext(file.FullName, err)
					}
					_, _ = fails.GetOrSet(sf, err)
				} else {
					app.summariseSmashedFile(cache, stats, sf, elapsedMs)
				}
			}
		}()
	}
	wg.Wait()

	pss.Success("Finding duplicates...Done!")

	psr, _ := theme.FinaliseSpinner().WithWriter(pap.NewWriter()).Start("Finding smash hits...")

	summary := RunSummary{
		TotalFiles:      totalFiles,
		TotalFileErrors: int64(fails.Len()),
		UniqueFiles:     int64(cache.Len()),
		DuplicateFiles:  totalFiles - int64(cache.Len()),
		ElapsedTime:     time.Now().UnixMilli() - appStartTime,
	}

	psr.Success("Finding smash hits...Done!")
	pap.Stop()

	totalDuplicateSize := app.printSmashHits(cache)
	summary.DuplicateFileSize = totalDuplicateSize
	summary.DuplicateFileSizeF = humanize.Bytes(totalDuplicateSize)

	app.printSmashRunSummary(summary)

	return nil
}

func (app *App) summariseSmashedFile(cache *haxmap.Map[string, []SmashFile], stats slicer.SlicerStats, filename string, ms int64) {
	sf := SmashFile{
		Filename:    filename,
		Hash:        hex.EncodeToString(stats.Hash),
		FileSize:    stats.FileSize,
		FullHash:    stats.HashedFullFile,
		FileSizeF:   humanize.Bytes(stats.FileSize),
		ElapsedTime: ms,
	}
	if v, existing := cache.Get(sf.Hash); existing {
		v = append(v, sf)
		cache.Set(sf.Hash, v)
	} else {
		cache.Set(sf.Hash, []SmashFile{sf})
	}
}

func resolveFilename(file indexer.FileFS) string {
	if file.Path == "." {
		return filepath.Base(file.FullName)
	} else {
		return file.Path
	}
}
