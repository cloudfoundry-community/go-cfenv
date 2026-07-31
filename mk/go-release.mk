# go-release.mk — cross-compiled Go release fan-out (vendorable)
#
# Generates one release target per os/arch pair and the surrounding
# lifecycle:
#   make release-all         clean, build every TARGETS entry, list
#   make ci-release          release-all, but VERSION= must be given
#   make show-releases       list built artifacts
#   make release-clean       remove built artifacts
#
# Plus the GitHub release lifecycle (tag → publish / unpublish → untag):
#   make tag [VERSION=vX]    create + push the annotated release tag
#   make publish [DRAFT=yes] gh release create + upload artifacts
#   make unpublish TAG=vX    delete the GitHub release (assets included)
#   make untag TAG=vX        delete the tag (local + remote)
#
# Artifacts land in RELEASE_ROOT as
#   <PROJECT>-<version>+<os>.<arch>[.<meta>][.exe]
# each with a sibling checksum file (.sha1/.sha256/...).
#
# Settings:
#   PROJECT            Artifact base name. Default: directory name.
#   TARGETS            os/arch pairs. Default:
#                      linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
#   RELEASE_ROOT       Output directory. Default: releases
#   RELEASE_VERSION    Version in artifact names. Default: version.mk's
#                      SEMVER_VERSION (include version.mk first, or set
#                      this explicitly).
#   RELEASE_META       Optional extra artifact-name suffix (typically
#                      build metadata). Default: empty.
#   RELEASE_PACKAGES   Package path(s) passed to go build. Default: .
#   RELEASE_LDFLAGS    ldflags template, expanded per target with
#                      $(1)=os and $(2)=arch — so per-platform symbols
#                      work:
#                        RELEASE_LDFLAGS = $(GO_LDFLAGS) \
#                          -X 'main.GoOs=$(1)' -X 'main.GoArch=$(2)'
#                      Default: $(GO_LDFLAGS) (yours to define).
#   RELEASE_CHECKSUM   Checksum flavor, shaN form. Default: sha256
#                      (one flavor; add a second checksum in the repo
#                      if a distribution channel demands it)
#   CGO_ENABLED        Default: 0 (static-friendly cross builds)
#
# Windows artifacts get .exe; the checksum file sits next to the
# artifact under the un-suffixed name.

_HIDE ?= _

PROJECT          ?= $(notdir $(CURDIR))
TARGETS          ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
RELEASE_ROOT     ?= releases
RELEASE_VERSION  ?= $($(_HIDE)SEMVER_VERSION)
RELEASE_META     ?=
RELEASE_PACKAGES ?= .
RELEASE_LDFLAGS  ?= $(GO_LDFLAGS)
RELEASE_CHECKSUM ?= sha256
CGO_ENABLED      ?= 0

$(_HIDE)RELEASES := $(foreach target,$(TARGETS),release-$(target)-$(PROJECT))

.PHONY: release-all ci-release show-releases release-clean $(_HIDE)release-mkdir

release-all: release-clean $(_HIDE)release-mkdir $($(_HIDE)RELEASES) show-releases

# Requires VERSION given explicitly on the command line — a VERSION
# resolved by version.mk doesn't count, releases must be deliberate.
ci-release:
	@if [ -z "$(filter command line environment,$(origin VERSION))" ]; then \
		echo "VERSION must be set explicitly: make ci-release VERSION=x.y.z" >&2; \
		exit 1; \
	fi
	@$(MAKE) release-all

show-releases:
	@ls -lA $(RELEASE_ROOT)

release-clean:
	@rm -f $(RELEASE_ROOT)/$(PROJECT)-* || true
	@[ ! -d $(RELEASE_ROOT) ] || rmdir -p $(RELEASE_ROOT) 2>/dev/null || true

$(_HIDE)release-mkdir:
	@mkdir -p $(RELEASE_ROOT)

# One target per os/arch pair. Target-specific variables stay lazy (=)
# so RELEASE_VERSION resolves at build time (bump/build chains).
define $(_HIDE)release_target_impl
.PHONY: release-$(1)/$(2)-$(PROJECT)
release-$(1)/$(2)-$(PROJECT): $(_HIDE)REL_LDFLAGS = $$(call RELEASE_LDFLAGS,$(1),$(2))
release-$(1)/$(2)-$(PROJECT): $(_HIDE)REL_BASE = $$(RELEASE_ROOT)/$$(PROJECT)-$$(RELEASE_VERSION)+$(1).$(2)$$(if $$(RELEASE_META),.$$(RELEASE_META))
release-$(1)/$(2)-$(PROJECT): $(_HIDE)REL_EXE = $$($(_HIDE)REL_BASE)$(if $(patsubst windows,,$(1)),,.exe)
release-$(1)/$(2)-$(PROJECT):
	@echo "Building $$(PROJECT) $$(RELEASE_VERSION) for $(1)/$(2)..."
	@CGO_ENABLED=$$(CGO_ENABLED) GOOS=$(1) GOARCH=$(2) go build -o "$$($(_HIDE)REL_EXE)" -ldflags "$$($(_HIDE)REL_LDFLAGS)" $$(RELEASE_PACKAGES)
	@shasum -a $$(patsubst sha%,%,$$(RELEASE_CHECKSUM)) "$$($(_HIDE)REL_EXE)" > "$$($(_HIDE)REL_BASE).$$(RELEASE_CHECKSUM)"
endef

$(foreach target,$(TARGETS),$(eval $(call $(_HIDE)release_target_impl,$(word 1,$(subst /, ,$(target))),$(word 2,$(subst /, ,$(target))))))

# ── Release lifecycle (tag → publish / unpublish → untag) ─────
# Settings:
#   TAG          Release tag. Explicitly set, it applies to every verb.
#                Unset, each verb resolves what its intent needs: tag
#                derives the NEXT version (build metadata stripped —
#                tags are clean semver); publish takes the nearest
#                EXISTING tag; untag/unpublish refuse to guess and
#                require TAG= on the command line.
#   TAG_REMOTE   Remote for tag/untag. Default: origin
#   DRAFT        DRAFT=yes publishes a draft release.
#   NOTES        Notes file for publish. Default: gh --generate-notes.
#   TAG_NOTES_CMD
#                Command whose stdout becomes the annotated tag body.
#                Empty (default) tags with "Release <TAG>". Set it to
#                `$(MAKE) changelog-assemble` to have changelog.mk's
#                fragments become the tag body — then publish reads them
#                straight back and nothing is hand-edited at release time.
#   GH_ASSETS    Files publish uploads. Default: RELEASE_ROOT/PROJECT-*
#                (artifacts + their sibling checksum files).
#   DRYRUN       DRYRUN=yes echoes state-changing commands instead.
#
# gh authenticates from the environment (GH_TOKEN/GITHUB_TOKEN or a
# prior `gh auth login`) — credentials never appear on a command line.
# --prerelease derives from an alpha/beta/rc part in the tag — dev.N
# tags can be full releases. Rollback order: unpublish first, then
# untag, so a half-done rollback never orphans the tag.
#
# Notes precedence in publish: NOTES (a file) wins if set; otherwise, if
# the tag carries a body of its own, --notes-from-tag uses it; otherwise
# gh generates notes from the commit log. So a repo using changelog.mk
# gets fragment-derived notes with no extra arguments, and a repo using
# neither is unaffected.

# TAG_PREFIX selects the tag naming scheme. 'v' is the historical default.
# Set 'v-' in repos whose major version has outrun their Go module path:
# Go version-parses a 'vN.Y.Z' tag and refuses it when N>=2 without a
# matching '/vN' module path, but never version-parses 'v-N.Y.Z', so that
# form stays resolvable as a `go get` revision. version.mk parses both.
TAG_PREFIX    ?= v
TAG           ?=
TAG_REMOTE    ?= origin
DRAFT         ?=
NOTES         ?=
TAG_NOTES_CMD ?=
GH_ASSETS     ?= $(RELEASE_ROOT)/$(PROJECT)-*

# Prefix that turns state-changing commands into echoes under DRYRUN=yes
$(_HIDE)DRY := $(if $(filter yes,$(DRYRUN)),@echo "DRYRUN:" )

# Per-verb tag resolution. tag CREATES the next release, so its default
# derives the next version. publish operates on a tag that already
# EXISTS — deriving "next" there is always one release ahead: the moment
# tag has run, next-version points past the tag just created, and a bare
# `make tag && make publish` aborts on --verify-tag having published
# nothing. So publish defaults to the nearest existing tag instead,
# resolved at recipe time (kept lazy) so it also sees a tag created
# earlier in the same invocation. Both stay recursive on purpose — see
# version.mk on why := would freeze pre-bump values.
$(_HIDE)NEXT_TAG    = $(TAG_PREFIX)$($(_HIDE)SEMVER_NOMETA)
$(_HIDE)LAST_TAG    = $(shell git describe --tags --abbrev=0 2>/dev/null)
$(_HIDE)CREATE_TAG  = $(or $(TAG),$($(_HIDE)NEXT_TAG))
$(_HIDE)PUBLISH_TAG = $(or $(TAG),$($(_HIDE)LAST_TAG))

.PHONY: tag untag publish unpublish

tag:
	@case "$($(_HIDE)CREATE_TAG)" in v[0-9]*.[0-9]*.[0-9]*) ;; *) echo "ERROR: '$($(_HIDE)CREATE_TAG)' does not look like a release tag (vX.Y.Z[-prerelease])" >&2; exit 1;; esac
ifeq ($(strip $(TAG_NOTES_CMD)),)
	$($(_HIDE)DRY)git tag -a "$($(_HIDE)CREATE_TAG)" -m "Release $($(_HIDE)CREATE_TAG)"
else
	@body=$$($(TAG_NOTES_CMD)) || { echo "ERROR: TAG_NOTES_CMD failed - refusing to tag" >&2; exit 1; }; \
	if [ -z "$$body" ]; then \
		echo "ERROR: TAG_NOTES_CMD produced no output - refusing to tag $($(_HIDE)CREATE_TAG) with empty release notes." >&2; \
		echo "       Add a fragment (make changelog-new), or unset TAG_NOTES_CMD for an untitled release." >&2; \
		exit 1; \
	fi; \
	$(if $(filter yes,$(DRYRUN)), \
		printf 'DRYRUN: git tag -a %s -F - <<EOF\n%s\nEOF\n' "$($(_HIDE)CREATE_TAG)" "$$body", \
		printf '%s\n\n%s\n' "Release $($(_HIDE)CREATE_TAG)" "$$body" | git tag -a "$($(_HIDE)CREATE_TAG)" -F -)
endif
	$($(_HIDE)DRY)git push $(TAG_REMOTE) "refs/tags/$($(_HIDE)CREATE_TAG)"

# Deletion verbs never guess: "whatever is newest" as a default is how a
# typo removes the wrong release. Both demand an explicit TAG.
untag:
	@[ -n "$(TAG)" ] || { echo "ERROR: untag deletes a tag - name it: make untag TAG=vX.Y.Z" >&2; exit 1; }
	@echo "Deleting tag $(TAG) locally and on $(TAG_REMOTE)..."
	-$($(_HIDE)DRY)git tag -d "$(TAG)"
	$($(_HIDE)DRY)git push $(TAG_REMOTE) --delete "refs/tags/$(TAG)"

publish:
	@TAG="$($(_HIDE)PUBLISH_TAG)"; \
	if [ -z "$$TAG" ]; then \
		echo "ERROR: no tag to publish - create one with 'make tag' or pass TAG=vX.Y.Z" >&2; \
		exit 1; \
	fi; \
	PRERELEASE=""; case "$$TAG" in *-alpha*|*-beta*|*-rc*) PRERELEASE="--prerelease";; esac; \
	NOTESARG="--generate-notes"; \
	$(if $(NOTES),NOTESARG="--notes-file $(NOTES)",\
		if [ -n "$$(git tag -l --format='%(contents:body)' "$$TAG" 2>/dev/null | tr -d '[:space:]')" ]; then \
			NOTESARG="--notes-from-tag"; \
		fi); \
	set -- gh release create "$$TAG" --title "$(PROJECT) $$TAG" --verify-tag $$PRERELEASE $(if $(filter yes,$(DRAFT)),--draft) $$NOTESARG $(GH_ASSETS); \
	$(if $(filter yes,$(DRYRUN)),echo "DRYRUN: $$*",echo "+ $$*"; "$$@")

unpublish:
	@[ -n "$(TAG)" ] || { echo "ERROR: unpublish deletes a release - name it: make unpublish TAG=vX.Y.Z" >&2; exit 1; }
	@echo "Release to delete from GitHub:"
	@gh release view "$(TAG)" --json tagName,name,isDraft,assets \
		--jq '"  " + .tagName + "  (" + .name + ")" + (if .isDraft then "  [draft]" else "" end), (.assets[] | "    " + .name)'
	$($(_HIDE)DRY)gh release delete "$(TAG)" --yes
