# smash
[![CI](https://github.com/thushan/smash/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thushan/smash/actions/workflows/ci.yml)

aka SmartHash, a tool to `smash` through to find duplicate files efficiently by slicing a file (or blob) into multiple segments and computing a hash using [xxhash](https://xxhash.com/).

It is ideally suited to finding duplicates on bandwidth constrained devices (or networks) but is ridiculously fast on SSDs/NVMe's where you want to quickly determine duplicate files.

# Usage

```
$  smash [flags] <locations-to-smash>

Flags:
  --max-threads         Maximum number of threads to utilise
  -h, --help            Shows help for smash
  -v, --version         Print version information
```

# Licence

Copyright (c) Thushan Fernando and licensed under Apache Licence 2.0
