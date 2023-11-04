package app

type Flags struct {
	Silent  bool `yaml:"silent"`
	Verbose bool `yaml:"verbose"`
}

var HashAlgorithms = map[int][]string{
	0: {"xxhash"},
}
