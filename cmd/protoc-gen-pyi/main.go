package main

import (
	"github.com/wasilibs/go-protoc/internal/runner"
	"github.com/wasilibs/go-protoc/internal/wasm"
)

func main() {
	runner.Run("protoc-gen-pyi", wasm.ProtocGenPyi)
}
