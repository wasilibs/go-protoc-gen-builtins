package main

import (
	"github.com/wasilibs/go-protoc-gen-builtins/internal/runner"
	"github.com/wasilibs/go-protoc-gen-builtins/internal/wasm"
)

func main() {
	runner.Run("protoc-gen-upb_minitable", wasm.ProtocGenUPBMinitable)
}
