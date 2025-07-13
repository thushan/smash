# Smash User Guide

This guide covers common use cases for using Smash and provides examples for basic operations, advanced filtering, performance tuning, and more. 

It is designed to help users effectively manage duplicate files across various scenarios.

## Table of Contents
- [Basic Operations](#basic-operations)
- [Advanced Filtering](#advanced-filtering)
- [Performance Tuning](#performance-tuning)
- [Working with Large Files](#working-with-large-files)
- [Report Analysis](#report-analysis)
- [Common Use Cases](#common-use-cases)
- [Using Docker](#using-docker)
- [Troubleshooting](#troubleshooting)

## Basic Operations

### Finding Duplicates in Current Directory
```bash
# Basic scan (current directory only)
smash

# Recursive scan
smash -r

# Scan specific directories
smash -r ~/Documents ~/Downloads
```

### Output Control
```bash
# Silent mode - no output except report
smash -r --silent -o report.json ~/data

# Verbose mode - detailed processing info
smash -r --verbose ~/data

# Disable progress indicators
smash -r --no-progress ~/data
```

### Report Generation
```bash
# Auto-generated filename
smash -r ~/data

# Custom report filename
smash -r -o duplicates.json ~/data

# Disable report generation (console output only)
smash -r --no-output ~/data
```

## Advanced Filtering

### Size-Based Filtering
```bash
# Only files larger than 1MB
smash -r --min-size=1048576 ~/data

# Only files smaller than 100MB
smash -r --max-size=104857600 ~/data

# Files between 1MB and 100MB
smash -r --min-size=1048576 --max-size=104857600 ~/data
```

### Exclusion Patterns
```bash
# Exclude multiple directories
smash -r --exclude-dir=.git,.svn,node_modules ~/projects

# Exclude file patterns
smash -r --exclude-file="*.log,*.tmp,*.cache" ~/data

# Complex exclusions
smash -r \
  --exclude-dir=.git,target,build \
  --exclude-file="*.class,*.o,*.pyc" \
  ~/workspace
```

### System Files
```bash
# Include hidden files
smash -r --ignore-hidden=false ~/data

# Include system files
smash -r --ignore-system=false ~/data

# Include empty files in analysis
smash -r --ignore-empty=false ~/data
```

## Performance Tuning

### Worker Configuration
```bash
# Low-resource system
smash -r --max-workers=4 --max-threads=4 ~/data

# High-performance system
smash -r --max-workers=32 --max-threads=32 ~/data

# I/O bound operations (network drives)
smash -r --max-workers=8 ~/network-share
```

### Slicing Configuration
```bash
# More slices for very large files
smash -r --slices=8 --slice-size=16384 ~/videos

# Disable slicing for small files only
smash -r --slice-threshold=1048576 ~/documents

# Full file hashing (no slicing)
smash -r --disable-slicing ~/critical-data
```

### Algorithm Selection

See [Algorithm Guide](./algorithms.md) for details on available algorithms.

```bash
# Fastest (default)
smash -r --algorithm=xxhash ~/data

# Balanced speed/distribution
smash -r --algorithm=murmur3 ~/data

# Cryptographic verification
smash -r --algorithm=sha256 ~/data
```

## Working with Large Files

### Video Libraries
```bash
# Optimized for large video files
smash -r \
  --slices=8 \
  --slice-size=32768 \
  --min-size=10485760 \
  --exclude-file="*.srt,*.sub,*.idx" \
  ~/videos
```

### Photo Collections
```bash
# Optimized for photos with metadata
smash -r \
  --disable-meta=false \
  --exclude-file="*.xmp,Thumbs.db,.DS_Store" \
  --min-size=10240 \
  ~/Photos
```

### Source Code Repositories
```bash
# Optimized for code repos
smash -r \
  --exclude-dir=.git,node_modules,target,dist \
  --exclude-file="*.log,*.lock,*.sum" \
  --disable-autotext=false \
  ~/repos
```

## Report Analysis

### Basic Duplicate Extraction
```bash
# List all duplicate files
jq -r '.analysis.dupes[].files[].path' report.json

# Count duplicate sets
jq '.analysis.dupes | length' report.json

# Total space wasted
jq '.analysis.summary.spaceWasted' report.json
```

### Finding Specific Duplicates
```bash
# Files larger than 100MB
jq '.analysis.dupes[] | select(.files[0].size > 104857600)' report.json

# Duplicates in specific directory
jq '.analysis.dupes[].files[] | select(.path | startswith("/home/user/Downloads"))' report.json

# Group by hash
jq -r '.analysis.dupes[] | "\(.files[0].hash): \(.files | map(.path) | join(", "))"' report.json
```

### Empty File Analysis
```bash
# List all empty files
jq -r '.analysis.empty[].path' report.json

# Count empty files by directory
jq -r '.analysis.empty[].path | split("/")[:-1] | join("/")' report.json | sort | uniq -c
```

### Creating Action Scripts
```bash
# Generate rm commands for duplicates (keeping first)
jq -r '.analysis.dupes[].files[1:][].path | "rm -i \"\(.)\""' report.json > remove_dupes.sh

# Create hardlink script
jq -r '.analysis.dupes[] | 
  .files[0].path as $keep | 
  .files[1:][] | 
  "ln -f \"\($keep)\" \"\(.path)\""' report.json > create_links.sh
```

## Common Use Cases

### Backup Deduplication
```bash
# Compare backup directories
smash -r --show-duplicates \
  /backup/2024-01 \
  /backup/2024-02 \
  /backup/2024-03

# Find duplicates across local and external
smash -r ~/Documents /mnt/usb-backup/Documents
```

### Media Library Cleanup
```bash
# Music library
smash -r \
  --algorithm=murmur3 \
  --exclude-file="*.jpg,*.png,*.txt,*.cue" \
  --min-size=1048576 \
  ~/Music

# Photo library with RAW files
smash -r \
  --slices=6 \
  --exclude-file="*.xmp,*.pp3" \
  ~/Pictures
```

### Development Cleanup
```bash
# Find duplicate dependencies
smash -r \
  --exclude-dir=.git \
  --show-top=50 \
  ~/go/pkg/mod \
  ~/.cargo/registry

# Clean build artifacts
smash -r \
  --exclude-dir=.git,src \
  --min-size=1048576 \
  ~/projects
```

## Using Docker

Smash can be run in a Docker container, which is useful for consistent environments, CI/CD pipelines, or when you don't want to install the binary directly. The official Docker image is available on GitHub Container Registry.

### Docker Basics

> [!TIP]
> Use the `-t` flag to allocate a pseudo-TTY for better output formatting with Docker.
> 
> Leave it out if you don't need TTY support (e.g., in scripts/pipelines).
>
> We use the `--rm` flag to automatically remove the container after it exits, keeping
> your environment clean in these examples.

```bash
# Pull the latest image
docker pull ghcr.io/thushan/smash:latest

# Basic scan of current directory
docker run -t --rm -v "$PWD:/data" ghcr.io/thushan/smash:latest -r /data

# Scan with JSON output
docker run -t --rm -v "$PWD:/data" ghcr.io/thushan/smash:latest -r -o /data/report.json /data
```

### Volume Mounting

Docker containers are isolated from your host filesystem. You need to mount directories as volumes to scan them:

```bash
# Mount current directory as /data (read-only recommended for scanning)
docker run -t --rm -v "$PWD:/data:ro" ghcr.io/thushan/smash:latest -r /data

# Mount multiple directories
docker run -t --rm \
  -v "$HOME/Documents:/docs:ro" \
  -v "$HOME/Pictures:/pics:ro" \
  ghcr.io/thushan/smash:latest -r /docs /pics

# Windows PowerShell syntax
docker run -t --rm -v "${PWD}:/data:ro" ghcr.io/thushan/smash:latest -r /data
```

### Output File Handling

When generating reports, the output file must be written to a mounted volume. The container includes a built-in `/output` directory that is already writable, but you need to mount it to a host directory to access the reports:

```bash
# Create a local output directory and mount it
mkdir -p ./output
docker run -t --rm \
  -v "$PWD:/data:ro" \
  -v "$PWD/output:/output" \
  ghcr.io/thushan/smash:latest -r -o /output/report.json /data

# Now you can read the report on your host
cat ./output/report.json | jq '.analysis.summary'

# Use different output directories for different scans
mkdir -p ./reports/photos ./reports/documents
docker run -t --rm \
  -v "$HOME/Pictures:/data:ro" \
  -v "$PWD/reports/photos:/output" \
  ghcr.io/thushan/smash:latest -r -o /output/duplicates.json /data

# Alternative: write directly to mounted data directory
docker run -t --rm \
  -v "$PWD:/data" \
  ghcr.io/thushan/smash:latest -r -o /data/report.json /data
# Note: This requires the data mount to NOT be read-only
```

> [!NOTE]
> The `/output` directory inside the container is pre-configured with appropriate permissions for the non-root user. 
> When mounting to a host directory, ensure the host directory exists and is accessible.

### Advanced Docker Usage

#### Custom Algorithm and Filtering
```bash
docker run -t --rm -v "$PWD:/data:ro" ghcr.io/thushan/smash:latest \
  -r --algorithm=murmur3 \
  --exclude-dir=.git,node_modules \
  --min-size=1048576 \
  -o /data/large-files.json /data
```

#### Performance Tuning in Docker
```bash
# Limit resources for container
docker run -t --rm \
  --cpus="2.0" \
  --memory="1g" \
  -v "$PWD:/data:ro" \
  ghcr.io/thushan/smash:latest \
  -r --max-workers=8 --max-threads=8 /data
```

#### Using Specific Versions
```bash
# Use a specific version tag
docker pull ghcr.io/thushan/smash:v1.0.0
docker run -t --rm -v "$PWD:/data:ro" ghcr.io/thushan/smash:v1.0.0 -r /data
```

### Docker in CI/CD

For CI/CD pipelines, Docker provides a consistent environment:

```yaml
# GitHub Actions example
- name: Find duplicates
  run: |
    docker run --rm \
      -v ${{ github.workspace }}:/workspace:ro \
      -v ${{ github.workspace }}/reports:/output \
      ghcr.io/thushan/smash:latest \
      -r --silent -o /output/duplicates.json /workspace
```

```bash
# GitLab CI example
find-duplicates:
  image: ghcr.io/thushan/smash:latest
  script:
    - smash -r --silent -o duplicates.json .
  artifacts:
    paths:
      - duplicates.json
```

### Docker Compose

For complex setups, you can use Docker Compose:

```yaml
# docker-compose.yml
version: '3.8'
services:
  smash:
    image: ghcr.io/thushan/smash:latest
    volumes:
      - ./data:/data:ro
      - ./reports:/output
    command: -r --silent -o /output/report.json /data
```

Run with: `docker-compose run --rm smash`

### Docker Tips and Best Practices

1. **Use `-t` for interactive runs** - This ensures proper output formatting and colours, but can be omitted in scripts/pipelines.
2. **Use `--rm` to auto-cleanup** - Prevents accumulation of stopped containers
3. **Mount as read-only (`:ro`)** - Smash only reads files, so use `:ro` for safety
4. **Create output directories first** - Ensure output directories exist and are writable
5. **Use specific versions in production** - Pin to specific tags rather than `latest`

### Docker Troubleshooting

**No coloured output**
```bash
# Always use -t flag
docker run -t --rm -v "$PWD:/data" ghcr.io/thushan/smash:latest -r /data
```

**Permission denied errors**
```bash
# Ensure output directory is writable
chmod 755 ./reports
docker run -t --rm -v "$PWD:/data:ro" -v "$PWD/reports:/output" \
  ghcr.io/thushan/smash:latest -r -o /output/report.json /data
```

**Can't find files**
```bash
# Check your volume mounts - paths inside container are different
docker run -t --rm -v "$HOME/Documents:/docs" ghcr.io/thushan/smash:latest ls /docs
```

## Troubleshooting

### Performance Issues

**Slow scanning on network drives**
```bash
# Reduce workers and increase slice size
smash -r --max-workers=4 --slice-size=65536 /mnt/nas
```

**High memory usage**
```bash
# Reduce concurrent operations
smash -r --max-workers=8 --show-top=10 ~/data
```

**Too many open files error**
```bash
# Reduce workers or increase system limits
ulimit -n 4096
smash -r --max-workers=8 ~/data
```

### Accuracy Concerns

**False positives**
```bash
# Use more slices and full file hashing for critical data
smash -r --slices=16 --disable-slicing ~/important
```

**Missing duplicates**
```bash
# Check exclusion rules and size limits
smash -r \
  --ignore-hidden=false \
  --ignore-system=false \
  --min-size=0 \
  ~/data
```

### Report Issues

**Report too large**
```bash
# Limit output to significant duplicates
smash -r --min-size=10485760 --show-top=100 ~/data
```

**Can't parse report**
```bash
# Validate JSON
jq . report.json > /dev/null || echo "Invalid JSON"

# Pretty print for reading
jq . report.json > report-formatted.json
```

## Environment Variables

Smash respects standard environment variables:

```bash
# Disable colors in output
NO_COLOR=1 smash -r ~/data

# Force colors in non-TTY environment
FORCE_COLOR=1 smash -r ~/data | tee scan.log
```

## Integration Examples

### Cron Job for Regular Scanning
```bash
#!/bin/bash
# /etc/cron.daily/smash-scan
REPORT_DIR="/var/reports/smash"
mkdir -p "$REPORT_DIR"
/usr/local/bin/smash -r --silent \
  -o "$REPORT_DIR/scan-$(date +%Y%m%d).json" \
  /data
```

### Git Hook for Repository Scanning
```bash
#!/bin/bash
# .git/hooks/pre-commit
smash --no-output --exclude-dir=.git . || {
  echo "Duplicate files detected!"
  exit 1
}
```

### Monitoring Script
```bash
#!/bin/bash
# monitor-duplicates.sh
THRESHOLD=1073741824  # 1GB
REPORT=$(mktemp)
smash -r --silent -o "$REPORT" "$1"
WASTE=$(jq '.analysis.summary.spaceWasted // 0' "$REPORT")
if [ "$WASTE" -gt "$THRESHOLD" ]; then
  echo "Warning: $WASTE bytes wasted in duplicates"
  jq -r '.analysis.dupes[] | "\(.files | length) copies: \(.files[0].path)"' "$REPORT"
fi
rm -f "$REPORT"
```
