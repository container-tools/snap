
GOLDFLAGS += -X main.GitCommit=$(GIT_COMMIT)
GOFLAGS = -ldflags "$(GOLDFLAGS)" -trimpath

default: build

build: build-resources build-snap test

build-snap:
	go build $(GOFLAGS) -o snap ./cmd/manager/*.go

build-resources:
	go run ./cmd/util/vfs-gen/ deploy

test:
	go test ./...

.PHONY: build build-resources build-snap test
