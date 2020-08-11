NAME=ledger
NAME_COMMAND_HANDLER=command-handler
VERSION=dev
OS ?= linux
PROJECT_PATH ?= github.com/stone-co/the-amazing-ledger
PKG ?= github.com/stone-co/the-amazing-ledger/cmd
REGISTRY ?= stone-co
TERM=xterm-256color
CLICOLOR_FORCE=true
RICHGO_FORCE_COLOR=1

.PHONY: test
test:
	@echo "==> Running Tests"
	go test -v ./...

.PHONY: compile
compile: clean
	@echo "==> Go Building CommandHandler"
	@env GOOS=${OS} GOARCH=amd64 go build -v -o build/${NAME_COMMAND_HANDLER} ${PKG}/${NAME_COMMAND_HANDLER}

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
	@echo "==> installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	go install ./...
	go test -i ./...
	$$(go env GOPATH)/bin/golangci-lint run -c ./.golangci.yml ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests"
	@richgo test -failfast -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: generate
generate:
	@echo "Go Generating"
	go get github.com/kevinburke/go-bindata/...
	go generate ./...
