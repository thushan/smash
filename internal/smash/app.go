package smash

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/thushan/smash/pkg/nerdstats"
	"github.com/thushan/smash/pkg/profiler"

	"github.com/puzpuzpuz/xsync/v4"

	"golang.org/x/term"

	"github.com/pterm/pterm"
	"github.com/thushan/smash/internal/theme"

	"github.com/thushan/smash/internal/algorithms"
	"github.com/thushan/smash/pkg/slicer"

	"github.com/thushan/smash/pkg/indexer"
)

type App struct {
	Flags     *Flags
	Session   *AppSession
	Runtime   *AppRuntime
	Summary   *RunSummary
	Output    *OutputManager
	Args      []string
	Locations []indexer.LocationFS
}
type AppSession struct {
	Dupes     *xsync.Map[string, *DuplicateFiles]
	Fails     *xsync.Map[string, error]
	Empty     *EmptyFiles
	StartTime int64
	EndTime   int64
}
type AppRuntime struct {
	Slicer        *slicer.Slicer
	SlicerOptions *slicer.Options
	IndexerConfig *indexer.IndexerConfig
	Files         chan *indexer.FileFS
}

const ReportOutputTemplate = "report-*.json"

func (app *App) Run() error {

	af := app.Flags

	// Initialize output manager first
	app.Output = NewOutputManager(af)

	if !af.Silent {
		PrintVersionInfo(af.ShowVersion)
		if af.ShowVersion {
			return nil
		}
		app.printConfiguration()
	}

	if af.Profile {
		profiler.InitialiseProfiler()
	}

	app.Session = &AppSession{
		Dupes: xsync.NewMap[string, *DuplicateFiles](),
		Fails: xsync.NewMap[string, error](),
		Empty: &EmptyFiles{
			Files:   []File{},
			RWMutex: sync.RWMutex{},
		},
		StartTime: time.Now().UnixNano(),
		EndTime:   -1,
	}

	// Validate and convert slice parameters
	if af.SliceSize < 0 || af.SliceThreshold < 0 {
		return fmt.Errorf("slice size and threshold must be non-negative")
	}
	if af.MinSize < 0 || af.MaxSize < 0 {
		return fmt.Errorf("min size and max size must be non-negative")
	}
	sl := slicer.NewConfigured(algorithms.Algorithm(af.Algorithm), af.Slices, uint64(af.SliceSize), uint64(af.SliceThreshold))
	wk := indexer.NewConfigured(af.ExcludeDir, af.ExcludeFile, af.IgnoreHidden, af.IgnoreSystem)
	slo := slicer.Options{
		DisableSlicing:  af.DisableSlicing,
		DisableMeta:     af.DisableMeta,
		DisableAutoText: af.DisableAutoText,
		MinSize:         uint64(af.MinSize),
		MaxSize:         uint64(af.MaxSize),
	}

	app.Runtime = &AppRuntime{
		Slicer:        &sl,
		SlicerOptions: &slo,
		IndexerConfig: wk,
		Files:         make(chan *indexer.FileFS),
	}

	app.setMaxThreads()
	app.checkTerminal()

	return app.Exec()
}
func (app *App) Exec() error {
	if err := app.validateArgs(); err != nil {
		return err
	}

	startStats := nerdstats.Snapshot()

	// Setup progress display
	pap := app.setupProgressDisplay()

	// Start indexing
	app.startIndexing(pap)

	// Process files
	totalFiles := app.processFiles(pap)

	// Finalize analysis
	app.finalizeAnalysis(pap, totalFiles)

	// Clean up progress display
	if pap != nil {
		pap.Stop()
	}

	// Print results and statistics
	app.printResultsAndStats(startStats)

	return nil
}

func (app *App) setupProgressDisplay() *pterm.MultiPrinter {
	if !app.Output.ShouldShowProgress() {
		return nil
	}
	mp := theme.MultiWriter()
	pap := &mp
	pap.Start()
	return pap
}

func (app *App) startIndexing(pap *pterm.MultiPrinter) {
	wk := app.Runtime.IndexerConfig
	files := app.Runtime.Files
	locations := app.Locations
	isVerbose := app.Output.IsVerbose()
	walkOptions := indexer.WalkConfig{Recurse: app.Flags.Recurse}
	session := app.Session

	psi := app.Output.StartSpinner(theme.IndexingSpinner(), "Indexing locations...", pap)
	go func() {
		defer func() {
			close(files)
			psi.Success("Indexing locations...Done!")
		}()
		for _, location := range locations {
			psi.UpdateText("Indexing location: " + location.Name)
			err := wk.WalkDirectory(location.FS, location.Name, walkOptions, files)

			if err != nil {
				if isVerbose {
					theme.WarnSkipWithContext(location.Name, err)
				}
				_, _ = session.Fails.LoadAndStore(location.Name, err)
			}
		}
	}()
}

func (app *App) processFiles(pap *pterm.MultiPrinter) int64 {
	sl := app.Runtime.Slicer
	slo := app.Runtime.SlicerOptions
	files := app.Runtime.Files
	isVerbose := app.Output.IsVerbose()
	session := app.Session

	totalFiles := xsync.NewCounter()
	pss := app.Output.StartSpinner(theme.SmashingSpinner(), "Finding duplicates...", pap)

	updateProgressTicker := make(chan bool)
	if app.Output.ShouldShowProgress() {
		app.updateDupeCount(updateProgressTicker, pss, totalFiles)
	}

	var wg sync.WaitGroup
	for i := 0; i < app.Flags.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range files {
				totalFiles.Inc()
				app.processFile(file, sl, slo, session, isVerbose)
			}
		}()
	}
	wg.Wait()

	if app.Output.ShouldShowProgress() {
		updateProgressTicker <- true
	}

	pss.Success("Finding duplicates...Done!")
	return totalFiles.Value()
}

func (app *App) processFile(file *indexer.FileFS, sl *slicer.Slicer, slo *slicer.Options, session *AppSession, isVerbose bool) {
	startTime := time.Now().UnixMilli()
	stats, err := sl.SliceFS(*file.FileSystem, file.Path, slo)
	elapsedMs := time.Now().UnixMilli() - startTime

	switch {
	case err != nil:
		if isVerbose {
			theme.WarnSkipWithContext(file.FullName, err)
		}
		_, _ = session.Fails.LoadOrStore(file.Path, err)
	case stats.IgnoredFile:
		// Check if it's an empty file that should be tracked
		if stats.EmptyFile {
			SummariseSmashedFile(stats, file, elapsedMs, session.Dupes, session.Empty)
		}
	default:
		SummariseSmashedFile(stats, file, elapsedMs, session.Dupes, session.Empty)
	}
}

func (app *App) finalizeAnalysis(pap *pterm.MultiPrinter, totalFiles int64) {
	app.Session.EndTime = time.Now().UnixNano()

	psr := app.Output.StartSpinner(theme.FinaliseSpinner(), "Finding smash hits...", pap)
	app.generateRunSummary(totalFiles)
	psr.Success("Finding smash hits...Done!")
}

func (app *App) printResultsAndStats(startStats nerdstats.NerdStats) {
	app.PrintRunAnalysis(app.Flags.IgnoreEmpty)

	midStats := nerdstats.Snapshot()
	app.ExportReport()
	exportStats := nerdstats.Snapshot()

	if !app.Output.IsSilent() {
		PrintRunSummary(*app.Summary, app.Flags)
	}

	endStats := nerdstats.Snapshot()

	if app.Flags.ShowNerdStats && !app.Output.IsSilent() {
		theme.StyleHeading.Println("---| Nerd Stats")
		PrintNerdStats(startStats, "> Initial")
		PrintNerdStats(midStats, "> Post-Analysis")
		PrintNerdStats(exportStats, "> Post-Summary")
		PrintNerdStats(endStats, "> Post-Report")
	}
}

func (app *App) updateDupeCount(updateProgressTicker chan bool, pss SpinnerHandle, totalFiles *xsync.Counter) {
	if !app.Output.ShouldShowProgress() {
		return
	}
	go func() {
		ticker := time.Tick(time.Duration(app.Flags.ProgressUpdate) * time.Second)
		for {
			select {
			case <-ticker:
				pss.UpdateText(fmt.Sprintf("Finding duplicates... (%s files smash'd)", pterm.Gray(totalFiles.Value())))
			case <-updateProgressTicker:
				return
			}
		}
	}()

}

func (app *App) checkTerminal() {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		pterm.DisableColor()
		pterm.DisableStyling()
	}
}

func (app *App) ExportReport() {
	if !app.Output.ShouldGenerateReport() {
		return
	}

	if filename, err := app.Export(app.Flags.OutputFile); err != nil {
		if !app.Output.IsSilent() {
			theme.Error.Println("Failed to export report because ", err)
		}
	} else {
		app.Summary.ReportFilename = filename
	}
}
