// +build tools

package tools

import (
	_ "github.com/kyoh86/richgo"           // enrich go test outputs with text decorations
	_ "github.com/resotto/gochk/cmd/gochk" // gochk architecture linter
	_ "golang.org/x/tools/cmd/goimports"   // Updates imports and formats code

	// GRPC Gateway - https://github.com/grpc-ecosystem/grpc-gateway#usage
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
