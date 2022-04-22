GO_DEPENDENCIES ?= golines go-junit-report golangci-lint goimports mockgen gotests git-chglog github-release go-mod-outdated protoc-gen-go

.PHONY: go.dependencies.install
go.dependencies.install: $(addprefix go.dependencies.verify., $(GO_DEPENDENCIES))

.PHONY: go.dependencies.install.%
go.dependencies.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) go.install.$*

go.dependencies.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) go.dependencies.install.$*; fi

.PHONY: go.install.golines
go.install.golines:
	@$(GO) install github.com/segmentio/golines@latest

.PHONY: go.install.go-junit-report
go.install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest

.PHONY: go.install.golangci-lint
go.install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
	@golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi

.PHONY: go.install.goimports
go.install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: go.install.mockgen
go.install.mockgen:
	@$(GO) install github.com/golang/mock/mockgen@latest

.PHONY: go.install.gotests
go.install.gotests:
	@$(GO) install github.com/cweill/gotests/...@latest

.PHONY: go.install.git-chglog
go.install.git-chglog:
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

.PHONY: go.install.github-release
go.install.github-release:
	@$(GO) install github.com/github-release/github-release@latest

.PHONY: go.install.go-mod-outdated
go.install.go-mod-outdated:
	@$(GO) install github.com/psampaz/go-mod-outdated@latest

.PHONY: go.install.protoc-gen-go
go.install.protoc-gen-go:
	@$(GO) install github.com/golang/protobuf/protoc-gen-go@latest