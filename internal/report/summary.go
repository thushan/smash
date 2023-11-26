package report

import (
	"fmt"
	"time"

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
	duration := time.Duration(rs.ElapsedTime)
	theme.Println(writeCategory("Total Time:"), theme.ColourTime(duration.Round(time.Second).String()))
	theme.Println(writeCategory("Total Analysed:"), theme.ColourNumber(rs.TotalFiles))
	theme.Println(writeCategory("Total Unique:"), theme.ColourNumber(rs.UniqueFiles), "(excludes empty files)")
	if rs.TotalFileErrors > 0 {
		theme.Println(writeCategory("Total Skipped:"), theme.ColourError(rs.TotalFileErrors))
	}
	theme.Println(writeCategory("Total Duplicates:"), theme.ColourNumber(rs.DuplicateFiles))
	if !ignoreEmptyFiles && rs.EmptyFiles > 0 {
		theme.Println(writeCategory("Total Empty Files:"), theme.ColourNumber(rs.EmptyFiles))
	}
	if rs.DuplicateFileSize > 0 {
		theme.Println(writeCategory("Space Reclaimable:"), theme.ColourFileSizeA(rs.DuplicateFileSizeF), "(approx)")
	}
}
func writeCategory(category string) string {
	return fmt.Sprintf("%20s", category)
}
