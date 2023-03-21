ROOT_MODULE=wailik.com

include scripts/make/common.mk
include scripts/make/go-depend.mk
include scripts/make/go-build.mk
include scripts/make/image.mk

define USAGE_OPTIONS
Options:
  PLATFORMS			The multiple platforms to build. Default is the host platform.
					ONLY support os darwin/linux/windows, arch 386/amd64/arm/arm64.
        			Example: make PLATFORMS="linux_amd64 linux_arm64"
  IMG_ARCHS		 	The image architectures. Default is amd64.
  					ONLY support arch 386/amd64/arm/arm64.
					Example: make IMG_ARCHS ="amd64 arm64"
  VERSION      		The version information compiled into binaries.
        			The default is obtained from gsemver or git.
  V					Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

.DEFAULT_GOAL := all
.PHONY: all
all: dependencies build 

.PHONY: build
build:
	@$(MAKE) go.build

.PHONY: image
image:
	@echo "===========> Create images"
	@$(MAKE) image.build

.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@$(MAKE) go.clean

.PHONY: lint
lint:
	@echo "===========> lint"
	@$(MAKE) go.lint

.PHONY: dependencies
dependencies:
	@echo "===========> install dependencies"
	@$(MAKE) go.dependencies.install
	@$(MAKE) go.tidy

.PHONY: update
update:
	@echo "===========> update environment"
	@$(MAKE) go.update

.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"