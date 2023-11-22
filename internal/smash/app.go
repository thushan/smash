package smash

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thushan/smash/internal/report"

	"golang.org/x/term"

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

func (app *App) Run() error {

	if !app.Flags.Silent {
		PrintVersionInfo(app.Flags.ShowVersion)
		if app.Flags.ShowVersion {
			return nil
		}
		app.printConfiguration()
	}

	var emptyFiles []report.SmashFile

	session := AppSession{
		Dupes:     haxmap.New[string, []report.SmashFile](),
		Fails:     haxmap.New[string, error](),
		Empty:     &emptyFiles,
		StartTime: time.Now().UnixMilli(),
		EndTime:   -1,
	}
	app.Session = &session
	app.setMaxThreads()

	locations := app.Locations
	excludeDirs := app.Flags.ExcludeDir
	excludeFiles := app.Flags.ExcludeFile
	disableSlicing := app.Flags.DisableSlicing
	isVerbose := app.Flags.Verbose && !app.Flags.Silent

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		pterm.DisableColor()
		pterm.DisableStyling()
	}

	files := make(chan indexer.FileFS)

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
				stats, err := sl.SliceFS(file.FileSystem, file.Path, &slo)
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
	pap.Stop()

	summary := app.generateRunSummary(totalFiles)
	report.PrintRunSummary(summary, app.Flags.IgnoreEmptyFiles)

	return nil
}
