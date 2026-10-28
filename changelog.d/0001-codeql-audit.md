[Chores]
- Maintainers: `make audit` now also runs CodeQL (`make codeql`) with
  the same `go-code-scanning` suite as the CodeQL workflow, and fails
  on any finding. It needs the CodeQL CLI locally
  (`brew install --cask codeql`); the database is built without the
  build tracer, so Apple Silicon needs no Rosetta. CI's audit job runs
  `make security modrot` instead, since `codeql.yml` already scans every
  PR. The target comes from the vendored `mk/security.mk`
  (GNUMakefile-Snippets).
