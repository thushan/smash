# smash

[![GitHub license](https://img.shields.io/github/license/thushan/smash)](https://github.com/thushan/smash/blob/master/LICENSE)
[![CI](https://github.com/thushan/smash/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thushan/smash/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thushan/smash)](https://goreportcard.com/report/github.com/thushan/smash)
[![GitHub release](https://img.shields.io/github/release/thushan/smash)](https://github.com/thushan/smash/releases/latest)

CLI tool to `smash` through to find duplicate files efficiently by slicing a file (or blob) into multiple segments
and computing a hash using a fast non-cryptographic algorithm such as [xxhash](https://xxhash.com/) or [murmur3](https://en.wikipedia.org/wiki/MurmurHash).

Amongst the highlights of `smash`:

* Super fast analysis of large files thanks to [slicing](./docs/slicing.md).
* Suited for finding duplicates on bandwidth constrained networks, devices or very large files but plenty capable on smaller ones!
* Supports a variety of non-cryptographic algorithms (see [algorithms supported](./docs/algorithms.md)).
* Read-only view of the underlying filesystem when analysing
* Reports on duplicate files & empty (0 byte) files
* Outputs a report in json, you can use tools like [jq](https://github.com/jqlang/jq) to operate on (see [examples](#examples) below or [the vhs tapes](./docs/demos.md))
* Used to dedupe multi-TB of astrophysics datasets, images and video content & run regularly to report duplicates

`smash` does not support pruning of duplicates or empty files natively and it's encouraged you vet the output report before pruning via automated tools.

<p align="center">
 <img src="https://vhs.charm.sh/vhs-6UTX5Yc6CIQ6Y3lzulLKYF.gif" alt="Made with VHS"><br/>
    <sub>
        <sup>Find duplicates in the <a href="https://github.com/torvalds/linux">linux/drivers</a> source tree with <code>smash</code> (see our <a href="docs/demos.md">🍿 other demos</a>). Made with <a href="https://vhs.charm.sh" target="_blank">vhs</a>!</sup>
    </sub>
</p>

The name comes from a prototype tool called SmartHash (written many years ago in C/ASM that's now lost in source & 
too hard to modernise). It operated on a similar concept of slicing and hashing (with CRC32 then later MD5).

# Installation

[![Operating Systems](https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux%20%7C%20freebsd-informational?style=for-the-badge)](https://github.com/thushan/smash/releases/latest)

You can download the latest binaries from [Github Releases](https://github.com/thushan/smash/releases) or via our [simple installer script](https://raw.githubusercontent.com/thushan/smash/main/install.sh) - which currently supports Linux, macos, FreeBSD & Windows:

```bash
bash <(curl -s https://raw.githubusercontent.com/thushan/smash/main/install.sh)
```

It will download the latest version & extract it to its own folder for you.

Alternatively, you can install it via go:

```bash
go install github.com/thushan/smash@latest
```

`smash` has been developed on Linux (Pop!_OS & Fedora), tested on macOS, FreeBSD & Windows.

# Usage

```bash
# Basic usage - scan current directory
smash

# Recursive scan
smash -r

# Scan multiple directories
smash -r ~/Documents ~/Downloads

# Silent mode with report
smash -r --silent -o report.json ~/data
```

For detailed usage, see the [User Guide](./docs/user-guide.md).

## Command Line Options

Key flags:
- `-r, --recurse` - Scan subdirectories (required for recursive scanning)
- `-o, --output-file` - Save results to JSON file
- `--silent` - Suppress all output except errors
- `--algorithm` - Choose hash algorithm (default: xxhash)
- `--exclude-dir` - Skip directories (comma-separated)
- `--exclude-file` - Skip files (comma-separated patterns)

Run `smash --help` for complete options.

## Quick Examples

### Find Duplicates
```bash
# In photos directory
smash -r ~/photos -o duplicates.json

# Across multiple drives
smash -r ~/Documents /mnt/backup/Documents

# Large video files only
smash -r --min-size=104857600 ~/Videos
```

### Filter and Exclude
```bash
# Skip git and node_modules
smash -r --exclude-dir=.git,node_modules ~/projects

# Include empty files
smash -r --ignore-empty=false ~/data
```

### Performance Tuning
```bash
# For network drives
smash -r --max-workers=4 /mnt/nas

# For many small files
smash -r --disable-slicing ~/documents
```

### Working with Reports
```bash
# Generate report
smash -r ~/data -o report.json

# List all duplicates
jq -r '.analysis.dupes[].files[].path' report.json

# Show space wasted
jq '.analysis.summary.spaceWasted' report.json
```

See the [User Guide](./docs/user-guide.md) for detailed examples and advanced usage.

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

# License

Copyright (c) Thushan Fernando and licensed under Apache License 2.0
