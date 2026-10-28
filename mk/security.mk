# security.mk — security scanning targets (vendorable, standalone)
#
# Provides:
#   make security      Run every scanner in SECURITY_SCANS
#   make govulncheck   Go dependency vulnerabilities (auto-installs)
#   make gosec         Go source security scanner (auto-installs)
#   make trivy         Filesystem vuln + secret scan (install prompted)
#   make gitleaks      Committed-secret scanner (install prompted)
#   make modrot        Archived/deprecated Go deps (auto-installs). Needs
#                      `gh auth token` + network, so it is NOT in the default
#                      SECURITY_SCANS — run it on its own.
#   make codeql        CodeQL static analysis of Go code (install prompted).
#                      Slow and a large install, so it is NOT in the
#                      default SECURITY_SCANS either.
#
# Settings:
#   SECURITY_SCANS     Scanners the aggregate target runs.
#                      Default: govulncheck trivy gosec gitleaks
#                      Non-Go repos drop the Go ones:
#                        SECURITY_SCANS := trivy gitleaks
#   GOSEC_EXCLUDE      Comma-separated gosec rule ids to exclude
#                      (e.g. G204,G304). Empty shows all findings.
#   TRIVY_SCANNERS     Trivy scanner set. Default: vuln,secret
#   TRIVY_SEVERITY     Trivy severities reported. Default: HIGH,CRITICAL
#   MODROT_FLAGS       Flags for the modrot scan.
#                      Default: --recursive --deprecated --resolve
#   GITLEAKS_VERSION   Version install-gitleaks fetches.
#   TOOLS_BIN          Where install-* put binaries. Default: ./bin
#   CODEQL_SUITE       CodeQL query suite. Default: the go-code-scanning
#                      suite github/codeql-action runs, so a local scan
#                      matches what code scanning reports.
#
# govulncheck, gosec and modrot are `go install`ed on demand; trivy,
# gitleaks and codeql are system packages, so a missing binary prints
# install instructions and fails rather than installing behind your back.
#
# That default is right for a workstation and wrong for CI: a plain
# runner has neither, so `make security` fails on a missing scanner
# rather than on a finding. `make install-gitleaks` fetches a pinned
# release and verifies it against the project's published checksums, so
# a CI job can run `make install-gitleaks security` without hand-rolling
# a download step or piping curl into a shell. Add TOOLS_BIN to PATH.

SECURITY_SCANS   ?= govulncheck trivy gosec gitleaks
GOSEC_EXCLUDE    ?=
TRIVY_SCANNERS   ?= vuln,secret
TRIVY_SEVERITY   ?= HIGH,CRITICAL
MODROT_FLAGS     ?= --recursive --deprecated --resolve
GITLEAKS_VERSION ?= 8.30.1
TOOLS_BIN        ?= $(CURDIR)/bin
CODEQL_SUITE     ?= codeql/go-queries:codeql-suites/go-code-scanning.qls

# sha256sum on Linux, shasum -a 256 on macOS/BSD.
$(_HIDE)TOOLS_SHA256SUM := $(shell command -v sha256sum >/dev/null 2>&1 && echo "sha256sum -c -" || echo "shasum -a 256 -c -")
TOOLS_SHA256SUM  ?= $($(_HIDE)TOOLS_SHA256SUM)

.PHONY: security govulncheck trivy gosec gitleaks modrot codeql install-gitleaks

security: $(SECURITY_SCANS)

govulncheck:
	@command -v govulncheck >/dev/null 2>&1 || { \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	}
	govulncheck ./...

gosec:
	@command -v gosec >/dev/null 2>&1 || { \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	}
	gosec -quiet $(if $(GOSEC_EXCLUDE),-exclude=$(GOSEC_EXCLUDE)) ./...

trivy:
	@command -v trivy >/dev/null 2>&1 || { \
		echo "trivy not found. Install:"; \
		echo "  macOS:  brew install trivy"; \
		echo "  Linux:  https://aquasecurity.github.io/trivy/latest/getting-started/installation/"; \
		exit 1; \
	}
	trivy fs --scanners $(TRIVY_SCANNERS) --severity $(TRIVY_SEVERITY) .

gitleaks:
	@command -v gitleaks >/dev/null 2>&1 || { \
		echo "gitleaks not found. Install:"; \
		echo "  macOS:  brew install gitleaks"; \
		echo "  Linux:  make install-gitleaks   (pinned release, checksum-verified)"; \
		echo "  or:     https://github.com/gitleaks/gitleaks#installing"; \
		exit 1; \
	}
	gitleaks detect --source . --no-banner --redact

# For CI. `gitleaks detect` reads git history, so a depth-1 checkout
# scans a single commit and reports clean whatever the history holds —
# fetch full depth in the job that runs this.
install-gitleaks:
	@set -e; \
	os=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	case "$$(uname -m)" in \
		x86_64|amd64) arch=x64 ;; \
		arm64|aarch64) arch=arm64 ;; \
		*) echo "ERROR: unsupported arch $$(uname -m)" >&2; exit 1 ;; \
	esac; \
	tarball="gitleaks_$(GITLEAKS_VERSION)_$${os}_$${arch}.tar.gz"; \
	base="https://github.com/gitleaks/gitleaks/releases/download/v$(GITLEAKS_VERSION)"; \
	tmp=$$(mktemp -d); trap 'rm -rf "$$tmp"' EXIT; \
	echo "Fetching $$tarball..."; \
	curl -fsSL -o "$$tmp/$$tarball" "$$base/$$tarball"; \
	curl -fsSL -o "$$tmp/checksums.txt" "$$base/gitleaks_$(GITLEAKS_VERSION)_checksums.txt"; \
	( cd "$$tmp" && grep " $$tarball$$" checksums.txt | $(TOOLS_SHA256SUM) ) \
		|| { echo "ERROR: checksum verification failed for $$tarball" >&2; exit 1; }; \
	tar -xzf "$$tmp/$$tarball" -C "$$tmp" gitleaks; \
	mkdir -p "$(TOOLS_BIN)"; \
	install -m 0755 "$$tmp/gitleaks" "$(TOOLS_BIN)/gitleaks"; \
	echo "Installed $$($(TOOLS_BIN)/gitleaks version) to $(TOOLS_BIN)/gitleaks"

modrot:
	@command -v modrot >/dev/null 2>&1 || { \
		echo "Installing modrot..."; \
		go install github.com/norman-abramovitz/modrot@latest; \
	}
	modrot $(MODROT_FLAGS) .

# The database is built with tracing off. `codeql database create`
# runs the Go autobuilder under CodeQL's build tracer, which fails with
# "Bad CPU type in executable" on Apple Silicon without Rosetta. Go
# extraction does not need the tracer: the autobuilder drives the
# extractor itself. analyze exits 0 whatever it finds, so the gate is
# the results file — CodeQL's CSV has no header, so empty means clean.
codeql:
	@command -v codeql >/dev/null 2>&1 || { \
		echo "codeql not found. Install:"; \
		echo "  macOS:  brew install --cask codeql"; \
		echo "  Linux:  https://github.com/github/codeql-action/releases (codeql-bundle)"; \
		exit 1; \
	}
	@set -e; \
	tmp=$$(mktemp -d); trap 'rm -rf "$$tmp"' EXIT; \
	log="$$tmp/codeql.log"; \
	{ codeql database init "$$tmp/db" --language=go --source-root=. && \
	  codeql database trace-command --no-tracing "$$tmp/db" -- \
		"$$(codeql resolve extractor --language=go)/tools/autobuild.sh" && \
	  codeql database finalize "$$tmp/db" && \
	  codeql database analyze "$$tmp/db" $(CODEQL_SUITE) --download \
		--format=csv --output="$$tmp/results.csv"; \
	} >"$$log" 2>&1 || { cat "$$log" >&2; echo "ERROR: codeql scan failed" >&2; exit 1; }; \
	if [ -s "$$tmp/results.csv" ]; then \
		cat "$$tmp/results.csv"; \
		echo "codeql: $$(wc -l < "$$tmp/results.csv" | tr -d ' ') finding(s)" >&2; \
		exit 1; \
	fi; \
	echo "codeql: no findings"
