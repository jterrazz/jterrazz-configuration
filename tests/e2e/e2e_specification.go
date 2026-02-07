// Package e2e provides shared utilities for end-to-end tests.
package e2e

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var (
	// BinaryPath is the absolute path to the built j binary. Set by TestMain.
	BinaryPath string

	// RepoRoot is the absolute path to the repository root.
	RepoRoot string

	// TemplatePath is the absolute path to the copier blueprint template.
	TemplatePath string
)

func init() {
	abs, err := filepath.Abs("../..")
	if err != nil {
		panic("failed to resolve repo root: " + err.Error())
	}
	RepoRoot = abs
	TemplatePath = filepath.Join(RepoRoot, "dotfiles", "blueprints")
}

// BuildBinary compiles the j CLI binary and sets BinaryPath.
// Call this from TestMain before m.Run().
func BuildBinary() {
	absPath, err := filepath.Abs("j_test_bin")
	if err != nil {
		panic("failed to get absolute path: " + err.Error())
	}
	BinaryPath = absPath

	cmd := exec.Command("go", "build", "-o", BinaryPath, "./src/cmd/j")
	cmd.Dir = RepoRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("failed to build binary: " + err.Error() + "\n" + string(output))
	}
}

// CleanupBinary removes the built test binary.
func CleanupBinary() {
	if BinaryPath != "" {
		os.Remove(BinaryPath)
	}
}

// RunCLI executes the j binary with the given arguments and returns stdout+stderr.
func RunCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command(BinaryPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command %v failed: %v\n%s", args, err, output)
	}
	return string(output)
}

// RunCLIInDir executes the j binary in the given directory.
func RunCLIInDir(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(BinaryPath, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// RequireCopier skips the test if copier is not installed.
func RequireCopier(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("copier"); err != nil {
		t.Skip("copier not installed, skipping")
	}
}

// CopierGenerate runs copier copy with the given data map into dest.
func CopierGenerate(t *testing.T, dest string, data map[string]string) {
	t.Helper()
	args := []string{"copy", "--trust", "--defaults", "--quiet"}
	for k, v := range data {
		args = append(args, "--data", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, TemplatePath, dest)

	cmd := exec.Command("copier", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("copier copy failed: %v\n%s", err, output)
	}
}

// ListFiles returns all relative file paths under dir, sorted.
func ListFiles(t *testing.T, dir string) []string {
	t.Helper()
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(dir, path)
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk %s: %v", dir, err)
	}
	sort.Strings(files)
	return files
}

// MergeFiles concatenates multiple string slices into one.
func MergeFiles(lists ...[]string) []string {
	var result []string
	for _, l := range lists {
		result = append(result, l...)
	}
	return result
}

// BlueprintFixtureDir returns the path to tests/e2e/output/<name>/ (committed fixtures).
func BlueprintFixtureDir(name string) string {
	return filepath.Join(RepoRoot, "tests", "e2e", "output", name)
}

// UpdateFixture controls whether tests overwrite committed fixtures with fresh output.
// Set via: go test ./e2e/ -run TestBlueprint -args -update
var UpdateFixture = flag.Bool("update", false, "update committed fixtures with fresh copier output")

// CompareWithFixture generates a blueprint into a temp dir, then compares it
// file-by-file against the committed fixture at tests/e2e/output/<name>/.
// If -update is passed, it overwrites the fixture instead.
func CompareWithFixture(t *testing.T, name string, data map[string]string) {
	t.Helper()

	fixtureDir := BlueprintFixtureDir(name)

	if *UpdateFixture {
		// Regenerate the fixture
		if err := os.RemoveAll(fixtureDir); err != nil {
			t.Fatalf("failed to clean fixture dir: %v", err)
		}
		if err := os.MkdirAll(fixtureDir, 0755); err != nil {
			t.Fatalf("failed to create fixture dir: %v", err)
		}
		CopierGenerate(t, fixtureDir, data)
		t.Logf("Updated fixture: %s", fixtureDir)
		return
	}

	// Generate into temp dir
	tmpDir := t.TempDir()
	CopierGenerate(t, tmpDir, data)

	// Compare generated files against fixture
	generatedFiles := ListFiles(t, tmpDir)
	fixtureFiles := ListFiles(t, fixtureDir)

	genSet := make(map[string]bool)
	for _, f := range generatedFiles {
		genSet[f] = true
	}
	fixSet := make(map[string]bool)
	for _, f := range fixtureFiles {
		fixSet[f] = true
	}

	// Check for missing/extra files
	for _, f := range generatedFiles {
		if !fixSet[f] {
			t.Errorf("new file not in fixture: %s (run with -update to accept)", f)
		}
	}
	for _, f := range fixtureFiles {
		if !genSet[f] {
			t.Errorf("fixture file no longer generated: %s (run with -update to accept)", f)
		}
	}

	// Compare content of shared files
	for _, f := range generatedFiles {
		if !fixSet[f] {
			continue
		}
		genContent, err := os.ReadFile(filepath.Join(tmpDir, f))
		if err != nil {
			t.Fatalf("failed to read generated %s: %v", f, err)
		}
		fixContent, err := os.ReadFile(filepath.Join(fixtureDir, f))
		if err != nil {
			t.Fatalf("failed to read fixture %s: %v", f, err)
		}
		if string(genContent) != string(fixContent) {
			t.Errorf("file %s differs from fixture (run with -update to accept)", f)
			// Show a short diff hint
			genLines := strings.Split(string(genContent), "\n")
			fixLines := strings.Split(string(fixContent), "\n")
			maxLines := len(genLines)
			if len(fixLines) > maxLines {
				maxLines = len(fixLines)
			}
			shown := 0
			for i := 0; i < maxLines && shown < 5; i++ {
				genLine := ""
				fixLine := ""
				if i < len(genLines) {
					genLine = genLines[i]
				}
				if i < len(fixLines) {
					fixLine = fixLines[i]
				}
				if genLine != fixLine {
					t.Errorf("  line %d:\n    fixture:   %q\n    generated: %q", i+1, fixLine, genLine)
					shown++
				}
			}
		}
	}
}
