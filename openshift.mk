## --------------------------------------
## Openshift specific make targets,
## intended to be included in root Makefile in this repository along with openshift folder.
## --------------------------------------

OPENSHIFT_DIR=./openshift
manifests_dir ?= $(OPENSHIFT_DIR)/manifests
manifests_prefix ?= 0000_30_openstack-resource-controller_

define manifest_name
    $(addsuffix ".yaml",$(addprefix $(manifests_dir)/$(manifests_prefix),$(1)))
endef

manifest_names = 04_infrastructure-components
infrastructure_components = $(OPENSHIFT_DIR)/cluster-capi-configmap/infrastructure-components.yaml

verify-generated: generate-openshift

.PHONY: generate-openshift
generate-openshift: $(foreach m,$(manifest_names),$(call manifest_name,$(m)))

$(infrastructure_components): $(KUSTOMIZE) ALWAYS
	$(KUSTOMIZE) build $(OPENSHIFT_DIR)/infrastructure-components > $@

$(call manifest_name,04_infrastructure-components): $(KUSTOMIZE) $(infrastructure_components) ALWAYS | $(manifests_dir)
	$(KUSTOMIZE) build $(OPENSHIFT_DIR)/cluster-capi-configmap > $@

$(manifests_dir):
	mkdir -p $(OPENSHIFT_DIR)/$@

#$(KUSTOMIZE):
#	$(MAKE) -C . kustomize

.PHONY: merge-bot
merge-bot: full-vendoring generate generate-openshift ## Runs targets that help merge-bot to rebase downstream ORC.

.PHONY: full-vendoring ## Runs commands that complete vendoring tasks for downstream ORC.
full-vendoring:
	go mod tidy && go mod vendor

.PHONY: ALWAYS
ALWAYS:
