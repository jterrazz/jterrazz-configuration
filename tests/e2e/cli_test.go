package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Run with: cd tests && go test -v ./e2e/

func TestMain(m *testing.M) {
	BuildBinary()
	code := m.Run()
	CleanupBinary()
	os.Exit(code)
}

func TestStatusCommand(t *testing.T) {
	t.Skip("skipped: j status launches a TUI that blocks without a terminal")
	out := RunCLI(t, "status")
	if !strings.Contains(out, "j status") {
		t.Error("expected 'j status' header in output")
	}
}

func TestInstallCommand(t *testing.T) {
	out := RunCLI(t, "install", "--help")
	if !strings.Contains(out, "install") {
		t.Error("expected 'install' in help output")
	}
}

func TestCleanCommand(t *testing.T) {
	out := RunCLI(t, "clean", "--help")
	if !strings.Contains(out, "clean") {
		t.Error("expected 'clean' in help output")
	}
}

func TestUpgradeCommand(t *testing.T) {
	out := RunCLI(t, "upgrade", "--help")
	if !strings.Contains(out, "upgrade") {
		t.Error("expected 'upgrade' in help output")
	}
}

func TestRunGitCommands(t *testing.T) {
	out := RunCLI(t, "run", "git", "--help")
	for _, sub := range []string{"feat", "fix", "chore", "push", "sync"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected '%s' in git help output", sub)
		}
	}
}

func TestRunDockerCommands(t *testing.T) {
	out := RunCLI(t, "run", "docker", "--help")
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
	out := RunCLI(t, "sync", "--help")
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
	tmpDir := t.TempDir()
	out, err := RunCLIInDir(t, tmpDir, "sync", "status")
	if err != nil {
		t.Fatalf("sync status failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Not linked") {
		t.Error("expected 'Not linked' for directory without .copier-answers.yml")
	}
	if !strings.Contains(out, "j sync init") {
		t.Error("expected hint to run 'j sync init'")
	}
}

func TestSyncStatusLinked(t *testing.T) {
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

	out, err := RunCLIInDir(t, tmpDir, "sync", "status")
	if err != nil {
		t.Fatalf("sync status failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Linked") {
		t.Error("expected 'Linked' for directory with .copier-answers.yml")
	}
	for _, want := range []string{"project_name", "test-project", "language"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in status output", want)
		}
	}
}

func TestSyncUpdateUnlinked(t *testing.T) {
	tmpDir := t.TempDir()
	out, _ := RunCLIInDir(t, tmpDir, "sync")
	if !strings.Contains(out, "No .copier-answers.yml") || !strings.Contains(out, "j sync init") {
		t.Error("expected warning about missing .copier-answers.yml with hint to run init")
	}
}

func TestSyncDiffUnlinked(t *testing.T) {
	tmpDir := t.TempDir()
	out, _ := RunCLIInDir(t, tmpDir, "sync", "diff")
	if !strings.Contains(out, "No .copier-answers.yml") {
		t.Error("expected warning about missing .copier-answers.yml")
	}
}

func TestSyncAllNoProjects(t *testing.T) {
	tmpHome := t.TempDir()
	devDir := filepath.Join(tmpHome, "Developer")
	os.MkdirAll(devDir, 0755)

	cmd := runCLICmdInDir(tmpHome, "sync", "--all")
	cmd.Env = append(os.Environ(), "HOME="+tmpHome)
	output, err := cmd.CombinedOutput()
	_ = err
	out := string(output)
	if !strings.Contains(out, "No projects") && !strings.Contains(out, "copier not installed") {
		t.Errorf("expected 'No projects' or 'copier not installed' message, got: %s", out)
	}
}

func TestSyncInitSubcommands(t *testing.T) {
	out := RunCLI(t, "sync", "init", "--help")
	if !strings.Contains(out, "Initialize project from template") {
		t.Error("expected description in sync init help")
	}
}

// runCLICmdInDir returns an exec.Cmd for the binary (for cases needing env customization).
func runCLICmdInDir(dir string, args ...string) *exec.Cmd {
	cmd := exec.Command(BinaryPath, args...)
	cmd.Dir = dir
	return cmd
}
