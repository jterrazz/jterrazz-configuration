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

func TestUpgradeCommand(t *testing.T) {
	cmd := exec.Command(binaryPath, "upgrade", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("upgrade --help failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "upgrade") {
		t.Error("expected 'upgrade' in help output")
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

// =============================================================================
// Sync command tests
// =============================================================================

func TestSyncHelp(t *testing.T) {
	cmd := exec.Command(binaryPath, "sync", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sync --help failed: %v\n%s", err, output)
	}
	out := string(output)
	for _, sub := range []string{"init", "status", "diff"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected '%s' subcommand in sync help output", sub)
		}
	}
	if !strings.Contains(out, "--all") {
		t.Error("expected '--all' flag in sync help output")
	}
}

func TestSyncStatusUnlinked(t *testing.T) {
	// Run in a temp directory with no .copier-answers.yml
	tmpDir := t.TempDir()

	cmd := exec.Command(binaryPath, "sync", "status")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sync status failed: %v\n%s", err, output)
	}
	out := string(output)
	if !strings.Contains(out, "Not linked") {
		t.Error("expected 'Not linked' for directory without .copier-answers.yml")
	}
	if !strings.Contains(out, "j sync init") {
		t.Error("expected hint to run 'j sync init'")
	}
}

func TestSyncStatusLinked(t *testing.T) {
	// Create a temp directory with a fake .copier-answers.yml
	tmpDir := t.TempDir()
	answersContent := `_commit: abc123
_src_path: /some/template
project_name: test-project
language: go
license: MIT
`
	err := os.WriteFile(filepath.Join(tmpDir, ".copier-answers.yml"), []byte(answersContent), 0644)
	if err != nil {
		t.Fatalf("failed to create answers file: %v", err)
	}

	cmd := exec.Command(binaryPath, "sync", "status")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sync status failed: %v\n%s", err, output)
	}
	out := string(output)
	if !strings.Contains(out, "Linked") {
		t.Error("expected 'Linked' for directory with .copier-answers.yml")
	}
	if !strings.Contains(out, "project_name") {
		t.Error("expected 'project_name' in status output")
	}
	if !strings.Contains(out, "test-project") {
		t.Error("expected 'test-project' value in status output")
	}
	if !strings.Contains(out, "language") {
		t.Error("expected 'language' in status output")
	}
}

func TestSyncUpdateUnlinked(t *testing.T) {
	// Running sync (update) in a directory without .copier-answers.yml should warn
	tmpDir := t.TempDir()

	cmd := exec.Command(binaryPath, "sync")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	// Should not fail hard, just warn
	_ = err
	out := string(output)
	if !strings.Contains(out, "No .copier-answers.yml") || !strings.Contains(out, "j sync init") {
		t.Error("expected warning about missing .copier-answers.yml with hint to run init")
	}
}

func TestSyncDiffUnlinked(t *testing.T) {
	tmpDir := t.TempDir()

	cmd := exec.Command(binaryPath, "sync", "diff")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	_ = err
	out := string(output)
	if !strings.Contains(out, "No .copier-answers.yml") {
		t.Error("expected warning about missing .copier-answers.yml")
	}
}

func TestSyncAllNoProjects(t *testing.T) {
	// Set HOME to a temp dir so ~/Developer doesn't exist or is empty
	tmpHome := t.TempDir()
	devDir := filepath.Join(tmpHome, "Developer")
	os.MkdirAll(devDir, 0755)

	cmd := exec.Command(binaryPath, "sync", "--all")
	cmd.Dir = tmpHome
	cmd.Env = append(os.Environ(), "HOME="+tmpHome)
	output, err := cmd.CombinedOutput()
	_ = err
	out := string(output)
	if !strings.Contains(out, "No projects") && !strings.Contains(out, "copier not installed") {
		t.Errorf("expected 'No projects' or 'copier not installed' message, got: %s", out)
	}
}

func TestSyncInitSubcommands(t *testing.T) {
	// Verify init subcommand help works
	cmd := exec.Command(binaryPath, "sync", "init", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sync init --help failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "Initialize project from template") {
		t.Error("expected description in sync init help")
	}
}
