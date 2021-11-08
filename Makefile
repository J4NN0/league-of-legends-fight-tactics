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
	go build cmd/$(app_name)/main.go

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

# === RUN =======================================================
# Run application.
run:
	go run cmd/$(app_name)/main.go

# Run application using linters: it runs linters in parallel, uses caching, supports yaml config, etc.
run-lint:
	golangci-lint run ./...

# === TOOLS =======================================================
# Get a decorated HTML presentation of cover file: showing the covered (green), uncovered (red), and uninstrumented (grey) source.
tool-read-cover:
	go tool cover -html=$(cover_profile_filename)

# === DEVELOPMENT =======================================================
dev-pre-commit: build run-lint test-cover
