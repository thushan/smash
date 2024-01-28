# smash

[![GitHub license](https://img.shields.io/github/license/thushan/smash)](https://github.com/thushan/smash/blob/master/LICENSE)
[![CI](https://github.com/thushan/smash/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thushan/smash/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thushan/smash)](https://goreportcard.com/report/github.com/thushan/smash)
[![GitHub release](https://img.shields.io/github/release/thushan/smash)](https://github.com/thushan/smash/releases/latest)

CLI tool to `smash` through to find duplicate files efficiently by slicing a file (or blob) into multiple segments
and computing a hash using a fast non-cryptographic algorithm such as [xxhash](https://xxhash.com/) or [murmur3](https://en.wikipedia.org/wiki/MurmurHash).

Amongst the highlights of `smash`:

* Super fast analysis of large files thanks to slicing.
* Suited for finding duplicates on bandwidth constrained networks, devices or very large files but plenty capable on smaller ones!
* Supports a variety of non-cryptographic algorithms (see [algorithms supported](./docs/algorithms.md)).
* Read-only view of the underlying filesystem when analysing
* Reports on duplicate files & empty (0 byte) files
* Outputs a report in json, you can use tools like [jq](https://github.com/jqlang/jq) to operate on (see [examples](#examples) below or [the vhs tapes](./docs/demos.md))
* Used to dedupe multi-TB of astrophysics datasets, images and video content & run regularly to report duplicates

`smash` does not support pruning of duplicates or empty files natively and it's encouraged you vet the output report before pruning via automated tools.

<p align="center">
 <img src="https://vhs.charm.sh/vhs-7BJdHGJLipNTwKQjJjDhXV.gif" alt="Made with VHS"><br/>
    <sub>
        <sup>Find duplicates in the <a href="https://github.com/torvalds/linux">linux/drivers</a> source tree with <code>smash</code> (see our <a href="docs/demos.md">🍿 other demos</a>). Made with <a href="https://vhs.charm.sh" target="_blank">vhs</a>!</sup>
    </sub>
</p>

The name comes from a prototype tool called SmartHash (written many years ago in C/ASM that's now lost in source & 
too hard to modernise). It operated on a similar concept of slicing and hashing (with CRC32 then later MD5).

# Installation

[![Operating Systems](https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux%20%7C%20freebsd-informational?style=for-the-badge)](https://github.com/thushan/smash/releases/latest)

You can download the latest binaries from the [Releases](https://github.com/thushan/smash/releases) page and extract & use on your chosen operating system. We
currently only support 64-bit binaries.

Alternatively, you can install it via go:

```bash
$ go install github.com/thushan/smash@latest
```

`smash` has been developed on Linux (Pop!_OS & Fedora), tested on macOS, FreeBSD & Windows.

# Usage

```
Usage:
  smash [flags] [locations-to-smash]

Flags:
      --algorithm algorithm    Algorithm to use to hash files. Supported: xxhash, murmur3, md5, sha512, sha256 (full list, see readme) (default xxhash)
      --base strings           Base directories to use for comparison Eg. --base=/c/dos,/c/dos/run/,/run/dos/run
      --disable-autotext       Disable detecting text-files to opt for a full hash for those
      --disable-meta           Disable storing of meta-data to improve hashing mismatches
      --disable-slicing        Disable slicing & hash the full file instead
      --exclude-dir strings    Directories to exclude separated by comma Eg. --exclude-dir=.git,.idea
      --exclude-file strings   Files to exclude separated by comma Eg. --exclude-file=.gitignore,*.csv
  -h, --help                   help for smash
      --ignore-empty           Ignore empty/zero byte files (default true)
      --ignore-hidden          Ignore hidden files & folders Eg. files/folders starting with '.' (default true)
      --ignore-system          Ignore system files & folders Eg. '$MFT', '.Trash' (default true)
  -p, --max-threads int        Maximum threads to utilise (default 16)
  -w, --max-workers int        Maximum workers to utilise when smashing (default 16)
      --nerd-stats             Show nerd stats
      --no-output              Disable report output
      --no-progress            Disable progress updates
      --no-top-list            Hides top x duplicates list
  -o, --output-file string     Export analysis as JSON (generated automatically otherwise)
      --profile                Enable Go Profiler - see localhost:1984/debug/pprof
      --progress-update int    Update progress every x seconds (default 5)
      --show-duplicates        Show full list of duplicates
      --show-top int           Show the top x duplicates (default 10)
  -q, --silent                 Run in silent mode
      --verbose                Run in verbose mode
  -v, --version                Show version information
```

See the [full list of algorithms](./docs/algorithms.md) supported.

## Examples

Examples are given in Unix format, but apply to Windows as well.

### Basic

To check for duplicates in a single path (Eg. `~/media/photos`) & output report to `report.json`

```bash
$ ./smash ~/media/photos -o report.json
```

You can then look at `report.json` with [jq](https://github.com/jqlang/jq) to check duplicates:

```bash 
$ jq '.analysis.dupes[]|[.location,.path,.filename]|join("/")' report.json | xargs wc -l
```

### Show Empty Files

By default, `smash` ignores empty files but can report on them with the `--ignore-empty=false` argument:

```bash
$ ./smash ~/media/photos --ignore-empty=false -o report.json
```

You can then look at `report.json` with [jq](https://github.com/jqlang/jq) to check empty files:

```bash 
$ jq '.analysis.empty[]|[.location,.path,.filename]|join("/")' report.json | xargs wc -l
```

### Show Top 50 Duplicates

By default, `smash` shows the top 10 duplicate files in the CLI and leaves the rest for the report, you can change that with the `--show-top=50` argument to show the top 50 instead.

```bash
$ ./smash ~/media/photos --show-top=50
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

By default, `smash` uses slicing to efficiently slice a file into mulitple segments and hash parts of the file. 

If you prefer not to use slicing for a run, you can disable slicing with:

```bash
$ ./smash --disable-slicing ~/media/photos
```

### Changing Hashing Algorithms

By default, smash uses `xxhash`, an extremely fast non-cryptographic hash algorithm. However, you can choose a variety
of algorithms [as documented](./docs/algorithms.md).

To use another supported algorithm, use the `--algorithm` switch:

```bash
$ ./smash --algorithm:murmur3 ~/media/photos
```

# Acknowledgements

This project was possible thanks to the following projects or folks.

* [@jqlang/jq](https://github.com/jqlang/jq) - without `jq` we'd be a bit lost!
* [@wader/fq](https://github.com/wader/fq) - countless nights of inspecting binary blobs!
* [@cespare/xxhash](https://github.com/cespare/xxhash) - xxhash implementation
* [@spaolacci/murmur3](https://github.com/spaolacci/murmur3) - murmur3 implementation
* [@puzpuzpuz/xsync](https://github.com/puzpuzpuz/xsync) - Amazingly efficient map implementation
* [@pterm/pterm](https://github.com/pterm/pterm) - Amazing TUI framework used
* [@spf13/cobra](https://github.com/spf13/cobra) - CLI Magic with Cobra
* [@golangci/golangci-lint](https://github.com/golangci/golangci-lint) - Go Linter
* [@dkorunic/betteralign](https://github.com/dkorunic/betteralign) - Go alignment checker

Testers - MarkB, JarredT, BenW, DencilW, JayT, ASV, TimW, RyanW, WilliamH, SpencerB, EmadA, ChrisE, AngelaB, LisaA, YousefI, JeffG, MattP

# Licence

Copyright (c) Thushan Fernando and licensed under Apache Licence 2.0
