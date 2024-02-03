package smash

import "errors"

type Flags struct {
	OutputFile      string   `yaml:"output"`
	Base            []string `yaml:"base"`
	ExcludeDir      []string `yaml:"exclude-dir"`
	ExcludeFile     []string `yaml:"exclude-file"`
	MinSize         int64    `yaml:"min-size"`
	MaxSize         int64    `yaml:"max-size"`
	Algorithm       int      `yaml:"algorithm"`
	MaxThreads      int      `yaml:"max-threads"`
	MaxWorkers      int      `yaml:"max-workers"`
	ProgressUpdate  int      `yaml:"progress-update"`
	ShowTop         int      `yaml:"show-top"`
	DisableSlicing  bool     `yaml:"disable-slicing"`
	DisableMeta     bool     `yaml:"disable-meta"`
	DisableAutoText bool     `yaml:"disable-autotext"`
	IgnoreEmpty     bool     `yaml:"ignore-empty"`
	IgnoreHidden    bool     `yaml:"ignore-hidden"`
	IgnoreSystem    bool     `yaml:"ignore-system"`
	ShowVersion     bool     `yaml:"version"`
	ShowNerdStats   bool     `yaml:"nerd-stats"`
	Recurse         bool     `yaml:"recurse"`
	ShowDuplicates  bool     `yaml:"show-duplicates"`
	Silent          bool     `yaml:"silent"`
	HideTopList     bool     `yaml:"no-top-list"`
	HideProgress    bool     `yaml:"no-progress"`
	HideOutput      bool     `yaml:"no-output"`
	Profile         bool     `yaml:"profile"`
	Verbose         bool     `yaml:"verbose"`
}

func (app *App) validateArgs() error {
	f := app.Flags
	if f.Silent && f.Verbose {
		return errors.New("cannot be verbose and silent")
	}
	if f.MaxThreads < 0 {
		return errors.New("maxthreads cannot be below zero")
	}
	if f.MaxWorkers < 0 {
		return errors.New("maxworkers cannot be below zero")
	}
	if f.MinSize < 0 {
		return errors.New("minSize cannot be below zero")
	}
	if f.MaxSize < 0 {
		return errors.New("maxSize cannot be below zero")
	}
	if f.MaxSize != 0 && f.MinSize > f.MaxSize {
		return errors.New("minSize cannot be greater than maxSize")
	}
	if f.ShowTop <= 1 {
		return errors.New("showtop should be greater than 1")
	}
	if f.ShowTop != 10 && f.HideTopList {
		return errors.New("cannot mix showtop x and hidetop")
	}
	if f.ProgressUpdate < 1 {
		return errors.New("updateseconds cannot be less than 1")
	}

	return nil
}
