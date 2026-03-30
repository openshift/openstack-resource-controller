## --------------------------------------
## Openshift specific make targets,
## intended to be included in root Makefile in this repository along with openshift folder.
## --------------------------------------

OPENSHIFT_DIR=./openshift

verify-generated: generate-openshift

.PHONY: generate-openshift
generate-openshift:
	$(MAKE) -C $(OPENSHIFT_DIR) ocp-manifests

.PHONY: merge-bot
merge-bot: full-vendoring generate generate-openshift ## Runs targets that help merge-bot to rebase downstream ORC.

.PHONY: full-vendoring ## Runs commands that complete vendoring tasks for downstream ORC.
full-vendoring:
	go mod tidy && go mod vendor
