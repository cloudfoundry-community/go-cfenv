# changelog.d — release-notes fragments

Each PR carries its own release-notes fragment here. At release time the
fragments are assembled into the annotated tag body (`make tag`),
published as the GitHub release notes (`make publish`), and then removed
(`make changelog-sweep`) — so this directory is empty right after every
release. Only this README is permanent.

## Adding a fragment to your PR

```bash
make changelog-new                          # names it after your branch
make changelog-new CHANGELOG_SLUG=my-slug   # or pick a slug
```

That creates `NNNN-<slug>.md`. Write your entry under the section(s) it
belongs to — one file can feed several:

```markdown
[Features]
- Added the frobnicator panel.

[BugFixes]
- Fixed the frobnicator crashing on empty input.
```

Valid sections, in the order they appear in published notes:

| Section | Use for |
|---------|---------|
| `[Breaking Changes]` | Anything a consumer must act on before upgrading |
| `[Features]` | New functionality |
| `[BugFixes]` | Fixes |
| `[Chores]` | Dependency bumps, refactors, CI/build changes |
| `[Security Updates]` | CVE fixes |

`make changelog-check` validates them, and CI runs it. An unknown section
name, content above the first header, or an unfilled stub is an error —
a malformed fragment must not silently vanish from a release.

## Why fragments instead of a CHANGELOG file

A single shared file is a merge-conflict magnet, and it cannot be both a
durable history and the payload for one release: publishing it whole
repeats every prior release inside the new release's notes. Release notes
for v1.19.0 and earlier live in the [GitHub releases](https://github.com/cloudfoundry-community/go-cfenv/releases).
