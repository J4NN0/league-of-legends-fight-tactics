# Makefile
#
# - config - tooling versions and variable configs
# - build - build application
# - test - run tests
# - run - run application
# - tools - golang useful tools
# - development - scripts for development

# === CONFIG =======================================================
app_name := league-of-legends-fight-tactics
cover_profile_filename := build/cover.out


# === BUILD =======================================================
build:
	go build -o build/league_of_legends_fight_tactics cmd/$(app_name)/main.go

# === TEST =======================================================
# Run all tests.
test:
	go test ./...

# Run all tests with verbose output.
test-verbose:
	go test -v ./...

# Test how much of a package’s code is exercised by running the package’s tests.
test-cover:
	go test -cover -coverprofile=$(cover_profile_filename) ./...

# Test data race
test-race:
	go test -race ./...

# Test data race and package coverage
test-pre-commit:
	go test -race -cover -coverprofile=$(cover_profile_filename) ./...

# === RUN =======================================================
# Run application.
run:
	go run cmd/$(app_name)/main.go -c1 $(c1) -c2 $(c2)

run-all:
	go run cmd/$(app_name)/main.go -all=true

# Run application using linters: it runs linters in parallel, uses caching, supports yaml config, etc.
run-lint:
	golangci-lint run ./...

# === TOOLS =======================================================
# Get a decorated HTML presentation of cover file: showing the covered (green), uncovered (red), and un-instrumented (grey) source.
tool-read-cover:
	go tool cover -html=$(cover_profile_filename)

# Fix go.mod and go.sum
tool-mod-tidy:
	go mod tidy

# Format go code
tool-fmt:
	go fmt ./...

# === DEVELOPMENT =======================================================
dev-pre-commit: tool-mod-tidy tool-fmt build run-lint test-pre-commit
