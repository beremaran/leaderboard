GO=/usr/bin/env go
ENTRYPOINT=cmd/leaderboard

all: deps tidy swagger build

build:
	cd $(ENTRYPOINT); \
	CGO_ENABLED=0 $(GO) build -tags netgo -a -v
deps:
	$(GO) mod download
tidy:
	$(GO) mod tidy
apidoc:
	swag init --parseInternal -g $(ENTRYPOINT)/main.go
.PHONY: all build deps tidy swagger