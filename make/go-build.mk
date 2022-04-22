GO	:=go

GO_SUPPORTED_VERSIONS ?= 1.18

ifeq ($(shell if ! which $(GO) &>/dev/null; then echo no;fi), 'no')	
	ERR	:=$(error golang not install)
endif

GO_LDFLAGS +=
GO_BUILD_FLAGS += $(GO_LDFLAGS)

GOOS :=$(shell go env GOOS)
ifeq ($(GOOS),)
	ERR	:= $(error go env GOOS not set)
endif

ifeq ($(GOOS),windows)
	GO_SUFFIX_EXE := .exe
endif

GOPATH :=$(shell go env GOPATH)
ifeq ($(GOPATH),)
	ERR := $(error go env GOPATH not set)
endif

GOBIN :=$(GOPATH)/bin

GO_CMDS_DIR :=$(ROOT_DIR)/cmd

GO_CMDS ?= $(foreach cmd, $(wildcard $(GO_CMDS_DIR)/*), $(notdir $(cmd)))
ifeq ($(GO_CMDS),)
  ERR := $(error commands not found in $(GO_CMDS_DIR))
endif

# EXCLUDE_TESTS=wailik.com/pkg/errors wailik.com/pkg/log

.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif

.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH) $(GO_BUILD_FLAGS)"
	@mkdir -p $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_SUFFIX_EXE) $(GO_CMDS_DIR)/$(COMMAND)

.PHONY: go.build
go.build: go.build.verify $(foreach p, $(PLATFORMS), $(addprefix go.build., $(addprefix $(p)., $(GO_CMDS))))

.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: go.lint
go.lint: go.dependencies.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes $(PLATFORMS)"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

# .PHONY: go.test
# go.test: go.dependencies.verify.go-junit-report
# 	@echo "===========> Run unit test"
# 	set -o pipefail;$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
# 		-timeout=10m -short -v `go list ./...|\
# 		egrep -v $(subst $(SPACE),'|',$(sort $(EXCLUDE_TESTS)))` 2>&1 | \
# 		tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)
# # sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out # remove mock_.*.go files from test coverage
# 	$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

# .PHONY: go.test.cover
# go.test.cover: go.test
# 	$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
# 		awk -v target=$(COVERAGE) -f $(ROOT_DIR)/scripts/coverage.awk

.PHONY: go.update
go.update: go.dependencies.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct

.PHONY: go-build-test
go-build-test:
	@echo $(GO)
	@echo $(GO_INSTALLED)

.PHONY: go.tidy
go.tidy:
	@$(GO) mod tidy