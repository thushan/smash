# Developer Guide

This guide covers building, testing and contributing to Smash. For architecture details, see [CONTRIBUTING.md](../CONTRIBUTING.md).

## Table of Contents

- [Prerequisites](#prerequisites)
- [Building from Source](#building-from-source)
- [Testing](#testing)
- [Development Workflow](#development-workflow)
- [Docker Development](#docker-development)
- [Release Process](#release-process)
- [Code Standards](#code-standards)
- [Debugging](#debugging)

## Prerequisites

### Required Tools

- **Go 1.24+** - Check with `go version`
- **Make** - For build automation
- **Git** - For version control
- **Docker** (optional) - For release builds and testing

### Recommended Tools

- **golangci-lint** - For comprehensive linting
- **better-align** - For checking code alignment
- **goreleaser** - For release builds (via Docker)

## Building from Source

### Quick Build

```bash
# Clone the repository
git clone https://github.com/thushan/smash.git
cd smash

# Build the binary
make build
```

This creates the `smash` binary in the project root.

### Build Options

```bash
# Build with specific version info
go build -ldflags="-X github.com/thushan/smash/internal/smash.Version=dev" ./cmd/smash

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 make build
GOOS=darwin GOARCH=arm64 make build
GOOS=windows GOARCH=amd64 make build

# Build all platforms (requires Docker)
make release
```

### Understanding the Build

The build process:
1. Compiles the `cmd/smash` package
2. Embeds version information via ldflags
3. Produces a statically linked binary (CGO_ENABLED=0)

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./pkg/analysis/
go test -v ./pkg/slicer/

# Run a specific test
go test -v -run TestDuplicateAnalysis ./pkg/analysis/
```

### Test Structure

Tests follow Go conventions:
- Test files named `*_test.go`
- Test functions named `Test*`
- Benchmark functions named `Benchmark*`
- Test data in `testdata/` or `artefacts/` directories

Key test areas:
- **pkg/slicer**: File slicing logic with test artifacts
- **pkg/analysis**: Duplicate detection algorithms
- **pkg/indexer**: File system traversal
- **internal/algorithms**: Hash algorithm implementations

### Writing Tests

Example test structure:

```go
func TestFeatureName(t *testing.T) {
    // Arrange
    input := setupTestData()
    
    // Act
    result, err := FunctionUnderTest(input)
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

## Development Workflow

### Pre-commit Checks

```bash
# Run linting and tests before committing
make ready

# Individual checks
make lint        # Run golangci-lint
make align       # Check code alignment
make test        # Run tests
make build       # Verify build
```

### Code Organisation

```
smash/
├── cmd/smash/          # Main entry point
├── internal/           # Private packages
│   ├── algorithms/     # Hash algorithms
│   ├── cli/           # Command interface
│   ├── smash/         # Core logic
│   └── theme/         # UI theming
├── pkg/               # Public packages
│   ├── analysis/      # Duplicate analysis
│   ├── indexer/       # File indexing
│   ├── nerdstats/     # Statistics
│   ├── profiler/      # Performance profiling
│   └── slicer/        # File slicing
└── docs/              # Documentation
```

### Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature
   ```

2. Make changes following the code style

3. Add tests for new functionality

4. Run checks:
   ```bash
   make ready
   ```

5. Commit with descriptive messages:
   ```bash
   git commit -m "pkg/slicer: add adaptive slicing for small files"
   ```

## Docker Development

### Building Docker Images

```bash
# Test Docker build locally
make release

# Build just the Docker image
docker build -t smash:dev .

# Test the image
docker run --rm -v "$PWD:/data" smash:dev -r /data
```

### Docker Architecture

The Dockerfile uses a minimal Alpine Linux base:
- Final image ~8MB
- Runs as non-root user
- Includes only the binary and CA certificates

GoReleaser handles multi-architecture builds:
- Builds for linux/amd64 and linux/arm64
- Creates manifest lists for automatic platform selection
- Pushes to GitHub Container Registry (ghcr.io)

## Release Process

### Creating a Release

1. Update version in code if needed
2. Commit all changes
3. Create and push a tag:
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

GitHub Actions automatically:
- Runs tests
- Builds binaries for all platforms
- Creates Docker images
- Publishes to GitHub Releases
- Pushes images to ghcr.io

### Testing Releases

Before tagging a release:

```bash
# Create a snapshot build
make release

# Test the artifacts
./dist/smash_linux_amd64_v1/smash --version

# Test Docker images
docker run --rm ghcr.io/thushan/smash:1.2.4-abc1234 --version
```

### Version Information

Version details are embedded at build time:
- **Version**: From git tag or "dev"
- **Commit**: Git commit hash
- **Date**: Build timestamp
- **User**: Builder identifier

Access via:
```bash
smash --version
```

## Code Standards

### Go Conventions

Follow standard Go practices:
- Run `gofmt` on all code
- Use meaningful variable names
- Keep functions focused and small
- Document exported functions
- Handle errors explicitly

### Project Conventions

- **Logging**: Use structured logging via slog
- **Errors**: Wrap errors with context
- **Concurrency**: Use worker pools with bounded parallelism
- **Testing**: Aim for >80% coverage on critical paths

### Performance Considerations

- Minimize allocations in hot paths
- Use sync.Pool for frequently allocated objects
- Profile before optimizing
- Benchmark critical functions

## Debugging

### Enable Debug Logging

```bash
# Verbose output
smash -r -vvv ~/test

# CPU profiling
smash -r --cpuprofile=cpu.prof ~/test
go tool pprof cpu.prof

# Memory profiling
smash -r --memprofile=mem.prof ~/test
go tool pprof mem.prof
```

### Common Issues

**Build fails with module errors**
```bash
# Clean module cache
go clean -modcache
go mod download
```

**Tests fail with missing test data**
```bash
# Ensure test artifacts are present
git checkout -- pkg/slicer/artefacts/
```

**Docker build fails**
```bash
# Ensure binary is built first for manual Docker builds
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o smash .
docker build -t smash:test .
```

### Getting Help

- Check existing issues on GitHub
- Review test cases for usage examples
- Enable verbose logging for debugging
- Use profiling for performance issues
