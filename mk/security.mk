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
#
# govulncheck, gosec and modrot are `go install`ed on demand; trivy and
# gitleaks are system packages, so a missing binary prints install
# instructions and fails rather than installing behind your back.

SECURITY_SCANS ?= govulncheck trivy gosec gitleaks
GOSEC_EXCLUDE  ?=
TRIVY_SCANNERS ?= vuln,secret
TRIVY_SEVERITY ?= HIGH,CRITICAL
MODROT_FLAGS   ?= --recursive --deprecated --resolve

.PHONY: security govulncheck trivy gosec gitleaks modrot

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
		echo "  Linux:  https://github.com/gitleaks/gitleaks#installing"; \
		exit 1; \
	}
	gitleaks detect --source . --no-banner --redact

modrot:
	@command -v modrot >/dev/null 2>&1 || { \
		echo "Installing modrot..."; \
		go install github.com/norman-abramovitz/modrot@latest; \
	}
	modrot $(MODROT_FLAGS) .
