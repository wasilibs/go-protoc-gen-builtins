package runner

import (
	"context"
	"crypto/rand"
	"log"
	"os"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	"github.com/wasilibs/wazero-helpers/allocator"

	"github.com/wasilibs/go-protoc-gen-builtins/internal/wasm"
)

func Run(name string, wasmBin []byte) {
	ctx := context.Background()
	ctx = experimental.WithMemoryAllocator(ctx, allocator.NewNonMoving())

	rt := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCoreFeatures(api.CoreFeaturesV2|experimental.CoreFeaturesThreads))

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	if _, err := rt.InstantiateWithConfig(ctx, wasm.Memory, wazero.NewModuleConfig().WithName("env")); err != nil {
		log.Fatal(err)
	}

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

	_, err := rt.InstantiateWithConfig(ctx, wasmBin, cfg)
	if err != nil {
		if sErr, ok := err.(*sys.ExitError); ok { //nolint:errorlint
			os.Exit(int(sErr.ExitCode()))
		}
		log.Fatal(err)
	}
}
