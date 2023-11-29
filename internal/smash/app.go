package smash

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/term"

	"github.com/thushan/smash/internal/report"

	"github.com/pterm/pterm"
	"github.com/thushan/smash/internal/theme"

	"github.com/alphadose/haxmap"

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
	Dupes     *haxmap.Map[string, *report.SmashFiles]
	Fails     *haxmap.Map[string, error]
	Empty     *[]report.SmashFile
	StartTime int64
	EndTime   int64
}
type AppRuntime struct {
	Slicer        *slicer.Slicer
	SlicerOptions *slicer.SlicerOptions
	IndexerConfig *indexer.IndexerConfig
	Files         chan indexer.FileFS
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

	app.Session = &AppSession{
		Dupes:     haxmap.New[string, *report.SmashFiles](),
		Fails:     haxmap.New[string, error](),
		Empty:     &[]report.SmashFile{},
		StartTime: time.Now().UnixNano(),
		EndTime:   -1,
	}

	sl := slicer.New(algorithms.Algorithm(af.Algorithm))
	wk := indexer.NewConfigured(af.ExcludeDir, af.ExcludeFile, af.IgnoreHidden, af.IgnoreSystem)
	slo := slicer.SlicerOptions{
		DisableSlicing:       af.DisableSlicing,
		DisableMeta:          false, // TODO: Flag this
		DisableFileDetection: false, // TODO: Flag this
	}

	app.Runtime = &AppRuntime{
		Slicer:        &sl,
		SlicerOptions: &slo,
		IndexerConfig: wk,
		Files:         make(chan indexer.FileFS),
	}

	app.setMaxThreads()
	app.checkTerminal()

	return app.Exec()
}
func (app *App) Exec() error {

	if err := app.validateArgs(); err != nil {
		return err
	}
	startStats := report.ReadNerdStats()
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
				_, _ = session.Fails.GetOrSet(location, err)
			}
		}
	}()

	totalFiles := int64(0)

	pss, _ := theme.SmashingSpinner().WithWriter(pap.NewWriter()).Start("Finding duplicates...")

	var wg sync.WaitGroup

	updateProgressTicker := make(chan bool)

	if showProgress {
		app.updateDupeCount(updateProgressTicker, pss, &totalFiles)
	}

	for i := 0; i < app.Flags.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range files {
				sf := resolveFilename(file)

				atomic.AddInt64(&totalFiles, 1)

				startTime := time.Now().UnixMilli()
				stats, err := sl.SliceFS(file.FileSystem, file.Path, slo)
				elapsedMs := time.Now().UnixMilli() - startTime

				if err != nil {
					if isVerbose {
						theme.WarnSkipWithContext(file.FullName, err)
					}
					_, _ = session.Fails.GetOrSet(sf, err)
				} else {
					report.SummariseSmashedFile(stats, sf, elapsedMs, session.Dupes, session.Empty)
				}
			}
		}()
	}
	wg.Wait()

	// Signal we're done
	updateProgressTicker <- true

	pss.Success("Finding duplicates...Done!")

	psr, _ := theme.FinaliseSpinner().WithWriter(pap.NewWriter()).Start("Finding smash hits...")
	app.generateRunSummary(totalFiles)
	psr.Success("Finding smash hits...Done!")

	pap.Stop()

	endStats := report.ReadNerdStats()

	app.PrintRunAnalysis(app.Flags.IgnoreEmpty)
	report.PrintRunSummary(*app.Summary, app.Flags.IgnoreEmpty)

	if app.Flags.ShowNerdStats {
		theme.StyleHeading.Println("---| Nerd Stats")
		report.PrintNerdStats(startStats, "Commenced analysis")
		report.PrintNerdStats(endStats, "Completed analysis")
	}

	return nil
}

func (app *App) updateDupeCount(updateProgressTicker chan bool, pss *pterm.SpinnerPrinter, totalFiles *int64) {
	if app.Flags.HideProgress {
		return
	}
	go func() {
		ticker := time.Tick(time.Duration(app.Flags.ProgressUpdate) * time.Second)
		for {
			select {
			case <-ticker:
				latestFileCount := atomic.LoadInt64(totalFiles)
				pss.UpdateText(fmt.Sprintf("Finding duplicates... (%s files smash'd)", pterm.Gray(latestFileCount)))
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
