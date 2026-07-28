# v1.19.1

No code changes. This release is repository and supply-chain tooling only;
the library is byte-for-byte identical to v1.19.0.

# Internal Changes/Improvements

* OpenSSF Scorecard now actually runs.

  The workflow referenced `ossf/scorecard-action@v2`. That action publishes
  no floating major tag, so the reference never resolved and every run since
  the workflow was added failed before starting — nothing was ever uploaded
  to the code-scanning tab. Pinned to v2.4.4.

* Dependabot's pinned-dependency guard now covers gomega.

  It documented the x/net and x/text pins but omitted onsi/gomega, which is
  held at v1.27.10 as the newest go1.18-compatible release. Dependabot duly
  proposed v1.42.1, which declares `go 1.25.0` and would have broken the
  declared floor.

* Added CODE_OF_CONDUCT.md, CONTRIBUTING.md and SECURITY.md, taken from the
  cloudfoundry and cloudfoundry-community org sources, so this repo points at
  the same CFF vulnerability reporting address and CLA process as the rest of
  the org.

* README carries CI, CodeQL, Scorecard, Go version, release and license
  badges. The Go Reference badge is removed. Go Report Card is deliberately
  absent: that project is archived read-only and its site now serves only a
  shutdown notice.

# Compatibility

No public API changes, and no changes to any `.go` file, `go.mod` or
`go.sum`. Consumers have no reason to upgrade from v1.19.0.

---

# v1.19.0

# Breaking Changes

* Minimum Go version is now 1.18 (was 1.11).

  Required by the new mapstructure dependency, and by module graph
  pruning. Pruning means this library's test-only dependencies no
  longer propagate into consumers' module graphs and go.sum files.

  This is a floor for consumers, not a build constraint — the library
  is built and tested with the current Go release.

# Improvements

* Both archived direct dependencies removed.

  mitchellh/mapstructure (archived 2024-06-25) is replaced by
  github.com/go-viper/mapstructure/v2 v2.5.0, the maintained fork.
  v2.5.0 specifically: earlier v2 releases could echo decoded values
  into error messages, which matters for a VCAP_SERVICES parser.

  joefitzgerald/rainbow-reporter (archived 2026-02-23) is removed.
  It coloured test output and carried no assertion logic, so tests
  now use standard go test output.

  For consumers this means archived code leaves your dependency graph:
  go-cfenv was the only path to archived mapstructure in some builds.

* One archived module remains, and it cannot reach your build.

  gopkg.in/yaml.v3 (archived 2025-04-01) is still present as an
  indirect, test-only dependency. Its only path is
  github.com/onsi/gomega/matchers, so it is not part of a consumer's
  build of this library, and nothing here imports it outside tests.

  It is not removable by upgrading: every published gomega release
  requires it. It is recorded in .modrotignore with that rationale, so
  the archived-dependency gate still fails on anything new.

* Dead ginkgo requirement dropped.

  Left over from the earlier move to sclevine/spec; it had no imports.

* Transitive x/net and x/text raised to their newest go1.18-compatible
  releases, clearing advisories reported by dependency scanners.

# Internal Changes/Improvements

* Makefile with build, test, quality, and audit targets.

* GitHub Actions replaces Travis CI, which had not run for years.
  Includes a compile-only job at the declared Go floor.

* golangci-lint, CodeQL code scanning, OpenSSF Scorecard, and
  Dependabot for both Go modules and workflow actions.

* modrot runs in the audit target, failing the build if an archived
  dependency is ever reintroduced.

# Compatibility

No public API changes. No exported identifiers added, removed, or
changed; error strings and JSON marshalling are unchanged.
