[Chores]
- Bumped the pinned `github/codeql-action` used by the CodeQL and
  Scorecard workflows from v4.39.1 to v4.39.5 (#65, #66, #67),
  `actions/setup-go` used by the CI and CodeQL workflows from v7.0.0 to
  v7.1.0 (#66), and `ossf/scorecard-action` from v2.5.0 to v2.6.0 (#66).
- Raised the indirect module requirements `golang.org/x/net` from
  v0.62.0 to v0.63.0, `golang.org/x/text` from v0.43.0 to v0.44.0, and
  `go.yaml.in/yaml/v3` from v3.0.5 to v3.0.6. None of them fixes a known
  vulnerability, and the `go` directive stays at 1.26.0.
