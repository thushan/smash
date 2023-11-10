package app

type Flags struct {
	Base        []string `yaml:"base"`
	ExcludeDir  []string `yaml:"exclude-dir"`
	ExcludeFile []string `yaml:"exclude-file"`
	Algorithm   int      `yaml:"algorithm"`
	MaxThreads  int      `yaml:"max-threads"`
	MaxWorkers  int      `yaml:"max-workers"`
	Silent      bool     `yaml:"silent"`
	Verbose     bool     `yaml:"verbose"`
}
