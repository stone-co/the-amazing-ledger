// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"                           // buf.build
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking" // gRPC breaking changes checker
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint"     // gRPC linter
	_ "github.com/golang/protobuf/protoc-gen-go"                  // protoc Go plugin
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"       // general linter
	_ "github.com/kevinburke/go-bindata"                          // converts any file into manageable Go source code
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"             // gRPC generator
)
