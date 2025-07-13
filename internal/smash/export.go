package smash

import (
	"encoding/json"
	"fmt"
	"os"
	user2 "os/user"
	"path/filepath"
	"time"

	"github.com/puzpuzpuz/xsync/v4"
	"github.com/thushan/smash/pkg/analysis"
)

type ReportOutput struct {
	Meta     ReportMeta    `json:"_meta"`
	Analysis ReportFiles   `json:"analysis"`
	Summary  ReportSummary `json:"summary"`
}
type ReportMeta struct {
	Timestamp time.Time `json:"timestamp"`
	Config    *Flags    `json:"config"`
	Version   string    `json:"version"`
	Commit    string    `json:"commit"`
	Host      string    `json:"host"`
	User      string    `json:"user"`
}
type ReportSummary struct {
	TopFiles          []ReportTopFilesSummary `json:"top"`
	DuplicateFileSize uint64                  `json:"duplicateFileSize"`
	TotalFiles        int64                   `json:"totalFiles"`
	TotalFileErrors   int64                   `json:"totalFileFails"`
	ElapsedTime       int64                   `json:"elapsedTime"`
	UniqueFiles       int64                   `json:"uniqueFiles"`
	EmptyFiles        int64                   `json:"emptyFiles"`
	DuplicateFiles    int64                   `json:"duplicateFiles"`
}
type ReportTopFilesSummary struct {
	Hash string `json:"hash"`
	Size uint64 `json:"size"`
}
type ReportFiles struct {
	Fails []ReportFailSummary      `json:"fails"`
	Empty []ReportFileBaseSummary  `json:"empty"`
	Dupes []ReportDuplicateSummary `json:"dupes"`
}

type ReportFailSummary struct {
	Filename string `json:"filename"`
	Error    string `json:"error"`
}

type ReportFileBaseSummary struct {
	Filename string `json:"filename"`
	Location string `json:"location"`
	Path     string `json:"path"`
}

type ReportFileSummary struct {
	ReportFileBaseSummary
	Hash     string `json:"hash"`
	Size     uint64 `json:"size"`
	FullHash bool   `json:"fullHash"`
}
type ReportDuplicateSummary struct {
	Duplicates []ReportFileSummary `json:"duplicates"`
	ReportFileSummary
}

func (app *App) Export(filePath string) (string, error) {

	var fs *os.File
	var err error

	if filePath == "" {
		if fs, err = os.CreateTemp(".", ReportOutputTemplate); err != nil {
			return "", fmt.Errorf("failed to report output: %w", err)
		}
	} else {
		if fs, err = os.Create(filePath); err != nil {
			return "", fmt.Errorf("failed to report output: %w", err)
		}
	}

	defer fs.Close()

	return fs.Name(), app.ExportFile(fs)
}

func (app *App) ExportFile(f *os.File) error {
	return json.NewEncoder(f).Encode(app.GenerateReportOutput())
}

func (app *App) GenerateReportOutput() ReportOutput {
	return ReportOutput{
		Summary:  summariseRunSummary(app.Summary),
		Analysis: summariseRunAnalysis(app.Session),
		Meta:     summariseMeta(app.Flags),
	}
}

func summariseMeta(flags *Flags) ReportMeta {
	return ReportMeta{
		Version:   Version,
		Commit:    Commit,
		Config:    flags,
		Host:      getHostName(),
		User:      getUsername(),
		Timestamp: time.Now(),
	}
}

func getUsername() string {
	if user, err := user2.Current(); err == nil {
		return user.Username
	}
	return "James Bond"
}

func getHostName() string {
	if host, err := os.Hostname(); err == nil {
		return host
	}
	return "Classified"
}

func summariseRunAnalysis(session *AppSession) ReportFiles {

	fails := summariseSmashFails(session.Fails)
	empty := summariseEmptyFiles(session.Empty.Files)
	dupes := transformDupes(session.Dupes)

	return ReportFiles{
		Fails: fails,
		Empty: empty,
		Dupes: dupes,
	}
}

func summariseSmashFails(fails *xsync.Map[string, error]) []ReportFailSummary {
	summary := make([]ReportFailSummary, fails.Size())
	var index = 0
	fails.Range(func(key string, value error) bool {
		summary[index] = ReportFailSummary{
			Filename: key,
			Error:    value.Error(),
		}
		index++
		return true
	})
	return summary
}

func transformDupes(duplicates *xsync.Map[string, *DuplicateFiles]) []ReportDuplicateSummary {
	dupes := make([]ReportDuplicateSummary, duplicates.Size())
	var index = 0
	duplicates.Range(func(hash string, dupe *DuplicateFiles) bool {
		root := dupe.Files[0]
		rest := dupe.Files[1:]
		dupes[index] = ReportDuplicateSummary{
			ReportFileSummary: summariseSmashedFile(root),
			Duplicates:        summariseSmashedFiles(rest),
		}
		index++
		return true
	})
	return dupes
}

func summariseEmptyFiles(files []File) []ReportFileBaseSummary {
	summary := make([]ReportFileBaseSummary, len(files))
	for i, file := range files {
		summary[i] = summariseSmashedFile(file).ReportFileBaseSummary
	}
	return summary
}
func summariseSmashedFiles(files []File) []ReportFileSummary {
	summary := make([]ReportFileSummary, len(files))
	for i, file := range files {
		summary[i] = summariseSmashedFile(file)
	}
	return summary
}
func summariseSmashedFile(file File) ReportFileSummary {
	return ReportFileSummary{
		ReportFileBaseSummary: ReportFileBaseSummary{
			Filename: file.Filename,
			Location: file.Location,
			Path:     filepath.Dir(file.Path),
		},
		Hash:     file.Hash,
		Size:     file.FileSize,
		FullHash: file.FullHash,
	}
}
func summariseRunSummary(summary *RunSummary) ReportSummary {
	return ReportSummary{
		TopFiles:          transformTopFiles(summary.TopFiles),
		DuplicateFileSize: summary.DuplicateFileSize,
		TotalFiles:        summary.TotalFiles,
		TotalFileErrors:   summary.TotalFileErrors,
		ElapsedTime:       summary.ElapsedTime,
		UniqueFiles:       summary.UniqueFiles,
		EmptyFiles:        summary.EmptyFiles,
		DuplicateFiles:    summary.DuplicateFiles,
	}
}

func transformTopFiles(files []analysis.Item) []ReportTopFilesSummary {
	items := make([]ReportTopFilesSummary, len(files))
	for i, file := range files {
		items[i] = ReportTopFilesSummary{
			Hash: file.Key,
			Size: file.Size,
		}
	}
	return items
}
