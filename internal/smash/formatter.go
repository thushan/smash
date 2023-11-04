package smash

import (
	"log"
)

const TreeLastChild = "└── "
const TreeNextChild = "├── "

func (app *App) printVerbose(message ...any) {
	if app.Flags.Verbose {
		log.Print(message...)
	}
}
