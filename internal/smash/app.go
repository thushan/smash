package smash

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/thushan/smash/pkg/nerdstats"
	"github.com/thushan/smash/pkg/profiler"

	"github.com/puzpuzpuz/xsync/v3"

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
}
type AppSession struct {
	Dupes     *xsync.MapOf[string, *DuplicateFiles]
	Fails     *xsync.MapOf[string, error]
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
		Dupes: xsync.NewMapOf[string, *DuplicateFiles](),
		Fails: xsync.NewMapOf[string, error](),
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
	isVerbose := app.Flags.Verbose && !app.Flags.Silent
	walkOptions := indexer.WalkConfig{Recurse: app.Flags.Recurse}
	showProgress := (!app.Flags.HideProgress && !app.Flags.Silent) || isVerbose

	pap := theme.MultiWriter()
	psi, _ := theme.IndexingSpinner().WithWriter(pap.NewWriter()).Start("Indexing locations...")

	pap.Start()
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

	pss, _ := theme.SmashingSpinner().WithWriter(pap.NewWriter()).Start("Finding duplicates...")

	var wg sync.WaitGroup

	updateProgressTicker := make(chan bool)

	if showProgress {
		app.updateDupeCount(updateProgressTicker, pss, *totalFiles)
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
					// Ignored counter
				default:
					SummariseSmashedFile(stats, file, elapsedMs, session.Dupes, session.Empty)
				}
			}
		}()
	}
	wg.Wait()

	// Signal we're done
	updateProgressTicker <- true
	app.Session.EndTime = time.Now().UnixNano()

	midStats := nerdstats.Snapshot()

	pss.Success("Finding duplicates...Done!")

	psr, _ := theme.FinaliseSpinner().WithWriter(pap.NewWriter()).Start("Finding smash hits...")
	app.generateRunSummary(totalFiles.Value())
	psr.Success("Finding smash hits...Done!")

	pap.Stop()

	app.PrintRunAnalysis(app.Flags.IgnoreEmpty)

	exportStats := nerdstats.Snapshot()

	app.ExportReport()

	endStats := nerdstats.Snapshot()

	PrintRunSummary(*app.Summary, app.Flags)

	if app.Flags.ShowNerdStats {
		theme.StyleHeading.Println("---| Nerd Stats")
		PrintNerdStats(startStats, "> Initial")
		PrintNerdStats(midStats, "> Post-Analysis")
		PrintNerdStats(exportStats, "> Post-Summary")
		PrintNerdStats(endStats, "> Post-Report")
	}
	return nil
}

func (app *App) updateDupeCount(updateProgressTicker chan bool, pss *pterm.SpinnerPrinter, totalFiles xsync.Counter) {
	if app.Flags.HideProgress {
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
	if app.Flags.HideOutput {
		return
	}

	if filename, err := app.Export(app.Flags.OutputFile); err != nil {
		theme.Error.Println("Failed to export report because ", err)
	} else {
		app.Summary.ReportFilename = filename
	}
}
