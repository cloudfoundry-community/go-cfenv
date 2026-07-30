# bump.mk — version maintenance (vendorable, standalone)
#
# Bumps the repo's persisted current version:
#   make bump major     x.y.z -> (x+1).0.0, prerelease stripped
#   make bump minor     x.y.z -> x.(y+1).0, prerelease stripped
#   make bump patch     x.y.z -> x.y.(z+1), prerelease stripped
#   make bump <label>   increment '<label>.N' prerelease counter, or
#                       start '<label>.1' when coming from a different
#                       label. Labels come from BUMP_LABELS. From a
#                       FINAL version the core advances first — see
#                       BUMP_PRERELEASE_STEP.
#   make bump final     strip the prerelease (x.y.z-rc.2 -> x.y.z)
#
# Settings:
#   BUMP_LABELS        Prerelease labels usable as bump modifiers.
#                      Default: dev alpha beta rc
#   BUMP_PRERELEASE_STEP
#                      Which core component advances when a label is
#                      applied to a FINAL version: major, minor or patch.
#                      Default: minor.
#
#                      Semver orders a prerelease BEFORE its own release,
#                      so 1.2.3 -> 1.2.3-dev.1 would be a version older
#                      than the release it follows: unreachable, and a
#                      Go module proxy would never select it. The label
#                      therefore has to target a later core. Set this to
#                      whatever the project's usual release step is, so
#                      `bump final` lands on the number actually wanted:
#                        BUMP_PRERELEASE_STEP := patch
#                      Only applies coming from a final version; counter
#                      advances and label switches leave the core alone.
#   BUMP_BASE_CMD      Command printing the persisted CURRENT version
#                      (not version.mk's "next release" resolution —
#                      bumping that would double-increment). Default:
#                      the most recently created tag reachable from
#                      HEAD, verbatim — creation-date order so chains
#                      of bumps on one commit (bump dev, bump dev)
#                      advance correctly, which `git describe` cannot
#                      do with several tags on the same commit. A
#                      leading 'v' is preserved onto the new version.
#                      No tag at all bumps from 0.0.0. File-versioned
#                      repos point this at their file, e.g.:
#                        BUMP_BASE_CMD := cat VERSION
#   VERSION_WRITE_CMD  Command persisting the new version; invoked with
#                      the full version string as its argument.
#                      Default: git tag -a -m Bumped — annotated
#                      because lightweight tags carry no creation date
#                      of their own, which would break same-commit
#                      bump chains under BUMP_BASE_CMD's ordering.
#   DRYRUN=yes         Print the new version without persisting.
#
# Build metadata is not carried through a bump — VERSION_WRITE_CMD may
# stamp its own. With the git-tag defaults, `make bump patch` creates
# exactly the version version.mk's next-patch default already resolves.

_HIDE ?= _

BUMP_LABELS       ?= dev alpha beta rc
BUMP_PRERELEASE_STEP ?= minor
# Last --sort key is primary: creation date desc, version-aware refname
# desc as the same-second tiebreak.
BUMP_BASE_CMD     ?= git tag --merged HEAD --sort=-v:refname --sort=-creatordate 2>/dev/null | head -n 1
VERSION_WRITE_CMD ?= git tag -a -m Bumped

$(_HIDE)BUMP_MOD := $(filter major minor patch final $(BUMP_LABELS),$(MAKECMDGOALS))

# No-op targets so bump modifiers don't error as unknown goals.
.PHONY: major minor patch final $(BUMP_LABELS)
major minor patch final $(BUMP_LABELS):
	@:

.PHONY: bump
bump:
	@set -- $($(_HIDE)BUMP_MOD); \
	if [ $$# -eq 0 ]; then \
		echo "Usage: make bump <major|minor|patch|final|$(subst $() ,|,$(BUMP_LABELS))>" >&2; \
		exit 1; \
	elif [ $$# -gt 1 ]; then \
		echo "Only one bump modifier allowed" >&2; \
		exit 1; \
	fi
	@base=$$($(BUMP_BASE_CMD)); base=$${base:-0.0.0}; \
	prefix=; case "$$base" in v*) prefix=v; base=$${base#v} ;; esac; \
	base=$${base%%+*}; \
	core=$${base%%-*}; \
	pre=; case "$$base" in *-*) pre=$${base#*-} ;; esac; \
	maj=$${core%%.*}; rest=$${core#*.}; min=$${rest%%.*}; pat=$${rest#*.}; \
	case "$$maj$$min$$pat" in *[!0-9]*|'') \
		echo "Cannot bump malformed current version '$$prefix$$base'" >&2; exit 1 ;; \
	esac; \
	mod='$($(_HIDE)BUMP_MOD)'; \
	case "$$mod" in \
		major) maj=$$((maj + 1)); min=0; pat=0; pre= ;; \
		minor) min=$$((min + 1)); pat=0; pre= ;; \
		patch) pat=$$((pat + 1)); pre= ;; \
		final) \
			if [ -z "$$pre" ]; then \
				echo "Already a final version: $$prefix$$base" >&2; exit 1; \
			fi; pre= ;; \
		*) \
			if [ -z "$$pre" ]; then \
				case "$(BUMP_PRERELEASE_STEP)" in \
					major) maj=$$((maj + 1)); min=0; pat=0 ;; \
					minor) min=$$((min + 1)); pat=0 ;; \
					patch) pat=$$((pat + 1)) ;; \
					*) echo "BUMP_PRERELEASE_STEP must be major, minor or patch, not '$(BUMP_PRERELEASE_STEP)'" >&2; exit 1 ;; \
				esac; \
			fi; \
			case "$$pre" in \
				"$$mod".*) n=$${pre#"$$mod".}; \
					case "$$n" in *[!0-9]*|'') n=0 ;; esac; \
					pre="$$mod.$$((n + 1))" ;; \
				*) pre="$$mod.1" ;; \
			esac ;; \
	esac; \
	new="$$prefix$$maj.$$min.$$pat$${pre:+-$$pre}"; \
	if [ "$(DRYRUN)" = "yes" ]; then \
		echo "DRYRUN: $$prefix$$base -> $$new"; \
	else \
		$(VERSION_WRITE_CMD) "$$new" && \
		echo "Bumped: $$prefix$$base -> $$new"; \
	fi
