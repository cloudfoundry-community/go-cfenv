[Breaking Changes]
- The minimum Go version is now **1.26** (was 1.25). Consumers on an
  older toolchain must stay on the v1.24 line. Go 1.25 is past upstream
  support, and the releases of `golang.org/x/net` that fix its current
  advisories declare `go 1.26.0`.

[Security Updates]
- `golang.org/x/net` moved to v0.62.0 and `golang.org/x/text` to
  v0.43.0, clearing the 3 known vulnerabilities their pinned versions
  carried (GO-2026-6860 and GO-2026-6866 in `x/net` v0.57.0,
  GO-2026-6629 in `x/text` v0.40.0). None were reachable from this
  package — both modules arrive only through `gomega` in the tests — so
  `govulncheck` stayed green, but scanners that report on presence
  rather than reachability did not, and consumers run those.
