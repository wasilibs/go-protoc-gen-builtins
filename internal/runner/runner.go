package runner

import (
	"crypto/rand"
	"log"
	"os"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	wzsys "github.com/tetratelabs/wazero/sys"

	"github.com/wasilibs/go-protoc-gen-builtins/internal/wasix_32v1"
	"github.com/wasilibs/wazero-helpers/allocator"
)

func Run(name string, wasm []byte) {
	ctx := wasix_32v1.BackgroundContext()
	ctx = experimental.WithMemoryAllocator(ctx, allocator.NewNonMoving())

	rt := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCoreFeatures(api.CoreFeaturesV2|experimental.CoreFeaturesThreads))

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)
	wasix_32v1.MustInstantiate(ctx, rt)

	args := []string{name}
	args = append(args, os.Args[1:]...)

	cfg := wazero.NewModuleConfig().
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithStderr(os.Stderr).
		WithStdout(os.Stdout).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithArgs(args...)
	for _, env := range os.Environ() {
		k, v, _ := strings.Cut(env, "=")
		cfg = cfg.WithEnv(k, v)
	}

	_, err := rt.InstantiateWithConfig(ctx, wasm, cfg)
	if err != nil {
		if sErr, ok := err.(*wzsys.ExitError); ok {
			os.Exit(int(sErr.ExitCode()))
		}
		log.Fatal(err)
	}
}
