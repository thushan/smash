package smash

type Flags struct {
	OutputFile        string   `yaml:"output"`
	Base              []string `yaml:"base"`
	ExcludeDir        []string `yaml:"exclude-dir"`
	ExcludeFile       []string `yaml:"exclude-file"`
	Algorithm         int      `yaml:"algorithm"`
	MaxThreads        int      `yaml:"max-threads"`
	MaxWorkers        int      `yaml:"max-workers"`
	UpdateSeconds     int      `yaml:"update-seconds"`
	ShowTop           int      `yaml:"show-top"`
	DisableSlicing    bool     `yaml:"disable-slicing"`
	IgnoreEmptyFiles  bool     `yaml:"ignore-emptyfiles"`
	IgnoreHiddenItems bool     `yaml:"ignore-hiddenitems"`
	ShowVersion       bool     `yaml:"show-version"`
	ShowDuplicates    bool     `yaml:"show-duplicates"`
	HideTop           bool     `yaml:"hide-top"`
	Silent            bool     `yaml:"silent"`
	NoProgress        bool     `yaml:"no-progress"`
	Verbose           bool     `yaml:"verbose"`
}
