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

.PHONY: setup
setup:
	@echo "==> Setup: Getting tools"
	go mod tidy
	GO111MODULE=on go install \
	github.com/bufbuild/buf/cmd/buf \
	github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking \
	github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	github.com/kevinburke/go-bindata \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	github.com/kyoh86/richgo \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
  golang.org/x/tools/cmd/goimports \
  github.com/kyoh86/richgo \
	github.com/resotto/gochk/cmd/gochk


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
	@buf generate --file ./proto/ledger/ledger.proto
	@go generate ./...

.PHONY: goimports
goimports:
	@echo "Go imports"
	@goimports -w $(shell \
                  	find . -not \( \( -name .git -o -name .go -o -name vendor \) -prune \) \
                  	-name '*.go')

.PHONY: pre/push
pre/push: goimports metalint archlint test
