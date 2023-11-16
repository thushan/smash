package report

import (
	"encoding/hex"
	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/slicer"
	"time"

	"github.com/alphadose/haxmap"
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

type SmashFile struct {
	Filename    string
	Hash        string
	FileSizeF   string
	FileSize    uint64
	ElapsedTime int64
	FullHash    bool
	EmptyFile   bool
}

func CalculateRunSummary(duplicates *haxmap.Map[string, []SmashFile], fails *haxmap.Map[string, error], emptyFiles *[]SmashFile, totalFiles int64, appStartTime int64) RunSummary {

	emptyFileCount := len(*emptyFiles)
	return RunSummary{
		TotalFiles:      totalFiles,
		TotalFileErrors: int64(fails.Len()),
		UniqueFiles:     int64(duplicates.Len()),
		EmptyFiles:      int64(emptyFileCount),
		DuplicateFiles:  totalFiles - (int64(duplicates.Len()) + int64(emptyFileCount)),
		ElapsedTime:     time.Now().UnixMilli() - appStartTime,
	}
}

func SummariseSmashedFile(stats slicer.SlicerStats, filename string, ms int64, duplicates *haxmap.Map[string, []SmashFile], emptyFiles *[]SmashFile) {
	sf := SmashFile{
		Hash:        hex.EncodeToString(stats.Hash),
		Filename:    filename,
		FileSize:    stats.FileSize,
		FullHash:    stats.HashedFullFile,
		EmptyFile:   stats.EmptyFile,
		FileSizeF:   humanize.Bytes(stats.FileSize),
		ElapsedTime: ms,
	}
	if sf.EmptyFile {
		*emptyFiles = append(*emptyFiles, sf)
	} else {
		hash := sf.Hash
		if v, existing := duplicates.Get(hash); existing {
			v = append(v, sf)
			duplicates.Set(hash, v)
		} else {
			duplicates.Set(hash, []SmashFile{sf})
		}
	}

}
