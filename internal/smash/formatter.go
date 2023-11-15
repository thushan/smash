package smash

import (
	"encoding/hex"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/slicer"

	"github.com/thushan/smash/internal/theme"

	"github.com/alphadose/haxmap"
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
	emptyFileHash := hex.EncodeToString(slicer.DefaultEmptyFileCookie)
	theme.StyleHeading.Println("---| Duplicates")
	cache.ForEach(func(hash string, files []SmashFile) bool {
		if hash == emptyFileHash {
			// Skip this for now
			return true
		}
		lastIndex := len(files)
		if lastIndex > 1 {
			root := files[0]
			printSmashHit(root, files, lastIndex)
			totalDuplicateSize += root.FileSize * uint64(lastIndex-1)
		} else {
			// prune unique files
			cache.Del(hash)
		}
		return true
	})
	if cache.Len() == 0 {
		theme.ColourSuccess("No duplicates found :-)")
	}
	if !app.Flags.IgnoreEmptyFiles {
		if files, ok := cache.Get(emptyFileHash); ok {
			theme.StyleHeading.Println("---| Empty Files")
			root := files[0]
			printSmashHit(root, files, len(files))
		}
	}

	return totalDuplicateSize
}

func printSmashHit(root SmashFile, duplicates []SmashFile, lastIndex int) {
	theme.Println(theme.ColourFilename(root.Filename), " ", theme.ColourFileSize(humanize.Bytes(root.FileSize)), " ", theme.ColourHash(root.Hash))
	for index, file := range duplicates[1:] {
		var subTree string
		if (index + 2) == lastIndex {
			subTree = TreeLastChild
		} else {
			subTree = TreeNextChild
		}
		theme.Println(theme.ColourFolderHierarchy(subTree), theme.ColourFilenameA(file.Filename))
	}
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
	if !app.Flags.IgnoreEmptyFiles && rs.EmptyFiles > 0 {
		theme.Println("Total Empty Files:  ", theme.ColourNumber(rs.EmptyFiles))
	}
	if rs.DuplicateFileSize > 0 {
		theme.Println("Space Reclaimable:  ", theme.ColourFileSizeA(rs.DuplicateFileSizeF), "(approx)")
	}

}
