package smash

import (
	"fmt"
	"github.com/thushan/smash/internal/report"
	"os"
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
	cache := haxmap.New[string, []report.SmashFile]()
	fails := haxmap.New[string, error]()

	sl := slicer.New(algorithms.Algorithm(app.Flags.Algorithm))
	wk := indexer.NewConfigured(excludeDirs, excludeFiles)

	slo := slicer.SlicerOptions{
		DisableSlicing:       disableSlicing,
		DisableMeta:          false, // TODO: Flag this
		DisableFileDetection: false, // TODO: Flag this
	}
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
				stats, err := sl.SliceFS(file.FileSystem, file.Path, &slo)
				elapsedMs := time.Now().UnixMilli() - startTime

				if err != nil {
					if app.Flags.Verbose {
						theme.WarnSkipWithContext(file.FullName, err)
					}
					_, _ = fails.GetOrSet(sf, err)
				} else {
					report.SummariseSmashedFile(cache, stats, sf, elapsedMs)
				}
			}
		}()
	}
	wg.Wait()

	pss.Success("Finding duplicates...Done!")

	psr, _ := theme.FinaliseSpinner().WithWriter(pap.NewWriter()).Start("Finding smash hits...")

	psr.Success("Finding smash hits...Done!")
	pap.Stop()

	summary := report.CalculateRunSummary(cache, fails, totalFiles, appStartTime)

	totalDuplicateSize := app.printSmashHits(cache)

	summary.DuplicateFileSize = totalDuplicateSize
	summary.DuplicateFileSizeF = humanize.Bytes(totalDuplicateSize)

	app.printSmashRunSummary(summary)

	return nil
}
