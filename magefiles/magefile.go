package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/google/go-github/v58/github"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/wasilibs/magefiles" // mage:import
)

func init() {
	magefiles.SetLibraryName("protoc")
}

func Snapshot() error {
	return sh.RunV("go", "run", fmt.Sprintf("github.com/goreleaser/goreleaser@%s", verGoReleaser), "release", "--snapshot", "--clean")
}

func Release() error {
	return sh.RunV("go", "run", fmt.Sprintf("github.com/goreleaser/goreleaser@%s", verGoReleaser), "release", "--clean")
}

func UpdateUpstream() error {
	currBytes, err := os.ReadFile(filepath.Join("buildtools", "wasm", "version.txt"))
	if err != nil {
		return err
	}
	curr := strings.TrimSpace(string(currBytes))

	gh, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	var releases []*github.RepositoryRelease
	if err := gh.Get("repos/protocolbuffers/protobuf/releases?per_page=1", &releases); err != nil {
		return err
	}

	if len(releases) == 0 {
		return errors.New("could not find releases")
	}

	latest := releases[0].GetTagName()
	if latest == curr {
		fmt.Println("up to date")
		return nil
	}

	fmt.Println("updating to", latest)
	if err := os.WriteFile(filepath.Join("buildtools", "wasm", "version.txt"), []byte(latest), 0o644); err != nil {
		return err
	}

	mg.Deps(magefiles.UpdateLibs)

	return nil
}
