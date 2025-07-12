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
	Args      []string
	Locations []indexer.LocationFS
	Output    *OutputManager
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
	session := app.Session

	wk := app.Runtime.IndexerConfig
	sl := app.Runtime.Slicer
	slo := app.Runtime.SlicerOptions

	files := app.Runtime.Files
	locations := app.Locations
	isVerbose := app.Output.IsVerbose()
	walkOptions := indexer.WalkConfig{Recurse: app.Flags.Recurse}

	var pap *pterm.MultiPrinter
	if app.Output.ShouldShowProgress() {
		mp := theme.MultiWriter()
		pap = &mp
		pap.Start()
	}
	
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

	totalFiles := xsync.NewCounter()

	pss := app.Output.StartSpinner(theme.SmashingSpinner(), "Finding duplicates...", pap)

	var wg sync.WaitGroup

	updateProgressTicker := make(chan bool)

	if app.Output.ShouldShowProgress() {
		app.updateDupeCount(updateProgressTicker, pss, totalFiles)
	}

	for i := 0; i < app.Flags.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range files {
				totalFiles.Inc()

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
		}()
	}
	wg.Wait()

	// Signal we're done
	if app.Output.ShouldShowProgress() {
		updateProgressTicker <- true
	}
	app.Session.EndTime = time.Now().UnixNano()

	midStats := nerdstats.Snapshot()

	pss.Success("Finding duplicates...Done!")

	psr := app.Output.StartSpinner(theme.FinaliseSpinner(), "Finding smash hits...", pap)
	app.generateRunSummary(totalFiles.Value())
	psr.Success("Finding smash hits...Done!")

	if pap != nil {
		pap.Stop()
	}

	app.PrintRunAnalysis(app.Flags.IgnoreEmpty)

	exportStats := nerdstats.Snapshot()

	app.ExportReport()

	endStats := nerdstats.Snapshot()

	if !app.Output.IsSilent() {
		PrintRunSummary(*app.Summary, app.Flags)
	}

	if app.Flags.ShowNerdStats && !app.Output.IsSilent() {
		theme.StyleHeading.Println("---| Nerd Stats")
		PrintNerdStats(startStats, "> Initial")
		PrintNerdStats(midStats, "> Post-Analysis")
		PrintNerdStats(exportStats, "> Post-Summary")
		PrintNerdStats(endStats, "> Post-Report")
	}
	return nil
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
