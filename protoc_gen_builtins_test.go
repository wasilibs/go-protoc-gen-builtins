package protoc_gen_builtins

import (
	"bytes"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

//go:embed testdata/buf.gen.gorun.yaml
var bufGenGorunYaml []byte

//go:embed testdata/buf.gen.installed.yaml
var bufGenInstalledYaml []byte

func TestBuf(t *testing.T) {
	goExe := filepath.Join(runtime.GOROOT(), "bin", "go")
	if err := os.RemoveAll(filepath.Join("build", "buf")); err != nil {
		t.Fatalf("failed to remove build directory: %v", err)
	}

	pluginsDir := filepath.Join("build", "plugins")
	if err := os.RemoveAll(pluginsDir); err != nil {
		t.Fatalf("failed to remove build directory: %v", err)
	}

	plugins := []string{"cpp", "csharp", "java", "kotlin", "objc", "php", "pyi", "python", "ruby", "rust"}
	for _, plugin := range plugins {
		output := bytes.Buffer{}
		cmd := exec.Command(goExe, "build", "-o", filepath.Join(pluginsDir, "protoc-gen-"+plugin), "./cmd/protoc-gen-"+plugin)
		cmd.Stderr = &output
		cmd.Stdout = &output
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to build plugin %v: %v\n%v", plugin, err, output.String())
		}
	}

	tests := []struct {
		name       string
		bufGenYaml []byte
	}{
		{
			name:       "gorun",
			bufGenYaml: bufGenGorunYaml,
		},
		{
			name:       "installed",
			bufGenYaml: bufGenInstalledYaml,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// We can only have one buf.gen.yaml at a time since buf provides no way of
			// an alternate config file for generation. This also means this should never
			// be marked parallel.

			bufGenPath := filepath.Join("testdata", "buf.gen.yaml")
			if err := os.WriteFile(bufGenPath, tc.bufGenYaml, 0o644); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(bufGenPath)

			output := bytes.Buffer{}
			env := os.Environ()
			pluginsDirAbs, _ := filepath.Abs(pluginsDir)
			for i, val := range env {
				if strings.HasPrefix(val, "PATH=") {
					pathVal := pluginsDirAbs + string(os.PathListSeparator) + filepath.Join(runtime.GOROOT(), "bin")
					env[i] = "PATH=" + pathVal
				}
				println(env[i])
			}
			cmd := exec.Command(goExe, "run", "github.com/bufbuild/buf/cmd/buf@v1.28.1", "generate")
			cmd.Stderr = &output
			cmd.Stdout = &output
			cmd.Env = env
			cmd.Dir = "testdata"
			if err := cmd.Run(); err != nil {
				t.Fatalf("failed to run buf: %v\n%v", err, output.String())
			}

			outDir := filepath.Join("build", "buf", tc.name)
			for _, path := range []string{
				filepath.Join(outDir, "cpp", "helloworld.pb.cc"),
				filepath.Join(outDir, "csharp", "Helloworld.cs"),
				filepath.Join(outDir, "java", "io", "grpc", "examples", "helloworld", "HelloReply.java"),
				filepath.Join(outDir, "kotlin", "io", "grpc", "examples", "helloworld", "HelloReplyKt.kt"),
				filepath.Join(outDir, "objc", "Helloworld.pbobjc.m"),
				filepath.Join(outDir, "php", "Helloworld", "HelloReply.php"),
				filepath.Join(outDir, "python", "helloworld_pb2.py"),
				filepath.Join(outDir, "python", "helloworld_pb2.pyi"),
				filepath.Join(outDir, "ruby", "helloworld_pb.rb"),
			} {
				if _, err := os.Stat(path); err != nil {
					t.Errorf("failed to stat %v: %v", path, err)
				}
			}
		})
	}
}
