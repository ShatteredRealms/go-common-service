//go:build tools
// +build tools

package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest"
	_ "github.com/mitranim/gow@latest"
	_ "github.com/onsi/ginkgo/v2/ginkgo@latest"
	_ "github.com/spf13/cobra-cli@latest"
	_ "go.uber.org/mock/mockgen@latest"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go@latest"
)
