package protoc

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuf(t *testing.T) {
	if err := os.RemoveAll(filepath.Join("build", "buf")); err != nil {
		t.Fatalf("failed to remove build directory: %v", err)
	}

	output := bytes.Buffer{}
	cmd := exec.Command("go", "run", "github.com/bufbuild/buf/cmd/buf@v1.28.1", "generate")
	cmd.Stderr = &output
	cmd.Stdout = &output
	cmd.Dir = "testdata"
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run buf: %v\n%v", err, output.String())
	}

	for _, path := range []string{
		filepath.Join("build", "buf", "cpp", "helloworld.pb.cc"),
		filepath.Join("build", "buf", "csharp", "Helloworld.cs"),
		filepath.Join("build", "buf", "java", "io", "grpc", "examples", "helloworld", "HelloReply.java"),
		filepath.Join("build", "buf", "kotlin", "io", "grpc", "examples", "helloworld", "HelloReplyKt.kt"),
		filepath.Join("build", "buf", "objc", "Helloworld.pbobjc.m"),
		filepath.Join("build", "buf", "php", "Helloworld", "HelloReply.php"),
		filepath.Join("build", "buf", "python", "helloworld_pb2.py"),
		filepath.Join("build", "buf", "python", "helloworld_pb2.pyi"),
		filepath.Join("build", "buf", "ruby", "helloworld_pb.rb"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("failed to stat %v: %v", path, err)
		}
	}
}

func TestProtoc(t *testing.T) {
	protosDir := filepath.Join("testdata", "protos")
	protosDirAbs, _ := filepath.Abs(protosDir)

	outDir := filepath.Join("build", "protoc", "python")
	outDirAbs, _ := filepath.Abs(outDir)

	// protoc requires directory created in advance
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "relative paths",
			args: []string{
				"-I" + protosDir, "--python_out=" + outDir, filepath.Join(protosDir, "helloworld.proto"),
			},
		},
		{
			name: "absolute paths",
			args: []string{
				"-I" + protosDirAbs, "--python_out=" + outDirAbs, filepath.Join(protosDirAbs, "helloworld.proto"),
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if err := os.RemoveAll(outDir); err != nil {
				t.Fatalf("failed to create directory: %v", err)
			}

			// protoc requires directory created in advance
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				t.Fatalf("failed to create directory: %v", err)
			}

			output := bytes.Buffer{}
			args := []string{"go", "run", "./cmd/protoc"}
			args = append(args, tc.args...)
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stderr = &output
			cmd.Stdout = &output
			if err := cmd.Run(); err != nil {
				t.Fatalf("failed to run protoc: %v\n%v", err, output.String())
			}

			path := filepath.Join(outDir, "helloworld_pb2.py")
			if _, err := os.Stat(path); err != nil {
				t.Errorf("failed to stat %v: %v", path, err)
			}
		})
	}
}
