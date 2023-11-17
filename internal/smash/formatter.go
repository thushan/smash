package smash

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/internal/report"

	"github.com/thushan/smash/internal/theme"
)

const (
	TreeNextChild = "├─"
	TreeLastChild = "└─"
)

func (app *App) printVerbose(message ...any) {
	if app.Flags.Verbose {
		theme.Verbose.Println(message...)
	}
}

// generateSmashHits Generates the smash hits of duplicates and returns the total size of duplicates.
func (app *App) generateSmashHits(totalFiles int64) report.RunSummary {
	session := *app.Session
	duplicates := session.Dupes
	emptyFiles := *session.Empty

	totalDuplicates := 0
	totalUniqueFiles := int64(duplicates.Len())
	totalDuplicateSize := uint64(0)
	totalFailFileCount := int64(session.Fails.Len())
	totalEmptyFileCount := int64(len(emptyFiles))

	theme.StyleHeading.Println("---| Duplicates")
	duplicates.ForEach(func(hash string, files []report.SmashFile) bool {
		duplicateFiles := len(files) - 1
		if duplicateFiles == 0 {
			// prune unique files
			duplicates.Del(hash)
		} else {
			root := files[0]
			dupes := files[1:]
			theme.Println(theme.ColourFilename(root.Filename), " ", theme.ColourFileSize(humanize.Bytes(root.FileSize)), " ", theme.ColourHash(root.Hash))
			printSmashHits(dupes)
			totalDuplicates += len(dupes)
			totalDuplicateSize += root.FileSize * uint64(duplicateFiles)
		}
		return true
	})

	if totalDuplicates == 0 {
		theme.Println(theme.ColourSuccess("No duplicates found :-)"))
	}

	if !app.Flags.IgnoreEmptyFiles && totalEmptyFileCount != 0 {
		theme.StyleHeading.Println("---| Empty Files (", totalEmptyFileCount, ")")
		printSmashHits(emptyFiles)
	}

	return report.RunSummary{
		TotalFiles:         totalFiles,
		TotalFileErrors:    totalFailFileCount,
		UniqueFiles:        totalUniqueFiles,
		EmptyFiles:         totalEmptyFileCount,
		DuplicateFiles:     int64(totalDuplicates),
		DuplicateFileSize:  totalDuplicateSize,
		DuplicateFileSizeF: humanize.Bytes(totalDuplicateSize),
		ElapsedTime:        time.Now().UnixMilli() - app.Session.StartTime,
	}
}

func printSmashHits(files []report.SmashFile) {
	lastIndex := len(files) - 1
	for index, file := range files {
		var subTree string
		if index < lastIndex {
			subTree = TreeNextChild
		} else {
			subTree = TreeLastChild
		}
		theme.Println(theme.ColourFolderHierarchy(subTree), theme.ColourFilenameA(file.Filename))
	}
}
