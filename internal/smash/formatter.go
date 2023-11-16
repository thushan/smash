package smash

import (
	"fmt"
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

// printSmashHits Prints the smash hits of duplicates and returns the total size of duplicates.
func (app *App) printSmashHits() uint64 {
	duplicates := app.Session.Dupes
	totalDuplicateSize := uint64(0)
	totalDuplicateFileCount := duplicates.Len()

	theme.StyleHeading.Println("---| Duplicates (", totalDuplicateFileCount, ")")

	if totalDuplicateFileCount == 0 {
		theme.ColourSuccess("No duplicates found :-)")
	} else {
		duplicates.ForEach(func(hash string, files []report.SmashFile) bool {
			duplicateFiles := len(files) - 1
			if duplicateFiles == 0 {
				// prune unique files
				duplicates.Del(hash)
			} else {
				root := files[0]
				theme.Println(theme.ColourFilename(root.Filename), " ", theme.ColourFileSize(humanize.Bytes(root.FileSize)), " ", theme.ColourHash(root.Hash))
				printSmashHits(files[1:])
				totalDuplicateSize += root.FileSize * uint64(duplicateFiles)
			}
			return true
		})
	}

	emptyFiles := *app.Session.Empty
	totalEmptyFileCount := len(emptyFiles)

	if !app.Flags.IgnoreEmptyFiles && totalEmptyFileCount != 0 {
		theme.StyleHeading.Println("---| Empty Files (", totalEmptyFileCount, ")")
		printSmashHits(emptyFiles)
	}

	return totalDuplicateSize
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

func (app *App) printSmashRunSummary(rs report.RunSummary) {
	theme.StyleHeading.Println("---| Summary")

	theme.Println("Total Time:         ", theme.ColourTime(fmt.Sprintf("%dms", rs.ElapsedTime)))
	theme.Println("Total Files:        ", theme.ColourNumber(rs.TotalFiles))
	theme.Println("Total Unique:       ", theme.ColourNumber(rs.UniqueFiles))
	if rs.TotalFileErrors > 0 {
		theme.Println("Total Skipped:      ", theme.ColourError(rs.TotalFileErrors))
	}
	theme.Println("Total Duplicates:   ", theme.ColourNumber(rs.DuplicateFiles))
	if !app.Flags.IgnoreEmptyFiles && rs.EmptyFiles > 0 {
		theme.Println("Total Empty Files:  ", theme.ColourNumber(rs.EmptyFiles))
	}
	if rs.DuplicateFileSize > 0 {
		theme.Println("Space Reclaimable:  ", theme.ColourFileSizeA(rs.DuplicateFileSizeF), "(approx)")
	}

}
