package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Run with: cd tests/e2e && go test -v

var binaryPath string

func TestMain(m *testing.M) {
	// Get absolute path for output binary
	absPath, err := filepath.Abs("j_test_bin")
	if err != nil {
		panic("failed to get absolute path: " + err.Error())
	}
	binaryPath = absPath

	// Build from src directory
	srcDir, err := filepath.Abs("../../src")
	if err != nil {
		panic("failed to get src path: " + err.Error())
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/j")
	cmd.Dir = srcDir
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("failed to build binary: " + err.Error() + "\n" + string(output))
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove(binaryPath)

	os.Exit(code)
}

func TestStatusCommand(t *testing.T) {
	cmd := exec.Command(binaryPath, "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("status command failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "j status") {
		t.Error("expected 'j status' header in output")
	}
}

func TestInstallCommand(t *testing.T) {
	cmd := exec.Command(binaryPath, "install", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("install --help failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "install") {
		t.Error("expected 'install' in help output")
	}
}

func TestCleanCommand(t *testing.T) {
	cmd := exec.Command(binaryPath, "clean", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("clean --help failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "clean") {
		t.Error("expected 'clean' in help output")
	}
}

func TestUpdateCommand(t *testing.T) {
	cmd := exec.Command(binaryPath, "update", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("update --help failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "update") {
		t.Error("expected 'update' in help output")
	}
}

func TestRunGitCommands(t *testing.T) {
	cmd := exec.Command(binaryPath, "run", "git", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run git --help failed: %v\n%s", err, output)
	}
	out := string(output)
	for _, sub := range []string{"feat", "fix", "chore", "push", "sync"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected '%s' in git help output", sub)
		}
	}
}

func TestRunDockerCommands(t *testing.T) {
	cmd := exec.Command(binaryPath, "run", "docker", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run docker --help failed: %v\n%s", err, output)
	}
	out := string(output)
	for _, sub := range []string{"rm", "rmi", "clean", "reset"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected '%s' in docker help output", sub)
		}
	}
}
