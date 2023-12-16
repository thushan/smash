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

	"github.com/thushan/smash/internal/report"

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
	Summary   *report.RunSummary
	Args      []string
	Locations []string
}
type AppSession struct {
	Dupes     *xsync.MapOf[string, *report.DuplicateFiles]
	Fails     *xsync.MapOf[string, error]
	Empty     *report.EmptyFiles
	StartTime int64
	EndTime   int64
}
type AppRuntime struct {
	Slicer        *slicer.Slicer
	SlicerOptions *slicer.Options
	IndexerConfig *indexer.IndexerConfig
	Files         chan *indexer.FileFS
}

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
		Dupes: xsync.NewMapOf[string, *report.DuplicateFiles](),
		Fails: xsync.NewMapOf[string, error](),
		Empty: &report.EmptyFiles{
			Files:   []report.SmashFile{},
			RWMutex: sync.RWMutex{},
		},
		StartTime: time.Now().UnixNano(),
		EndTime:   -1,
	}

	sl := slicer.New(algorithms.Algorithm(af.Algorithm))
	wk := indexer.NewConfigured(af.ExcludeDir, af.ExcludeFile, af.IgnoreHidden, af.IgnoreSystem)
	slo := slicer.Options{
		DisableSlicing:  af.DisableSlicing,
		DisableMeta:     af.DisableMeta,
		DisableAutoText: af.DisableAutoText,
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
			psi.UpdateText("Indexing location: " + location)
			err := wk.WalkDirectory(os.DirFS(location), location, files)

			if err != nil {
				if isVerbose {
					theme.WarnSkipWithContext(location, err)
				}
				_, _ = session.Fails.LoadAndStore(location, err)
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

				if err != nil {
					if isVerbose {
						theme.WarnSkipWithContext(file.FullName, err)
					}
					_, _ = session.Fails.LoadOrStore(file.Path, err)
				} else {
					report.SummariseSmashedFile(stats, file, elapsedMs, session.Dupes, session.Empty)
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

	report.PrintRunSummary(*app.Summary, app.Flags.IgnoreEmpty)

	if app.Flags.ShowNerdStats {
		theme.StyleHeading.Println("---| Nerd Stats")
		report.PrintNerdStats(startStats, "> Initial")
		report.PrintNerdStats(midStats, "> Post-Analysis")
		report.PrintNerdStats(exportStats, "> Post-Summary")
		report.PrintNerdStats(endStats, "> Post-Report")
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
	if app.Flags.OutputFile == "" {
		theme.Warn.Println("Could not output report.")
		return
	}

	if err := app.Export(app.Flags.OutputFile); err != nil {
		theme.Error.Println("Failed to export report because ", err)
	} else {
		app.Summary.ReportFilename = app.Flags.OutputFile
	}
}
