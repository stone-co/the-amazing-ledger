NAME=ledger
NAME_COMMAND_HANDLER=server
VERSION=dev
OS ?= linux
PROJECT_PATH ?= github.com/stone-co/the-amazing-ledger
PKG ?= github.com/stone-co/the-amazing-ledger/cmd
REGISTRY ?= stoneopenbank
TERM=xterm-256color
CLICOLOR_FORCE=true
RICHGO_FORCE_COLOR=1
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_BUILD_TIME=$(shell date '+%Y-%m-%d__%I:%M:%S%p')

BUF_VERSION=v0.43.2

.PHONY: setup
setup:
	@echo "==> Setup: Tidying modules"
	go mod tidy
	@echo "==> Setup: Getting dependencies"
	go mod download
	@echo "==> Setup: Getting tools"
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %


.PHONY: test
test:
	@echo "==> Test: Running Tests"
	gotest -v ./...

.PHONY: compile
compile: clean
	@echo "==> Go Building CommandHandler"
	@env GOOS=${OS} GOARCH=amd64 go build -v -o build/${NAME_COMMAND_HANDLER} \
		-ldflags "-X main.BuildGitCommit=$(GIT_COMMIT) -X main.BuildTime=$(GIT_BUILD_TIME)" \
		${PKG}/${NAME_COMMAND_HANDLER} 

.PHONY: build
build: compile
	@echo "==> Building Docker CommandHandler image"
	@docker build -t ${REGISTRY}/${NAME}-${NAME_COMMAND_HANDLER}:${VERSION} build -f build/Dockerfile-command_handler

.PHONY: push
push:
	@echo "==>Push Docker CommandHandler image"
	@docker push ${REGISTRY}/${NAME}-${NAME_COMMAND_HANDLER}:${VERSION}

.PHONY: clean
clean:
	@echo "==> Cleaning releases"
	@GOOS=${OS} go clean -i -x ./...
	@rm -f build/${NAME_COMMAND_HANDLER}

.PHONY: metalint
metalint:
	@echo "==> Running Linters"
	go test -i ./...
	golangci-lint run -c ./.golangci.yml ./...
	buf check lint

.PHONY: archlint
archlint:
	@echo "==> Running architecture linter(gochk)"
	gochk -c ./gochk-arch-linter.json

.PHONY: test-coverage
test-coverage:
	@echo "Running tests"
	@richgo test -failfast -coverprofile=coverage.out ./...
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
