#####################
# smash make
# make smash
#####################
BINARY_NAME=smash
VERSION=$(shell git describe --tags --always --dirty)
GCOMMIT=$(shell git rev-parse --short HEAD)
TODAY=$(shell date --iso-8601)

.PHONY: all
all: lint test build

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build -ldflags " -X github.com/thushan/smash/internal/smash.Date=$(TODAY) \
                      -X github.com/thushan/smash/internal/smash.User=make \
                      -X github.com/thushan/smash/internal/smash.Version=$(VERSION) \
                      -X github.com/thushan/smash/internal/smash.Commit=$(GCOMMIT)" \
                      -o dist/$(BINARY_NAME) .

.PHONY: release
release:
	MSYS_NO_PATHCONV=1 docker run -ti -v "$(PWD):/app" -w "/app" goreleaser/goreleaser:latest release --snapshot --clean

.PHONY: clean
clean:
	go clean
	rm -rf dist/

.PHONY: clean-reports
clean-reports:
	rm -rf report-*.json
