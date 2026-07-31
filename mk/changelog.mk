# changelog.mk — per-PR release-notes fragments (vendorable, standalone)
#
# Provides:
#   make changelog-new       Create changelog.d/NNNN-<slug>.md for this PR
#   make changelog-assemble  Print all fragments merged into release-notes
#                            layout on stdout (empty when there are none)
#   make changelog-sweep     git rm the fragments after a release
#   make changelog-check     Fail if any fragment is malformed or a stub
#   make changelog-deps      Draft the dependency fragment from the bot's
#                            commits since the last release tag
#   make changelog-deps-check  Report how many of those have landed, and
#                            warn when no fragment mentions any
#
# Settings:
#   CHANGELOG_DIR       Fragment directory. Default: changelog.d
#   CHANGELOG_SLUG      Slug for changelog-new. Default: current branch.
#   CHANGELOG_SECTIONS  Section names, in published order, '|'-separated.
#                       Default: Breaking Changes|Features|BugFixes|Chores|Security Updates
#   CHANGELOG_DEPS_GREP    Commit-subject pattern identifying a dependency
#                       bump (POSIX basic regex). Default: ^chore(deps
#   CHANGELOG_DEPS_SINCE   Start of the release window. Default: the
#                       nearest 'vX.Y.Z' or 'v-X.Y.Z' tag, i.e. the last
#                       release. Both forms are matched: a repo whose
#                       major version has outrun its Go module path tags
#                       releases 'v-X.Y.Z', because Go refuses to resolve
#                       a 'vN.Y.Z' tag for N>=2 without a matching '/vN'
#                       module path, but never version-parses the 'v-'
#                       form. Matching only 'v[0-9]*' there would find no
#                       tag and silently widen the window to the entire
#                       history rather than the last release.
#   CHANGELOG_DEPS_SLUG    Slug for changelog-deps. Default: dependency-updates
#   CHANGELOG_DEPS_SECTION Section the draft is filed under. Default: Chores
#
# Why fragments instead of one CHANGELOG file: a shared changelog is a
# merge-conflict magnet, and a single file cannot be both a durable
# history and the payload for one release — publishing it whole repeats
# every prior release in the new release's notes. A fragment per PR has
# neither problem, and the notes for a release are assembled from
# exactly the PRs that went into it.
#
# Pairs with go-release.mk: `make tag` embeds the assembled notes in the
# annotated tag body and `make publish` reads them back with
# --notes-from-tag, so nothing is hand-edited at release time. This file
# is useful on its own too — `make changelog-assemble > notes.md`.
#
# Fragment format — one or more sections, each with list items:
#
#     [Features]
#     - Added the frobnicator panel.
#
#     [BugFixes]
#     - Fixed the frobnicator crashing on empty input.

CHANGELOG_DIR      ?= changelog.d
CHANGELOG_SLUG     ?=
CHANGELOG_SECTIONS ?= Breaking Changes|Features|BugFixes|Chores|Security Updates

# The unbalanced '(' is deliberate and safe: this is a plain assignment, not
# a $(call) argument, and git's --grep is a POSIX basic regex where '(' is a
# literal. Do NOT "fix" it to \( — in a BRE that opens a group.
CHANGELOG_DEPS_GREP    ?= ^chore(deps
CHANGELOG_DEPS_SINCE   ?=
CHANGELOG_DEPS_SLUG    ?= dependency-updates
CHANGELOG_DEPS_SECTION ?= Chores

.PHONY: changelog-new changelog-assemble changelog-sweep changelog-check \
        changelog-deps changelog-deps-check

# Fragments sort by their numeric prefix. LC_ALL=C keeps that stable
# regardless of the caller's locale.
$(_HIDE)CHANGELOG_FIND = find "$(CHANGELOG_DIR)" -maxdepth 1 -name '[0-9]*.md' 2>/dev/null | LC_ALL=C sort

changelog-new:
	@set -e; \
	dir="$(CHANGELOG_DIR)"; \
	slug="$(CHANGELOG_SLUG)"; \
	[ -n "$$slug" ] || slug=$$(git branch --show-current 2>/dev/null || true); \
	slug=$$(printf '%s' "$$slug" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '-' | sed 's/^-*//; s/-*$$//'); \
	if [ -z "$$slug" ]; then \
		echo "ERROR: empty slug (detached HEAD?) - pass one: make changelog-new CHANGELOG_SLUG=<slug>" >&2; \
		exit 1; \
	fi; \
	mkdir -p "$$dir"; \
	max=0; \
	for f in $$($($(_HIDE)CHANGELOG_FIND)); do \
		n=$$(basename "$$f"); n=$$(( 10#$${n%%-*} )); \
		[ "$$n" -gt "$$max" ] && max=$$n || true; \
	done; \
	file=$$(printf '%s/%04d-%s.md' "$$dir" $$((max + 1)) "$$slug"); \
	if [ -e "$$file" ]; then echo "ERROR: $$file already exists" >&2; exit 1; fi; \
	printf '[Features]\n- \n' > "$$file"; \
	echo "$$file"

# Merges every fragment into one document, section by section, in
# CHANGELOG_SECTIONS order. An unknown section name or content sitting
# above the first [Section] header is a hard error: a malformed fragment
# must not silently vanish from a release's notes.
changelog-assemble:
	@set -e; \
	files=$$($($(_HIDE)CHANGELOG_FIND)); \
	[ -n "$$files" ] || exit 0; \
	awk -v order='$(CHANGELOG_SECTIONS)' ' \
	BEGIN { norder = split(order, sections, "|"); for (i = 1; i <= norder; i++) valid[sections[i]] = 1 } \
	FNR == 1 { section = ""; delete pending } \
	/^\[[^]]+\][[:space:]]*$$/ { \
		s = $$0; sub(/^\[/, "", s); sub(/\][[:space:]]*$$/, "", s); \
		if (!(s in valid)) { printf "ERROR: %s: unknown section [%s]\n", FILENAME, s > "/dev/stderr"; err = 1; exit 1 } \
		section = s; delete pending; next \
	} \
	{ \
		if (section == "") { \
			if ($$0 ~ /^[[:space:]]*$$/) next; \
			printf "ERROR: %s: content before any [Section] header\n", FILENAME > "/dev/stderr"; err = 1; exit 1 \
		} \
		if ($$0 ~ /^[[:space:]]*$$/) { pending[section] = pending[section] "\n"; next } \
		buf[section] = buf[section] pending[section] $$0 "\n"; delete pending[section] \
	} \
	END { \
		if (err) exit 1; \
		for (i = 1; i <= norder; i++) { \
			s = sections[i]; b = buf[s]; \
			sub(/^\n+/, "", b); sub(/\n+$$/, "\n", b); \
			if (b !~ /[^[:space:]]/) continue; \
			if (out) printf "\n"; \
			printf "[%s]\n%s", s, b; out = 1 \
		} \
	} \
	' $$files

# Guard for CI: catches the stub fragment `make changelog-new` writes,
# which is easy to add and forget to fill in.
changelog-check:
	@set -e; \
	files=$$($($(_HIDE)CHANGELOG_FIND)); \
	if [ -z "$$files" ]; then echo "changelog.d: no fragments"; exit 0; fi; \
	rc=0; \
	for f in $$files; do \
		if ! grep -qE '^-[[:space:]]+[^[:space:]]' "$$f"; then \
			echo "ERROR: $$f: no filled-in list item (still a stub?)" >&2; rc=1; \
		fi; \
	done; \
	$(MAKE) --no-print-directory changelog-assemble >/dev/null || rc=1; \
	[ "$$rc" -eq 0 ] && echo "changelog.d: $$(echo "$$files" | wc -l | tr -d ' ') fragment(s) OK" || exit 1

# ── Dependency updates ───────────────────────────────────────
#
# A bot that opens dependency PRs (Dependabot, Renovate) will never run
# changelog-new, so its bumps reach a release only if someone writes them up
# by hand — and nothing reports the omission until the notes come out thin.
# These two targets read them out of the commit log instead, keyed on the
# subject prefix the bot is configured to use. Pin that prefix in the bot's
# config; if it drifts, the bumps silently stop being found.
#
# Both see only what has LANDED on this branch. A bump sitting in an open
# PR is not counted.

# Resolves $$range and $$subjects: the window, and the deduped bump subjects
# in it with the prefix stripped, oldest first.
$(_HIDE)CHANGELOG_DEPS_SCAN = \
	since='$(CHANGELOG_DEPS_SINCE)'; \
	[ -n "$$since" ] || since=$$(git describe --tags --match 'v[0-9]*' --match 'v-[0-9]*' --abbrev=0 2>/dev/null || true); \
	range=$${since:+$$since..}HEAD; \
	subjects=$$(git log --reverse --no-merges --format='%s' --grep='$(CHANGELOG_DEPS_GREP)' "$$range" -- 2>/dev/null \
		| sed -E 's/^[^:]*: *(bump )?//I' \
		| awk 'NF && !seen[$$0]++')

# Writes a DRAFT, one bump per clause — rewrite it into prose before the
# release. Filed under Chores because that is where dependency work belongs;
# a security bump usually wants moving to [Security Updates] by hand.
changelog-deps:
	@set -e; \
	$($(_HIDE)CHANGELOG_DEPS_SCAN); \
	if [ -z "$$subjects" ]; then echo "$(CHANGELOG_DIR): no dependency bumps in $$range"; exit 0; fi; \
	file=$$($(MAKE) --no-print-directory changelog-new CHANGELOG_SLUG='$(CHANGELOG_DEPS_SLUG)'); \
	{ echo '[$(CHANGELOG_DEPS_SECTION)]'; \
	  printf -- '- Dependency updates: %s.\n' "$$(echo "$$subjects" | paste -sd ';' - | sed 's/;/; /g')"; \
	} > "$$file"; \
	echo "$$file"; \
	echo "Drafted from $$(echo "$$subjects" | wc -l | tr -d ' ') bump(s) in $$range - edit into prose before release." >&2

# Prints the count unconditionally, so this doubles as an any-time status
# command: "has enough dependency work piled up to be worth a build?" needs
# answering between releases, not only at tag time.
#
# It warns rather than fails, and is deliberately not wired into
# changelog-check. A hard error here would fire on pipeline fixes, docs-only
# changes, reverts and the bumps themselves; a gate with that false-positive
# rate gets routed around rather than obeyed.
#
# "Covered" is a loose text match, not an exact one against the subjects:
# changelog-deps writes a draft meant to be rewritten, and prose that no
# longer quotes the raw subjects would otherwise warn on every well-curated
# release. The cost is that a dependency fragment left over from an
# already-published window silences the report for the current one —
# changelog-sweep after each release is what prevents that.
changelog-deps-check:
	@set -e; \
	$($(_HIDE)CHANGELOG_DEPS_SCAN); \
	n=$$(test -n "$$subjects" && echo "$$subjects" | wc -l | tr -d ' ' || echo 0); \
	echo "$(CHANGELOG_DIR): $$n dependency bump(s) in $$range"; \
	[ "$$n" -gt 0 ] || exit 0; \
	files=$$($($(_HIDE)CHANGELOG_FIND)); \
	if [ -n "$$files" ] && echo "$$files" | xargs grep -qiE 'depend|bump' 2>/dev/null; then exit 0; fi; \
	echo "WARNING: none of them are mentioned in any fragment, so they will not" >&2; \
	echo "         appear in the release notes. Draft them: make changelog-deps" >&2

# Run after publishing. The removal is left staged-but-uncommitted so it
# rides the next PR rather than requiring a direct push to the default
# branch.
changelog-sweep:
	@set -e; \
	files=$$($($(_HIDE)CHANGELOG_FIND)); \
	if [ -z "$$files" ]; then echo "changelog.d: nothing to sweep"; exit 0; fi; \
	n=$$(echo "$$files" | wc -l | tr -d ' '); \
	$(if $(filter yes,$(DRYRUN)), \
		echo "DRYRUN: git rm $$(echo $$files | tr '\n' ' ')"; \
		echo "DRYRUN: would sweep $$n fragment(s).", \
		echo "$$files" | xargs git rm -q --; \
		echo "Swept $$n fragment(s) - commit this with the next PR.")
