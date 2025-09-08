package protocgenbuiltins

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
	if err := os.RemoveAll(filepath.Join("out", "buf")); err != nil {
		t.Fatalf("failed to remove out directory: %v", err)
	}

	pluginsDir := filepath.Join("out", "plugins")
	if err := os.RemoveAll(pluginsDir); err != nil {
		t.Fatalf("failed to remove out directory: %v", err)
	}

	plugins := []string{"cpp", "csharp", "java", "kotlin", "objc", "php", "pyi", "python", "ruby", "rust", "upb", "upbdefs", "upb_minitable"}
	for _, plugin := range plugins {
		output := bytes.Buffer{}
		cmd := exec.Command("go", "build", "-o", filepath.Join(pluginsDir, "protoc-gen-"+plugin), "./cmd/protoc-gen-"+plugin)
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
		t.Run(tc.name, func(t *testing.T) {
			// We can only have one buf.gen.yaml at a time since buf provides no way of
			// an alternate config file for generation. This also means this should never
			// be marked parallel.

			if tc.name == "installed" && runtime.GOOS == "windows" {
				// Currently this is not working on Windows, will need a real machine to
				// debug. Since gorun works and installed works on other OS's, it seems
				// likely just an environment issue or an issue with buf on windows.
				t.Skip("skipping on windows")
			}

			bufGenPath := filepath.Join("testdata", "buf.gen.yaml")
			if err := os.WriteFile(bufGenPath, tc.bufGenYaml, 0o644); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(bufGenPath)

			output := bytes.Buffer{}
			env := os.Environ()
			pluginsDirAbs, _ := filepath.Abs(pluginsDir)
			for i, val := range env {
				k, v, _ := strings.Cut(val, "=")
				if k == "PATH" {
					pathVal := pluginsDirAbs + string(os.PathListSeparator) + v
					env[i] = "PATH=" + pathVal
				}
			}
			cmd := exec.Command("go", "run", "github.com/bufbuild/buf/cmd/buf@"+verBuf, "generate")
			cmd.Stderr = &output
			cmd.Stdout = &output
			cmd.Env = env
			cmd.Dir = "testdata"
			if err := cmd.Run(); err != nil {
				t.Fatalf("failed to run buf: %v\n%v", err, output.String())
			}

			outDir := filepath.Join("out", "buf", tc.name)
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
				filepath.Join(outDir, "upb", "helloworld.upb_minitable.c"),
				filepath.Join(outDir, "upb", "helloworld.upb.c"),
				filepath.Join(outDir, "upb", "helloworld.upbdefs.c"),
			} {
				if _, err := os.Stat(path); err != nil {
					t.Errorf("failed to stat %v: %v", path, err)
				}
			}
		})
	}
}
