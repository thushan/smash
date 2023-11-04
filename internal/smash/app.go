package smash

import (
	"github.com/thushan/smash/internal/app"
)

type App struct {
	Args      []string
	Locations []string
	Flags     *app.Flags
}

func (app *App) Run() error {

	if !app.Flags.Silent {
		app.printConfiguration()
	}
	return nil
}
