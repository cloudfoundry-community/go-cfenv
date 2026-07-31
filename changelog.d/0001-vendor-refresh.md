[Chores]
- Refreshed the vendored build tooling: `make publish` now resolves the
  tag that exists rather than deriving the next one, the deletion verbs
  require an explicit tag, and dependency bumps landed since the last
  release can be drafted into release notes with `make changelog-deps`.
