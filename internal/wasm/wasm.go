package wasm

import _ "embed"

//go:embed protoc-gen-cpp.wasm
var ProtocGenCpp []byte

//go:embed protoc-gen-csharp.wasm
var ProtocGenCsharp []byte

//go:embed protoc-gen-java.wasm
var ProtocGenJava []byte

//go:embed protoc-gen-kotlin.wasm
var ProtocGenKotlin []byte

//go:embed protoc-gen-objc.wasm
var ProtocGenObjc []byte

//go:embed protoc-gen-php.wasm
var ProtocGenPhp []byte

//go:embed protoc-gen-pyi.wasm
var ProtocGenPyi []byte

//go:embed protoc-gen-python.wasm
var ProtocGenPython []byte

//go:embed protoc-gen-ruby.wasm
var ProtocGenRuby []byte

//go:embed protoc-gen-rust.wasm
var ProtocGenRust []byte
