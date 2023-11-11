package smash

import (
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

func (app *App) printSmashHits(cache *haxmap.Map[string, []SmashFile]) {
	cache.ForEach(func(hash string, file []SmashFile) bool {
		mainFile := file[0]
		lastIndex := len(file[0:]) - 1
		if lastIndex > 0 {
			log.Println(aurora.Magenta(mainFile.Filename), " ", aurora.Cyan(humanize.Bytes(mainFile.FileSize)), " ", aurora.Blue(mainFile.Hash))
			for _, file := range file[1:] {
				log.Println(aurora.BrightYellow("> "), file.Filename)
			}
		}
		return true
	})
}
