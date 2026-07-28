[Chores]
- Release notes are now assembled from per-PR `changelog.d/` fragments into the annotated tag body, so cutting a release no longer requires hand-editing a notes file. `CHANGELOG.md` is removed; notes for v1.19.1 and earlier remain in the GitHub releases.
- The CI security job installs gitleaks with `make install-gitleaks` — a pinned, checksum-verified download from the shared make snippets — replacing the download step that was inlined in the workflow.
