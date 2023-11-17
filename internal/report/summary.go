package report

import (
	"fmt"

	"github.com/thushan/smash/internal/theme"
)

type RunSummary struct {
	DuplicateFileSizeF string
	DuplicateFileSize  uint64
	TotalFiles         int64
	TotalFileErrors    int64
	ElapsedTime        int64
	UniqueFiles        int64
	EmptyFiles         int64
	DuplicateFiles     int64
}

func PrintRunSummary(rs RunSummary, ignoreEmptyFiles bool) {
	theme.StyleHeading.Println("---| Analysis Summary")

	theme.Println("Total Time:         ", theme.ColourTime(fmt.Sprintf("%dms", rs.ElapsedTime)))
	theme.Println("Total Analysed:     ", theme.ColourNumber(rs.TotalFiles))
	theme.Println("Total Unique:       ", theme.ColourNumber(rs.UniqueFiles), "(excludes empty files)")
	if rs.TotalFileErrors > 0 {
		theme.Println("Total Skipped:      ", theme.ColourError(rs.TotalFileErrors))
	}
	theme.Println("Total Duplicates:   ", theme.ColourNumber(rs.DuplicateFiles))
	if !ignoreEmptyFiles && rs.EmptyFiles > 0 {
		theme.Println("Total Empty Files:  ", theme.ColourNumber(rs.EmptyFiles))
	}
	if rs.DuplicateFileSize > 0 {
		theme.Println("Space Reclaimable:  ", theme.ColourFileSizeA(rs.DuplicateFileSizeF), "(approx)")
	}

}
