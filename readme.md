# smash

[![CI](https://github.com/thushan/smash/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thushan/smash/actions/workflows/ci.yml)
[![Lint](https://github.com/thushan/smash/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/thushan/smash/actions/workflows/lint.yml)
[![GitHub license](https://img.shields.io/github/license/thushan/smash)](https://github.com/thushan/smash/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/thushan/smash)](https://github.com/thushan/smash/releases/latest)

Tool to `smash` through to find duplicate files efficiently by slicing a file (or blob) into multiple segments and computing a hash using [xxhash](https://xxhash.com/) or another algorithm. The name comes from a prototype tool called SmartHash (written many years ago in C+ASM that's now lost in source & too hard to modernise) that operated on a similar concept.

It is ideally suited to finding duplicates on bandwidth constrained devices (or networks) or very large files but is ridiculously fast on SSDs/NVMe's where you want to quickly determine duplicate files.

# Usage

```
Usage:
  smash [flags] [locations-to-smash]

Flags:
      --algorithm algorithm    Algorithm to use, can be 'xxhash' (default), 'fnv128', 'fnv128a' (default xxhash)
      --disable-slicing        Disable slicing (hashes full file).
      --exclude-dir strings    Directories to exclude separated by comma. Eg. --exclude-dir=.git,.idea
      --exclude-file strings   Files to exclude separated by comma. Eg. --exclude-file=.gitignore,*.csv
  -h, --help                   help for smash
  -p, --max-threads int        Maximum threads to utilise. (default 16)
  -w, --max-workers int        Maximum workers to utilise when smashing. (default 8)
  -q, --silent                 Run in silent mode.
      --verbose                Run in verbose mode.
  -v, --version                version for smash
```

# Acknowledgements

This project was possible thanks to the following projects or folks.

* [@wader/fq](https://github.com/wader/fq) - countless nights of inspecting binary blobs!
* [@cespare/xxhash](https://github.com/cespare/xxhash) - xxhash implementation
* [@spf13/cobra](https://github.com/spf13/cobra) - CLI Magic with Cobra
* [@golangci/golangci-lint](https://github.com/golangci/golangci-lint) - Go Linter
* [@dkorunic/betteralign](https://github.com/dkorunic/betteralign) - Go alignment checker

Testers - MarkB, JarredT, BenW, DencilW, JayT, ASV

# Licence

Copyright (c) Thushan Fernando and licensed under Apache Licence 2.0
