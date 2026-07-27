# help.mk — self-documenting help target (vendorable, standalone)
#
# `make help` lists every target annotated with a `##` comment, in
# declaration order, grouped by `##@ Section` headers:
#
#   ##@ Build
#   build: ## Build the binary for the current platform
#   install: build ## Build and install
#
# Targets without a `##` comment stay out of the listing, so internal
# targets need no other hiding. Annotations are scraped from every
# included makefile ($(MAKEFILE_LIST)) — annotated snippet targets
# show up automatically.
#
# Make it the default goal in the including Makefile if wanted:
#   .DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_.%-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[33m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
