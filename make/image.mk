DOCKER 			:= docker
REGISTRY_PREFIX ?= 192.168.1.200:80
PROJECT_PREFIX	:= basic-service
BASE_IMAGE 		:= alpine:3.15

ifeq ($(shell if ! which $(DOCKER) &>/dev/null; then echo no;fi), 'no')	
	ERR	:=$(error docker has not been installed)
endif

ifeq ($(shell ps -ef | grep "docker" | grep "server" > /dev/null; echo $?), 1)
	ERR :=$(error docker is not running)
endif

# Determine image files by looking into build/docker/*/Dockerfile
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)
# Determine images names by stripping out the dir names
IMAGES ?= $(foreach image,${IMAGES_DIR},$(notdir ${image}))

ifeq (${IMAGES},)
  ERR := $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.build
image.build: $(foreach a,$(IMG_ARCHS),$(addprefix image.build., $(addprefix linux_$(a)., $(IMAGES))))

.PHONY: image.build.%
image.build.%: go.build.%
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	mkdir -p $(TEMP_DIR)/$(IMAGE)
	cat $(ROOT_DIR)/build/docker/$(IMAGE)/Dockerfile\
		| sed "s#BASE_IMAGE#$(BASE_IMAGE)#g" >$(TEMP_DIR)/$(IMAGE)/Dockerfile
	cp $(OUTPUT_DIR)/platforms/$(IMAGE_PLAT)/$(IMAGE) $(TEMP_DIR)/$(IMAGE)/
	DST_DIR=$(TEMP_DIR)/$(IMAGE) $(ROOT_DIR)/build/docker/$(IMAGE)/build.sh 2>/dev/null || true
	$(eval BUILD_SUFFIX := --pull -t $(REGISTRY_PREFIX)/$(PROJECT_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION) $(TEMP_DIR)/$(IMAGE))
	@echo $(BUILD_SUFFIX)
	$(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)
	@rm -rf $(TEMP_DIR)/$(IMAGE)

.PHONY: image.push
image.push: $(foreach a,$(IMG_ARCHS),$(addprefix image.push., $(addprefix linux_$(a)., $(IMAGES))))

.PHONY: image.push.%
image.push.%: image.build.%
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Pushing image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) push $(REGISTRY_PREFIX)/$(PROJECT_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION)

.PHONY: image.pull
image.pull: $(foreach a,$(IMG_ARCHS),$(addprefix image.pull., $(addprefix linux_$(a)., $(IMAGES))))

.PHONY: image.pull.%
image.pull.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> pulling image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) pull $(REGISTRY_PREFIX)/$(PROJECT_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION)

.PHONY: image.clean
image.clean: $(foreach a,$(IMG_ARCHS),$(addprefix image.clean., $(addprefix linux_$(a)., $(IMAGES))))

.PHONY: image.clean.%
image.clean.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> cleaning image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) rmi $(REGISTRY_PREFIX)/$(PROJECT_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION)	