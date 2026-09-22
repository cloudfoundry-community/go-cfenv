[Chores]
- Bumped `github.com/onsi/gomega` from v1.42.1 to v1.43.1, and the
  pinned `github/codeql-action` used by the CodeQL and Scorecard
  workflows from v4.37.7 to v4.38.1.
- Dependabot now prefixes its commit subjects with `chore(deps)`, so
  `make changelog-deps` finds the bumps in a release window instead of
  reporting none. Earlier bumps had to be written up by hand.
- Recorded in SECURITY.md why the OpenSSF Scorecard `Code-Review` check
  scores 0 for a single-maintainer library, and what is enforced on
  `master` in place of a second reviewer.
