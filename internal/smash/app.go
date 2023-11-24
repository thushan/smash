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
	Dupes     *haxmap.Map[string, []report.SmashFile]
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

	if !app.Flags.Silent {
		PrintVersionInfo(app.Flags.ShowVersion)
		if app.Flags.ShowVersion {
			return nil
		}
		app.printConfiguration()
	}

	app.Session = &AppSession{
		Dupes:     haxmap.New[string, []report.SmashFile](),
		Fails:     haxmap.New[string, error](),
		Empty:     &[]report.SmashFile{},
		StartTime: time.Now().UnixMilli(),
		EndTime:   -1,
	}

	sl := slicer.New(algorithms.Algorithm(app.Flags.Algorithm))
	wk := indexer.NewConfigured(app.Flags.ExcludeDir, app.Flags.ExcludeFile)
	slo := slicer.SlicerOptions{
		DisableSlicing:       app.Flags.DisableSlicing,
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

	session := app.Session

	wk := app.Runtime.IndexerConfig
	sl := app.Runtime.Slicer
	slo := app.Runtime.SlicerOptions

	files := app.Runtime.Files
	locations := app.Locations
	isVerbose := app.Flags.Verbose && !app.Flags.Silent

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
	updateTicker := int64(1000)

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

	pss.Success("Finding duplicates...Done!")

	psr, _ := theme.FinaliseSpinner().WithWriter(pap.NewWriter()).Start("Finding smash hits...")
	app.generateRunSummary(totalFiles)
	psr.Success("Finding smash hits...Done!")

	pap.Stop()

	app.PrintRunAnalysis(app.Flags.IgnoreEmptyFiles)
	report.PrintRunSummary(*app.Summary, app.Flags.IgnoreEmptyFiles)

	return nil
}

func (app *App) checkTerminal() {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		pterm.DisableColor()
		pterm.DisableStyling()
	}
}
