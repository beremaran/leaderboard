GO=/usr/bin/env go
ENTRYPOINT=cmd/leaderboard

all: deps tidy apidoc build

build: deps tidy apidoc
	cd $(ENTRYPOINT); \
	CGO_ENABLED=0 $(GO) build -tags netgo -a -v
deps:
	$(GO) mod download
tidy:
	$(GO) mod tidy
apidoc:
	swag init --parseInternal -g $(ENTRYPOINT)/main.go
.PHONY: all build deps tidy swagger

test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress
