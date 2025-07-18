package smash

import (
	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/analysis"

	"github.com/thushan/smash/internal/theme"
)

const (
	TreeNextChild = "├─"
	TreeLastChild = "└─"
)

func (app *App) printVerbose(message ...any) {
	if app.Output.IsVerbose() {
		theme.Verbose.Println(message...)
	}
}

func (app *App) PrintRunAnalysis(ignoreEmptyFiles bool) {
	if app.Output.IsSilent() {
		return
	}

	duplicates := app.Session.Dupes
	emptyFiles := app.Session.Empty.Files
	topFiles := app.Summary.TopFiles

	totalDuplicates := app.Summary.DuplicateFiles

	theme.StyleHeading.Println("---| Duplicate Files (", totalDuplicates, ")")

	if duplicates.Size() == 0 || len(topFiles) == 0 {
		theme.Println(theme.ColourSuccess("No duplicates found :-)"))
	} else {

		if !app.Flags.HideTopList {
			theme.StyleSubHeading.Println("---[ Top ", app.Flags.ShowTop, " Duplicates ]---")
			for _, tf := range topFiles {
				if files, ok := duplicates.Load(tf.Key); ok {
					displayFiles(files.Files)
				}
			}
		}

		if app.Flags.ShowDuplicates {
			theme.StyleSubHeading.Println("---[ All Duplicates ]---")
			duplicates.Range(func(hash string, files *DuplicateFiles) bool {
				displayFiles(files.Files)
				return true
			})
		}
	}

	if !ignoreEmptyFiles && len(emptyFiles) != 0 {
		theme.StyleHeading.Println("---| Empty Files (", len(emptyFiles), ")")
		printSmashHits(emptyFiles)
	}

}

func displayFiles(files []File) {
	duplicateFiles := len(files) - 1
	if duplicateFiles != 0 {
		root := files[0]
		dupes := files[1:]
		var dupeSize string
		if len(files) > 2 {
			// Since len(files) > 2, we know len(files) >= 3
			// Calculate duplicate count safely
			duplicateCount := len(files) - 1
			if duplicateCount > 0 && duplicateCount < len(files) {
				// #nosec G115 -- duplicateCount is guaranteed positive by the checks above
				fileCount := uint64(duplicateCount)
				totalDupeSize := fileCount * root.FileSize
				dupeSize = "(" + theme.ColourFileSizeDupe(humanize.Bytes(totalDupeSize)) + ")"
			} else {
				dupeSize = " "
			}
		} else {
			dupeSize = " "
		}
		theme.Println(theme.ColourFilename(root.Path), " ", theme.ColourFileSize(root.FileSizeF), dupeSize, theme.ColourHash(root.Hash))
		printSmashHits(dupes)
	}
}

func printSmashHits(files []File) {
	lastIndex := len(files) - 1
	for index, file := range files {
		var subTree string
		if index < lastIndex {
			subTree = TreeNextChild
		} else {
			subTree = TreeLastChild
		}
		theme.Println(theme.ColourFolderHierarchy(subTree), theme.ColourFilenameA(file.Path))
	}
}

// generateRunSummary Generates the smash hits of duplicates and returns the total size of duplicates.
func (app *App) generateRunSummary(totalFiles int64) {
	session := *app.Session
	duplicates := session.Dupes
	emptyFiles := session.Empty.Files

	topFiles := analysis.NewSummary(app.Flags.ShowTop)

	totalDuplicates := 0
	totalUniqueFiles := int64(duplicates.Size())
	totalDuplicateSize := uint64(0)
	totalFailFileCount := int64(session.Fails.Size())
	totalEmptyFileCount := int64(len(emptyFiles))

	duplicates.Range(func(hash string, df *DuplicateFiles) bool {
		files := df.Files
		duplicateFiles := len(files) - 1
		if duplicateFiles == 0 {
			// prune unique files
			duplicates.Delete(hash)
		} else {
			root := files[0]

			topFiles.Add(analysis.Item{Key: hash, Size: root.FileSize})

			totalDuplicates += duplicateFiles
			if duplicateFiles >= 0 {
				totalDuplicateSize += root.FileSize * uint64(duplicateFiles)
			}
		}
		return true
	})
	summary := RunSummary{
		TopFiles:           topFiles.All(),
		TotalFiles:         totalFiles,
		TotalFileErrors:    totalFailFileCount,
		UniqueFiles:        totalUniqueFiles,
		EmptyFiles:         totalEmptyFileCount,
		DuplicateFiles:     int64(totalDuplicates),
		DuplicateFileSize:  totalDuplicateSize,
		DuplicateFileSizeF: humanize.Bytes(totalDuplicateSize),
		ElapsedTime:        app.Session.EndTime - app.Session.StartTime,
	}
	app.Summary = &summary
}
