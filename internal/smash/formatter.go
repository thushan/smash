package smash

import (
	"fmt"

	"github.com/thushan/smash/internal/theme"

	"github.com/alphadose/haxmap"
	"github.com/dustin/go-humanize"
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

func (app *App) printSmashHits(cache *haxmap.Map[string, []SmashFile]) uint64 {
	totalDuplicateSize := uint64(0)
	theme.StyleHeading.Println("---| Duplicates")
	cache.ForEach(func(hash string, files []SmashFile) bool {
		mainFile := files[0]
		lastIndex := len(files)
		if lastIndex > 1 {
			theme.Println(theme.ColourFilename(mainFile.Filename), " ", theme.ColourFileSize(humanize.Bytes(mainFile.FileSize)), " ", theme.ColourHash(mainFile.Hash))
			for index, file := range files[1:] {
				var subTree string
				if (index + 2) == lastIndex {
					subTree = TreeLastChild
				} else {
					subTree = TreeNextChild
				}
				theme.Println(theme.ColourFolderHierarchy(subTree), theme.ColourFilenameA(file.Filename))
			}
			totalDuplicateSize += mainFile.FileSize * uint64(lastIndex-1)
		} else {
			// prune unique files
			cache.Del(hash)
		}
		return true
	})
	if cache.Len() == 0 {
		theme.ColourSuccess("No duplicates found :-)")
	}
	return totalDuplicateSize
}

func (app *App) printSmashRunSummary(rs RunSummary) {
	theme.StyleHeading.Println("---| Summary")

	theme.Println("Total Time:         ", theme.ColourTime(fmt.Sprintf("%dms", rs.ElapsedTime)))
	theme.Println("Total Files:        ", theme.ColourNumber(rs.TotalFiles))
	theme.Println("Total Unique:       ", theme.ColourNumber(rs.UniqueFiles))
	if rs.TotalFileErrors > 0 {
		theme.Println("Total Skipped:      ", theme.ColourError(rs.TotalFileErrors))
	}
	theme.Println("Total Duplicates:   ", theme.ColourNumber(rs.DuplicateFiles))
	if rs.DuplicateFileSize > 0 {
		theme.Println("Approx Reclaimable: ", theme.ColourFileSizeA(rs.DuplicateFileSizeF))
	}
	if rs.EmptyFiles > 0 {
		theme.Println("Total Empty Files:  ", theme.ColourNumber(rs.EmptyFiles))
	}

}
