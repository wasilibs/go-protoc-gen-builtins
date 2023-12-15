package main

import (
	"context"
	"crypto/rand"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/wasilibs/go-protoc-gen-grpc/internal/wasix_32v1"
	"github.com/wasilibs/go-protoc-gen-grpc/internal/wasm"
	wazero "github.com/wasilibs/wazerox"
	"github.com/wasilibs/wazerox/api"
	"github.com/wasilibs/wazerox/experimental"
	"github.com/wasilibs/wazerox/experimental/sys"
	"github.com/wasilibs/wazerox/experimental/sysfs"
	"github.com/wasilibs/wazerox/imports/wasi_snapshot_preview1"
	wzsys "github.com/wasilibs/wazerox/sys"
)

func main() {
	ctx := context.Background()

	rt := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCoreFeatures(api.CoreFeaturesV2|experimental.CoreFeaturesThreads))

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)
	wasix_32v1.MustInstantiate(ctx, rt)

	args := []string{"protoc"}
	args = append(args, os.Args[1:]...)
	for _, arg := range args {
		println(arg)
	}

	fsCfg := wazero.NewFSConfig().(sysfs.FSConfig).WithSysFSMount(cmdFS{cwd: sysfs.DirFS("."), root: sysfs.DirFS("/")}, "/")
	fsCfg = fsCfg.(sysfs.FSConfig).WithRawPaths()

	cfg := wazero.NewModuleConfig().
		WithStartFunctions(). // Manually start after setting global state
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithStderr(os.Stderr).
		WithStdout(os.Stdout).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithArgs(args...).
		WithFSConfig(fsCfg)
	for _, env := range os.Environ() {
		k, v, _ := strings.Cut(env, "=")
		cfg = cfg.WithEnv(k, v)
	}

	mod, err := rt.InstantiateWithConfig(ctx, wasm.Protoc, cfg)
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	malloc := mod.ExportedFunction("malloc")
	res, err := malloc.Call(ctx, uint64(len(wd)+1))
	if err != nil {
		log.Fatal(err)
	}
	newWDPtr := uint32(res[0])
	buf, ok := mod.Memory().Read(newWDPtr, uint32(len(wd)+1))
	if !ok {
		log.Fatal("failed to read cwd allocation")
	}
	copy(buf, wd)
	buf[len(wd)] = 0

	main := mod.ExportedFunction("_start")
	_, err = main.Call(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

type cmdFS struct {
	cwd  sys.FS
	root sys.FS
}

func (fs cmdFS) OpenFile(path string, flag sys.Oflag, perm fs.FileMode) (sys.File, sys.Errno) {
	return fs.fs(path).OpenFile(path, flag, perm)
}

func (fs cmdFS) Lstat(path string) (wzsys.Stat_t, sys.Errno) {
	return fs.fs(path).Lstat(path)
}

func (fs cmdFS) Stat(path string) (wzsys.Stat_t, sys.Errno) {
	return fs.fs(path).Stat(path)
}

func (fs cmdFS) Mkdir(path string, perm fs.FileMode) sys.Errno {
	return fs.fs(path).Mkdir(path, perm)
}

func (fs cmdFS) Chmod(path string, perm fs.FileMode) sys.Errno {
	return fs.fs(path).Chmod(path, perm)
}

func (fs cmdFS) Rename(from string, to string) sys.Errno {
	return fs.fs(from).Rename(from, to)
}

func (fs cmdFS) Rmdir(path string) sys.Errno {
	return fs.fs(path).Rmdir(path)
}
func (fs cmdFS) Unlink(path string) sys.Errno {
	return fs.fs(path).Unlink(path)
}

func (fs cmdFS) Link(oldPath string, newPath string) sys.Errno {
	return fs.fs(oldPath).Link(oldPath, newPath)
}

func (fs cmdFS) Symlink(oldPath string, linkName string) sys.Errno {
	return fs.fs(oldPath).Symlink(oldPath, linkName)
}

func (fs cmdFS) Readlink(path string) (string, sys.Errno) {
	return fs.fs(path).Readlink(path)
}

func (fs cmdFS) Utimens(path string, atim int64, mtim int64) sys.Errno {
	return fs.fs(path).Utimens(path, atim, mtim)
}

func (fs cmdFS) fs(path string) sys.FS {
	if len(path) > 0 && path[0] != '/' {
		return fs.cwd
	}
	return fs.root
}
