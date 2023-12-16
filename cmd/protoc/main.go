package main

import (
	"context"
	"crypto/rand"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/wasilibs/go-protoc/internal/wasix_32v1"
	"github.com/wasilibs/go-protoc/internal/wasm"
	wazero "github.com/wasilibs/wazerox"
	"github.com/wasilibs/wazerox/api"
	"github.com/wasilibs/wazerox/experimental"
	"github.com/wasilibs/wazerox/experimental/sys"
	"github.com/wasilibs/wazerox/experimental/sysfs"
	"github.com/wasilibs/wazerox/imports/wasi_snapshot_preview1"
	wazerosys "github.com/wasilibs/wazerox/sys"
	wzsys "github.com/wasilibs/wazerox/sys"
)

func main() {
	ctx := context.Background()

	rt := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCoreFeatures(api.CoreFeaturesV2|experimental.CoreFeaturesThreads))

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)
	wasix_32v1.MustInstantiate(ctx, rt)

	args := []string{"protoc"}
	args = append(args, os.Args[1:]...)

	fsCfg := wazero.NewFSConfig().(sysfs.FSConfig).WithSysFSMount(cmdFS{cwd: sysfs.DirFS("."), root: sysfs.DirFS("/")}, "/")
	fsCfg = fsCfg.(sysfs.FSConfig).WithRawPaths()

	cfg := wazero.NewModuleConfig().
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

	_, err := rt.InstantiateWithConfig(ctx, wasm.Protoc, cfg)
	if err != nil {
		if sErr, ok := err.(*wazerosys.ExitError); ok {
			os.Exit(int(sErr.ExitCode()))
		}
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
