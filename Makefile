NAME=ledger
NAME_COMMAND_HANDLER=server
VERSION=dev
OS ?= linux
PROJECT_PATH ?= github.com/stone-co/the-amazing-ledger
PKG ?= github.com/stone-co/the-amazing-ledger/cmd
REGISTRY ?= stoneopenbank
TERM=xterm-256color
CLICOLOR_FORCE=true
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_BUILD_TIME=$(shell date '+%Y-%m-%d__%I:%M:%S%p')

BUF_VERSION=0.43.2
GOLANGCI_LINT_VERSION=1.41.1

.PHONY: setup
setup:
	@echo "==> Setup: Tidying modules"
	go mod tidy
	@echo "==> Setup: Getting dependencies"
	go mod download
	@echo "==> Setup: Getting tools"
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
	@echo "==> Setup: Checking buf installation"
ifeq (, $(shell command -v buf 2> /dev/null))
	@echo "    ==> Setup: Buf not installed, please follow the instructions on https://docs.buf.build/installation and install version ${BUF_VERSION}"
else
	@echo "    ==> Setup: Buf already installed"
endif
	@echo "==> Setup: Checking golangci-lint installation"
ifeq (, $(shell command -v golangci-lint 2> /dev/null))
	@echo "    ==> Setup: Installing golangci-lint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_LINT_VERSION}
else
	@echo "    ==> Setup: Golangci-lint already installed"
endif

.PHONY: test
test:
	@echo "==> Test: Running Tests"
	gotest -v ./...

.PHONY: compile
compile: clean
	@echo "==> Compile: Go Building CommandHandler"
	@env GOOS=${OS} GOARCH=amd64 go build -v -o build/${NAME_COMMAND_HANDLER} \
		-ldflags "-X main.BuildGitCommit=$(GIT_COMMIT) -X main.BuildTime=$(GIT_BUILD_TIME)" \
		${PKG}/${NAME_COMMAND_HANDLER} 

.PHONY: build
build: compile
	@echo "==> Build: Building Docker CommandHandler image"
	@docker build -t ${REGISTRY}/${NAME}-${NAME_COMMAND_HANDLER}:${VERSION} build -f build/Dockerfile-command_handler

.PHONY: push
push:
	@echo "==> Push: Push Docker CommandHandler image"
	@docker push ${REGISTRY}/${NAME}-${NAME_COMMAND_HANDLER}:${VERSION}

.PHONY: clean
clean:
	@echo "==> Clean: Cleaning releases"
	@GOOS=${OS} go clean -i -x ./...
	@rm -f build/${NAME_COMMAND_HANDLER}

.PHONY: lint
lint: metalint protolint

.PHONY: metalint
metalint:
	@echo "==> Metalint: Running Linters"
	golangci-lint run -c ./.golangci.yml ./...

.PHONY: archlint
archlint:
	@echo "==> Running architecture linter(gochk)"
	gochk -c ./gochk-arch-linter.json

.PHONY: protolint
protolint:
	@echo "==> Proto lint: Running protofile linters"
	@buf lint
	@buf breaking --against '.git#branch=main'

.PHONY: test-coverage
test-coverage:
	@echo "Running tests"
	@gotest -failfast -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: generate
generate:
	@echo "Go Generating"
	@rm -rf gen/*
	@buf generate --path proto/ledger
	@go generate ./...

.PHONY: goimports
goimports:
	@echo "Go imports"
	@goimports -w $(shell \
                  	find . -not \( \( -name .git -o -name .go -o -name vendor \) -prune \) \
                  	-name '*.go')

.PHONY: pre/push
pre/push: goimports metalint archlint test

.PHONY: update-buf-dependencies
update-buf-dependencies:
	@buf beta mod update
