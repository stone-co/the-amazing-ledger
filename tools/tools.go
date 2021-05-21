// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"                                   // buf.build
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking"         // gRPC breaking changes checker
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint"             // gRPC linter
	_ "github.com/golang/protobuf/protoc-gen-go"                          // protoc Go plugin
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"               // general linter
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway" // gRPC gateway
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"    // Generates swagger from gRPC definitions
	_ "github.com/kyoh86/richgo"                                          // enrich go test outputs with text decorations
	_ "github.com/resotto/gochk/cmd/gochk"                                // gochk architecture linter
	_ "golang.org/x/tools/cmd/goimports"                                  //Updates imports and formats code
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"                     // gRPC generator
)
