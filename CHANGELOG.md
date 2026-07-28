# v1.19.0

# Breaking Changes

* Minimum Go version is now 1.18 (was 1.11).

  Required by the new mapstructure dependency, and by module graph
  pruning. Pruning means this library's test-only dependencies no
  longer propagate into consumers' module graphs and go.sum files.

  This is a floor for consumers, not a build constraint — the library
  is built and tested with the current Go release.

# Improvements

* Both archived dependencies removed.

  mitchellh/mapstructure (archived 2024-06-25) is replaced by
  github.com/go-viper/mapstructure/v2 v2.5.0, the maintained fork.
  v2.5.0 specifically: earlier v2 releases could echo decoded values
  into error messages, which matters for a VCAP_SERVICES parser.

  joefitzgerald/rainbow-reporter (archived 2026-02-23) is removed.
  It coloured test output and carried no assertion logic, so tests
  now use standard go test output.

  For consumers this means archived code leaves your dependency graph:
  go-cfenv was the only path to archived mapstructure in some builds.

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
