package smash

import (
	"github.com/thushan/smash/internal/app"
)

type App struct {
	Args      []string
	Locations []string
	Flags     *app.Flags
}

func (a *App) Run() error {
	return nil
}
