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
}

func CalculateRunSummary(cache *haxmap.Map[string, []SmashFile], fails *haxmap.Map[string, error], totalFiles int64, appStartTime int64) RunSummary {

	emptyFileCount := getEmptyFileCount(cache)

	return RunSummary{
		TotalFiles:      totalFiles,
		TotalFileErrors: int64(fails.Len()),
		UniqueFiles:     int64(cache.Len()),
		EmptyFiles:      int64(emptyFileCount),
		DuplicateFiles:  totalFiles - (int64(cache.Len()) + int64(emptyFileCount)),
		ElapsedTime:     time.Now().UnixMilli() - appStartTime,
	}
}
func getEmptyFileCount(cache *haxmap.Map[string, []SmashFile]) int {

	zeroByteCookie := hex.EncodeToString(slicer.DefaultEmptyFileCookie)
	if v, ok := cache.Get(zeroByteCookie); ok {
		return len(v)
	} else {
		return 0
	}
}

func SummariseSmashedFile(cache *haxmap.Map[string, []SmashFile], stats slicer.SlicerStats, filename string, ms int64) {
	sf := SmashFile{
		Filename:    filename,
		Hash:        hex.EncodeToString(stats.Hash),
		FileSize:    stats.FileSize,
		FullHash:    stats.HashedFullFile,
		FileSizeF:   humanize.Bytes(stats.FileSize),
		ElapsedTime: ms,
	}
	hash := sf.Hash
	if v, existing := cache.Get(hash); existing {
		v = append(v, sf)
		cache.Set(hash, v)
	} else {
		cache.Set(hash, []SmashFile{sf})
	}
}
