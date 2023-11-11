package smash

import (
	"fmt"
	"log"

	"github.com/alphadose/haxmap"
	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora/v3"
)

const (
	TreeLastChild = "└── "
	TreeNextChild = "├── "
)

func (app *App) printVerbose(message ...any) {
	if app.Flags.Verbose {
		log.Print(message...)
	}
}

func (app *App) printSmashHits(cache *haxmap.Map[string, []SmashFile]) uint64 {
	totalDuplicateSize := uint64(0)
	cache.ForEach(func(hash string, files []SmashFile) bool {
		mainFile := files[0]
		lastIndex := len(files)
		if lastIndex > 1 {
			log.Println(aurora.Magenta(mainFile.Filename), " ", aurora.Cyan(humanize.Bytes(mainFile.FileSize)), " ", aurora.Blue(mainFile.Hash))
			for index, file := range files[1:] {
				var subTree string
				if (index + 2) == lastIndex {
					subTree = TreeLastChild
				} else {
					subTree = TreeNextChild
				}
				log.Println(aurora.BrightYellow(subTree), file.Filename)
			}
			totalDuplicateSize += mainFile.FileSize * uint64(lastIndex-1)
		} else {
			// prune unique files
			cache.Del(hash)
		}
		return true
	})
	return totalDuplicateSize
}

func (app *App) printSmashRunSummary(rs RunSummary) {
	log.Println("Total Time:   ", aurora.Green(fmt.Sprintf("%dms", rs.ElapsedTime)))
	log.Println("Total Files:  ", aurora.Blue(rs.TotalFiles))
	log.Println("Total Unique: ", aurora.Blue(rs.UniqueFiles))
	log.Println("Total Duplicates: ", aurora.Blue(rs.DuplicateFiles), "(", aurora.Cyan(rs.DuplicateFileSizeF), " can be reclaimed).")
}
