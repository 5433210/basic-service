HOST_OS			:=$(shell uname -s)
HOST_ARCH		:=$(shell uname -m)

ROOT_DIR		:=$(abspath $(shell cd $(dir $(MAKEFILE_LIST)) && pwd -P))
OUTPUT_DIR 		:=$(ROOT_DIR)/_output
TEMP_DIR		:=$(ROOT_DIR)/_temp
TOOL_DIR 		:=$(OUTPUT_DIR)/tool

SHELL 			:=$(shell which bash)
VERSION			:=0.0.1
DATE			:=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
COVERAGE 		:=60
ROOT_MODULE		:=wailik.com

ifneq ($(findstring $(HOST_OS), "Darwin darwin"), )
	DEFAULT_PLATFORM_OS := darwin
endif

ifneq ($(findstring $(HOST_OS), "Linux linux"), )
	DEFAULT_PLATFORM_OS := linux
endif

ifneq ($(filter i%86, $(HOST_ARCH)),)
	DEFAULT_PLATFORM_ARCH := 386
endif

ifneq ($(filter x86_64 i%86_64 amd64, $(HOST_ARCH)),)
	DEFAULT_PLATFORM_ARCH := amd64
endif

ifneq ($(filter arm, $(HOST_ARCH)),)
	DEFAULT_PLATFORM_ARCH := arm
endif

ifneq ($(filter aarch64 arm64, $(HOST_ARCH)),)
	DEFAULT_PLATFORM_ARCH := arm64
endif

ifeq ($(DEFAULT_PLATFORM_OS),)
	ERR := $(error unsupported host os:$(HOST_OS))
endif

ifeq ($(DEFAULT_PLATFORM_ARCH),)
	ERR := $(error unsupported host arch:$(HOST_ARCH))
endif

ifeq ($(PLATFORMS),)
	PLATFORMS := $(DEFAULT_PLATFORM_OS)_$(DEFAULT_PLATFORM_ARCH)
endif

ifneq ($(IMG_ARCHS),)
	ps := $(foreach a, $(IMG_ARCHS), linux_$(a))
	PLATFORMS += $(ps)
endif

PLATFORMS := $(shell echo $(PLATFORMS) | xargs -n 1 | sort -u | xargs) #去掉重复项

COMMA 			:= ,
SPACE 			:=
SPACE 			+=

.PHONY: common-test
common-test:
	@echo $(ROOT_DIR)
	@echo $(SHELL)
	@echo $(DATE)
	@echo $(PLATFORMS)
	@echo $()
	@echo $(ROOT_MODULE)
	@echo $(COVERAGE)
	@echo $(VERSION)

