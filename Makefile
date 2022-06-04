# binary aliases
GO  = GO111MODULE=on go

# variables
BUILD_DIR=./build

##@ Dependencies
.PHONY: install install-tools intsall-go

install: install-tools ## Install all dependencies

install-tools: ## Install runtimes
	@echo "Installing tooling"
	asdf install

##@ Build
.PHONY: run build

run: ## Run locally
	$(GO) run *.go

build: ## Build binaries
	@echo "Buiding for linux"
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/steam-screenshot-sorter_linux_amd64
	@echo "Buiding for windows"
	GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/steam-screenshot-sorter_windows_amd64.exe

##@ Helpers
.PHONY: help

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)