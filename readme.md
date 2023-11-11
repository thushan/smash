# smash

[![CI](https://github.com/thushan/smash/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thushan/smash/actions/workflows/ci.yml)
[![Lint](https://github.com/thushan/smash/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/thushan/smash/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thushan/smash)](https://goreportcard.com/report/github.com/thushan/smash)
[![Maintainability](https://api.codeclimate.com/v1/badges/944834a9d91128fa690d/maintainability)](https://codeclimate.com/github/thushan/smash/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/944834a9d91128fa690d/test_coverage)](https://codeclimate.com/github/thushan/smash/test_coverage)
[![GitHub license](https://img.shields.io/github/license/thushan/smash)](https://github.com/thushan/smash/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/thushan/smash?sort=semver)](https://github.com/thushan/smash/tags)
[![GitHub release](https://img.shields.io/github/release/thushan/smash)](https://github.com/thushan/smash/releases/latest)

Tool to `smash` through to find duplicate files efficiently by slicing a file (or blob) into multiple segments
and computing a hash using [xxhash](https://xxhash.com/) or another algorithm. Algorithms used are non-cryptographic and are utilised
for its speed & efficiency over other attributes. You can read [further about xxhash](https://xxhash.com/).

The name comes from a prototype tool called SmartHash (written many years ago in C+ASM that's now lost in source & too hard to modernise) that operated on a similar concept.

It is ideally suited to finding duplicates on bandwidth constrained devices (or networks) or very large files but is ridiculously fast on SSDs/NVMe's where you want to quickly determine duplicate files.

# Installation

[![Operating Systems](https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux%20%7C%20freebsd-informational?style=for-the-badge)](https://github.com/thushan/smash/releases/latest)


You can download the latest binaries from the [Releases](https://github.com/thushan/smash/releases) page and extract & use on your chosen operating system.

Alternatively, you can clone the repo and compile it from source - go will download dependencies.

```bash
$ go run .
```

`smash` has been developed on Linux (Pop!_OS & Fedora), tested on macOS, FreeBSD & Windows.

# Usage

```
Usage:
  smash [flags] [locations-to-smash]

Flags:
      --algorithm algorithm    Algorithm to use, can be 'xxhash', 'fnv128', 'fnv128a' (default xxhash)
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

## Examples

Examples are given in Unix format, but apply to Windows as well.

### Simplest

To check for duplicates in a single path (Eg. `~/media/photos`)

```bash
$ ./smash ~/media/photos
```

### Multiple Directories

To check across multiple directories - which can be different drives, or mounts (Eg. `~/media/photos` and `/mnt/my-usb-drive/photos`):

```bash
$ ./smash ~/media/photos /mnt/my-usb-drive/photos
```

Smash will find and report all duplicates within any number of directories passed in.

### Exclude Files or Directories

You can exclude certain directories or files with the `--exclude-dir` and `--exclude-file` switches including wildcard characters:

```bash
$ ./smash --exclude-dir=.git,.svn --exclude-file=.gitignore,*.csv ~/media/photos
```

For example, to ignore all hidden files on unix (those that start with `.` such as `.config` or `.gnome` folders):

```bash
$ ./smash --exclude-dir=.config,.gnome ~/media/photos
```

### Disabling Slicing & Getting Full Hash

By default, smash uses slicing to efficiently slice a file into mulitple segments and hash parts of the file. 

If you prefer not to use slicing for a run, you can disable slicing with:

```bash
$ ./smash --disable-slicing ~/media/photos
```

### Changing Hashing Algorithms

By default, smash uses `xxhash`, an extremely fast non-cryptographic hash algorithm 
(which you can [read about further](https://xxhash.com/)). 

To use another supported algorithm, use the `--algorithm` switch:

```bash
$ ./smash --algorithm:fnv128a ~/media/photos
```

# Acknowledgements

This project was possible thanks to the following projects or folks.

* [@wader/fq](https://github.com/wader/fq) - countless nights of inspecting binary blobs!
* [@cespare/xxhash](https://github.com/cespare/xxhash) - xxhash implementation
* [@alphadose/haxmap](https://github.com/alphadose/haxmap) - Amazingly efficient map implementation
* [@spf13/cobra](https://github.com/spf13/cobra) - CLI Magic with Cobra
* [@golangci/golangci-lint](https://github.com/golangci/golangci-lint) - Go Linter
* [@dkorunic/betteralign](https://github.com/dkorunic/betteralign) - Go alignment checker

Testers - MarkB, JarredT, BenW, DencilW, JayT, ASV

# Licence

Copyright (c) Thushan Fernando and licensed under Apache Licence 2.0
