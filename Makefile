# === CONFIG =======================================================
PROJECT_NAME="league-of-legends-tactics-cli"
COVER_PROFILE="build/cover.out"


# === BUILD =======================================================
build-lol-tactics:
	@echo "---> Building $(PROJECT_NAME)"
	go build -o build/league_of_legends_tactics_cli cmd/$(PROJECT_NAME)/main.go
.PHONY: build-lol-tactics

# === INSTALL =======================================================
install-lol-tactics:
	@echo "--> Installing $(PROJECT_NAME)"
	@go install  ./cmd/league-of-legends-tactics-cli
.PHONY: install-lol-tactics

# === RUN =======================================================
# Shows a list of commands or help for one command
run-lol-tactics-help:
	go run ./cmd/league-of-legends-tactics-cli/main.go -h
.PHONY: run-lol-tactics-help

# Generate all fight tactics
run-lol-tactics-fight:
	go run ./cmd/league-of-legends-tactics-cli/main.go -f $(f1) -f $(f2)
.PHONY: run-lol-tactics-fight

# Generate fight tactics between specific league of legends champions
run-lol-tactics-all:
	go run ./cmd/league-of-legends-tactics-cli/main.go -t
.PHONY: run-lol-tactics-all

# Download and update a specific league of legends champion
run-lol-tactics-download:
	go run ./cmd/league-of-legends-tactics-cli/main.go -d $(d)
.PHONY: run-lol-tactics-download

# Download and update all league of legends champions
run-lol-tactics-download-all:
	go run ./cmd/league-of-legends-tactics-cli/main.go -a
.PHONY: run-lol-tactics-download-all

# === TEST =======================================================
test:
	@echo "---> Running all tests"
	go test -race -cover -coverprofile=$(COVER_PROFILE) ./...
.PHONY: test

# === CLEAN =======================================================
# Clean lol fights
clean-lol-fights:
	@echo "---> Cleaning lol fights"
	rm fights/*
.PHONY: clean-lol-fights

# Clean lol champions
clean-lol-champions:
	@echo "---> Cleaning lol champions data"
	rm champions/lol/*
.PHONY: clean-lol-champions

# === TOOLS =======================================================
# Get a decorated HTML presentation of cover file: showing the covered (green), uncovered (red), and un-instrumented (grey) source.
tool-cover:
	go tool cover -html=$(COVER_PROFILE)
.PHONY: tool-cover

# Fix go.mod and go.sum
tool-mod-tidy:
	@echo "---> Checking module requirements"
	go mod tidy
.PHONY: tool-mod-tidy

# Format go code
tool-fmt:
	@echo "---> Formatting code"
	go fmt ./...
.PHONY: tool-fmt

# Examine Go source code and reports suspicious constructs
tool-vet:
	@echo "---> Checking Go source code"
	go vet ./...
.PHONY: tool-vet

# Run application using linters: it runs linters in parallel, uses caching, supports yaml config, etc.
run-lint:
	@echo "---> Running linter"
	golangci-lint run ./... --timeout=3m
.PHONY: run-lint

# === DEVELOPMENT =======================================================
pre-commit: tool-mod-tidy tool-fmt build-lol-tactics run-lint tool-vet test
