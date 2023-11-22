package smash

type Flags struct {
	OutputFile       string   `yaml:"output"`
	Base             []string `yaml:"base"`
	ExcludeDir       []string `yaml:"exclude-dir"`
	ExcludeFile      []string `yaml:"exclude-file"`
	Algorithm        int      `yaml:"algorithm"`
	MaxThreads       int      `yaml:"max-threads"`
	MaxWorkers       int      `yaml:"max-workers"`
	DisableSlicing   bool     `yaml:"disable-slicing"`
	IgnoreEmptyFiles bool     `yaml:"ignore-emptyfiles"`
	ShowVersion      bool     `yaml:"show-version"`
	Silent           bool     `yaml:"silent"`
	Verbose          bool     `yaml:"verbose"`
}
