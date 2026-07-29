[Breaking Changes]
- The minimum Go version is now **1.25** (was 1.18). Consumers on an older toolchain must stay on v1.22.0. The old floor had become what blocked dependency maintenance: every current release of `golang.org/x/net`, `golang.org/x/text` and `gomega` declares a directive above it. (#38)

[Security Updates]
- `golang.org/x/net` and `golang.org/x/text` moved to current releases, clearing the 16 known vulnerabilities their pinned versions carried between them (15 in `x/net` v0.35.0, 1 in `x/text` v0.22.0). None were reachable from this package's code, so `govulncheck` was always green — but scanners that report on presence rather than reachability were not, and consumers run those. (#38)

[Chores]
- `gomega` moved to v1.42.1, which depends on the maintained `go.yaml.in/yaml/v3` instead of the archived `gopkg.in/yaml.v3`. The `.modrotignore` exemption that covered the archived module recorded exactly this as its removal condition, and is now gone — the archived-dependency gate runs with no exemptions at all.
- The `copyloopvar` and `intrange` linters now run. They were already enabled in `.golangci.yml` but golangci-lint disabled them at runtime for being newer than the project's go directive.
- Added `make fuzz`, a fuzz target over the `VCAP_APPLICATION` and `VCAP_SERVICES` parsing that `New` performs. Both documents come from the platform rather than the app, so they are the package's untrusted boundary. (#39)
