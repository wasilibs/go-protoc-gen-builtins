package main

import (
	"github.com/curioswitch/go-build"
	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/boot"
	"github.com/goyek/x/cmd"
	"github.com/wasilibs/tools/tasks"
)

func main() {
	tasks.Define(tasks.Params{
		LibraryName: "protoc",
		LibraryRepo: "protocolbuffers/protobuf",
		GoReleaser:  true,
	})
	runBuf := "go run github.com/bufbuild/buf/cmd/buf@" + verBuf
	build.RegisterCommandDownloads(runBuf + " --version")
	goyek.Define(goyek.Task{
		Name: "example",
		Action: func(a *goyek.A) {
			cmd.Exec(a, runBuf+" generate", cmd.Dir("example"))
		},
	})
	boot.Main()
}
