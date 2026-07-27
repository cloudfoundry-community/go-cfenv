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
#   TAG          Release tag. Default: the version with build metadata
#                stripped — tags are clean semver (vX.Y.Z[-prerelease]).
#   TAG_REMOTE   Remote for tag/untag. Default: origin
#   DRAFT        DRAFT=yes publishes a draft release.
#   NOTES        Notes file for publish. Default: gh --generate-notes.
#   GH_ASSETS    Files publish uploads. Default: RELEASE_ROOT/PROJECT-*
#                (artifacts + their sibling checksum files).
#   DRYRUN       DRYRUN=yes echoes state-changing commands instead.
#
# gh authenticates from the environment (GH_TOKEN/GITHUB_TOKEN or a
# prior `gh auth login`) — credentials never appear on a command line.
# --prerelease derives from an alpha/beta/rc part in the tag — dev.N
# tags can be full releases. Rollback order: unpublish first, then
# untag, so a half-done rollback never orphans the tag.

TAG        ?= v$($(_HIDE)SEMVER_NOMETA)
TAG_REMOTE ?= origin
DRAFT      ?=
NOTES      ?=
GH_ASSETS  ?= $(RELEASE_ROOT)/$(PROJECT)-*

# Prefix that turns state-changing commands into echoes under DRYRUN=yes
$(_HIDE)DRY := $(if $(filter yes,$(DRYRUN)),@echo "DRYRUN:" )

.PHONY: tag untag publish unpublish

tag:
	@case "$(TAG)" in v[0-9]*.[0-9]*.[0-9]*) ;; *) echo "ERROR: '$(TAG)' does not look like a release tag (vX.Y.Z[-prerelease])" >&2; exit 1;; esac
	$($(_HIDE)DRY)git tag -a "$(TAG)" -m "Release $(TAG)"
	$($(_HIDE)DRY)git push $(TAG_REMOTE) "refs/tags/$(TAG)"

untag:
	@echo "Deleting tag $(TAG) locally and on $(TAG_REMOTE)..."
	-$($(_HIDE)DRY)git tag -d "$(TAG)"
	$($(_HIDE)DRY)git push $(TAG_REMOTE) --delete "refs/tags/$(TAG)"

publish:
	@TAG="$(TAG)"; \
	PRERELEASE=""; case "$$TAG" in *-alpha*|*-beta*|*-rc*) PRERELEASE="--prerelease";; esac; \
	set -- gh release create "$$TAG" --title "$(PROJECT) $$TAG" --verify-tag $$PRERELEASE $(if $(filter yes,$(DRAFT)),--draft) $(if $(NOTES),--notes-file "$(NOTES)",--generate-notes) $(GH_ASSETS); \
	$(if $(filter yes,$(DRYRUN)),echo "DRYRUN: $$*",echo "+ $$*"; "$$@")

unpublish:
	@echo "Release to delete from GitHub:"
	@gh release view "$(TAG)" --json tagName,name,isDraft,assets \
		--jq '"  " + .tagName + "  (" + .name + ")" + (if .isDraft then "  [draft]" else "" end), (.assets[] | "    " + .name)'
	$($(_HIDE)DRY)gh release delete "$(TAG)" --yes
