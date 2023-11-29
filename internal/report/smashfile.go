package report

import (
	"encoding/hex"
	"github.com/thushan/smash/internal/theme"

	"github.com/alphadose/haxmap"
	"github.com/dustin/go-humanize"
	"github.com/thushan/smash/pkg/slicer"
)

type SmashFile struct {
	Filename    string
	Hash        string
	FileSizeF   string
	FileSize    uint64
	ElapsedTime int64
	FullHash    bool
	EmptyFile   bool
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
		if ov, loaded := duplicates.GetOrSet(hash, []SmashFile{sf}); loaded {
			v := append(ov, sf)
			if swapped := duplicates.CompareAndSwap(hash, ov, v); !swapped {
				theme.Error.Println("Swap failed for ", hash, ". old: ", len(ov), " | new: ", len(v))
			}
		}
	}

}
