# go-cfenv — build, test, and release automation.

PROJECT := go-cfenv

# This is a library: no binaries are produced. Both overrides are required,
# not cosmetic. go-release.mk defaults GH_ASSETS to a glob over built
# artifacts; with none, the unmatched glob reaches gh as a literal '*' and
# the release fails. TARGETS empty suppresses the cross-compile fan-out.
TARGETS   :=
GH_ASSETS :=

SECURITY_SCANS := govulncheck gosec gitleaks

include mk/help.mk
include mk/version.mk
include mk/bump.mk
include mk/security.mk
include mk/go-release.mk

##@ Build

.PHONY: build
build: ## Compile all packages
	go build ./...

##@ Test

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-race
test-race: ## Run tests with the race detector
	go test -race ./...

.PHONY: cover
cover: ## Run tests and report coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

##@ Quality

.PHONY: fmt
fmt: ## Report files needing gofmt
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { \
	  echo "golangci-lint not found. See https://golangci-lint.run/welcome/install/"; exit 1; }
	golangci-lint run ./...

.PHONY: tidy
tidy: ## Fail if go mod tidy would change anything
	@cp go.mod go.mod.bak; cp go.sum go.sum.bak
	@go mod tidy || { mv go.mod.bak go.mod; mv go.sum.bak go.sum; exit 1; }
	@if ! diff -q go.mod go.mod.bak >/dev/null || ! diff -q go.sum go.sum.bak >/dev/null; then \
	  mv go.mod.bak go.mod; mv go.sum.bak go.sum; \
	  echo "go mod tidy would modify go.mod/go.sum — commit the tidied files"; exit 1; \
	fi
	@rm -f go.mod.bak go.sum.bak

.PHONY: check
check: fmt vet lint tidy ## Run all quality gates

##@ Security

.PHONY: audit
audit: security modrot ## Run every scanner, including the archived-dependency check
