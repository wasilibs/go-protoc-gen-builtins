# go-protoc-gen-builtins

go-protoc-gen-builtins is a distribution of the code generator plugins from protoc, the 
[protocol buffers][1] compiler, that can be built with Go. It does not actually reimplement any
functionality of protoc in Go, instead compiling the original source code to WebAssembly, and 
executing with the pure Go Wasm runtime [wazero][2]. This means that `go install` or `go run`
can be used to execute it, with no need to rely on external package managers such as Homebrew,
on any platform that Go supports.

This project is primarily targeted at [Buf][3] users, or users of other alternative protocol buffer
compilers. `protoc` users should have no need for the plugins in this repository because they are
already built-in to `protoc`.

## Installation

Precompiled binaries are available in the [releases](https://github.com/wasilibs/go-protoc-gen-builtins/releases).
Alternatively, install the plugin you want using `go install`.

```bash
$ go install github.com/wasilibs/go-protoc-gen-builtins/cmd/protoc-gen-python@latest
```

As long as `$GOPATH/bin`, e.g. `~/go/bin` is on the `PATH`, `buf` should find it automatically.

```yaml
version: v1
plugins:
  - plugin: python
    out: out/python
```

To avoid installation entirely, it can be convenient to use `go run` with `path` instead.

```yaml
version: v1
plugins:
  - plugin: python
    out: out/python
    path: ["go", "run", "github.com/wasilibs/go-protoc-gen-builtins/cmd/protoc-gen-python@latest"]
```

If invoking `buf` itself with `go run`, it is possible to have full protobuf generation with no
installation of tools, besides Go itself, on any platform that Go supports. The above examples use
`@latest`, but it is recommended to specify a version, in which case all of the developers on your
codebase will use the same version of the tool with no special steps.

For gRPC plugins, also see [go-protoc-gen-grpc][4].

A full example is available at [example](./example/). To generate protos, enter the directory and run
`go run github.com/bufbuild/buf/cmd/buf@v1.30.0 generate`. As long as your machine has Go installed,
you will be able to generate protos. The first time using `go run` for a command, Go automatically builds
it making it slower, but subsequent invocations should be quite fast.

[1]: https://protobuf.dev/
[2]: https://wazero.io/
[3]: https://buf.build/
[4]: https://github.com/wasilibs/go-protoc-gen-builtins-gen-grpc
