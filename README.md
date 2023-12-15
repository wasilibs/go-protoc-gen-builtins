# go-protoc

go-protoc is a cut-down distribution of protoc, the [protocol buffers][1] compiler that can be built
with Go. It does not actually reimplement any functionality of protoc in Go, instead compiling the
original source code to WebAssembly, and executing with the pure Go Wasm runtime [wazero][2].
This means that `go install` or `go run` can be used to execute it, with no need to rely on external
package managers such as Homebrew.

Note that currently, go-protoc __DOES NOT__ support executing protoc plugins, such as gRPC plugins.
[Buf][3] is a pure Go tool that can execute protoc plugins. It is intended to use go-protoc together
with buf to generate code for languages built-in to protoc. View the Buf [documentation][4] on setting
protoc_path and the languages that are supported.

A full example of usage can be found as an [integration test][5].

[1]: https://protobuf.dev/
[2]: https://wazero.io/
[3]: https://buf.build/
[4]: https://buf.build/docs/configuration/v1/buf-gen-yaml#protoc_path
[5]: ./testdata/
