// +build tools

package tools

import (
	// Architecture linter
	_ "github.com/resotto/gochk/cmd/gochk"

	// Import linter
	_ "golang.org/x/tools/cmd/goimports"

	// Testing helper
	_ "github.com/rakyll/gotest"

	// GRPC Gateway - https://github.com/grpc-ecosystem/grpc-gateway#usage
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
