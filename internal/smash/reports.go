package smash

import (
	"encoding/hex"
	"path/filepath"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/indexer"
	"github.com/thushan/smash/pkg/slicer"
)

func calculateRunSummary(cache *haxmap.Map[string, []SmashFile], fails *haxmap.Map[string, error], totalFiles int64, appStartTime int64) RunSummary {

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
	if v, ok := cache.GetAndDel(zeroByteCookie); ok {
		return len(v)
	} else {
		return 0
	}
}

func (app *App) summariseSmashedFile(cache *haxmap.Map[string, []SmashFile], stats slicer.SlicerStats, filename string, ms int64) {
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

func resolveFilename(file indexer.FileFS) string {
	if file.Path == "." {
		return filepath.Base(file.FullName)
	} else {
		return file.Path
	}
}
