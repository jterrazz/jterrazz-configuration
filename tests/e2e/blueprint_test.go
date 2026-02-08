package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// =============================================================================
// Blueprint variant test infrastructure
// =============================================================================

type blueprintVariant struct {
	name          string
	data          map[string]string
	wantFiles     []string
	excludedFiles []string
}

func runBlueprintVariant(t *testing.T, v blueprintVariant) {
	t.Helper()
	RequireCopier(t)

	// Snapshot test: compare generated output against committed fixture
	CompareWithFixture(t, v.name, v.data)

	// Also verify expected/excluded files against the fixture
	fixtureDir := BlueprintFixtureDir(v.name)
	files := ListFiles(t, fixtureDir)

	fileSet := make(map[string]bool, len(files))
	for _, f := range files {
		fileSet[f] = true
	}

	for _, want := range v.wantFiles {
		if !fileSet[want] {
			t.Errorf("expected file %q not found in fixture", want)
		}
	}

	for _, excluded := range v.excludedFiles {
		if fileSet[excluded] {
			t.Errorf("file %q should not exist in fixture", excluded)
		}
	}
}

// =============================================================================
// Shared file lists
// =============================================================================

var (
	commonFiles = []string{
		".editorconfig",
		".gitignore",
		"LICENSE",
	}

	tsFiles = []string{
		".nvmrc",
		"tsconfig.json",
	}

	tsWithPkg = []string{
		"package.json",
		"vitest.config.ts",
	}

	goFiles = []string{
		"go.mod",
		"Makefile",
		".golangci.yml",
	}

	dockerFiles = []string{
		"Dockerfile",
		".dockerignore",
	}

	ciWorkflow      = ".github/workflows/ci.yml"
	releaseWorkflow = ".github/workflows/release.yml"
	deployWorkflow  = ".github/workflows/deploy.yml"
)

// =============================================================================
// None — license variants
// =============================================================================

func TestBlueprintNoneMit(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "none-mit",
		data: map[string]string{
			"project_name": "my-config",
			"language":     "none",
			"project_type": "none",
			"license":      "MIT",
			"ci":           "false",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: commonFiles,
		excludedFiles: MergeFiles(tsFiles, tsWithPkg, goFiles, dockerFiles, []string{
ciWorkflow, releaseWorkflow, deployWorkflow,
		}),
	})
}

func TestBlueprintNoneProprietary(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "none-proprietary",
		data: map[string]string{
			"project_name": "my-config",
			"language":     "none",
			"project_type": "none",
			"license":      "proprietary",
			"ci":           "false",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: commonFiles,
		excludedFiles: MergeFiles(tsFiles, tsWithPkg, goFiles, dockerFiles, []string{
ciWorkflow, releaseWorkflow, deployWorkflow,
		}),
	})
}

// =============================================================================
// TypeScript — language and project types
// =============================================================================

func TestBlueprintTypescriptNone(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "typescript-none",
		data: map[string]string{
			"project_name": "my-ts",
			"language":     "typescript",
			"project_type": "none",
			"license":      "MIT",
			"ci":           "false",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, tsFiles, tsWithPkg),
		excludedFiles: MergeFiles(goFiles, dockerFiles, []string{
ciWorkflow, releaseWorkflow, deployWorkflow,
		}),
	})
}

func TestBlueprintTypescriptLibrary(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "typescript-library",
		data: map[string]string{
			"project_name": "my-lib",
			"language":     "typescript",
			"project_type": "library",
			"license":      "MIT",
			"ci":           "true",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, tsFiles, tsWithPkg, []string{
ciWorkflow, releaseWorkflow,
		}),
		excludedFiles: MergeFiles(goFiles, dockerFiles, []string{deployWorkflow}),
	})
}

func TestBlueprintTypescriptApi(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "typescript-api",
		data: map[string]string{
			"project_name": "my-api",
			"language":     "typescript",
			"project_type": "api",
			"license":      "MIT",
			"ci":           "true",
			"docker":       "true",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, tsFiles, tsWithPkg, dockerFiles, []string{
ciWorkflow,
		}),
		excludedFiles: MergeFiles(goFiles, []string{releaseWorkflow, deployWorkflow}),
	})
}

func TestBlueprintTypescriptWeb(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "typescript-web",
		data: map[string]string{
			"project_name": "my-web",
			"language":     "typescript",
			"project_type": "web",
			"license":      "MIT",
			"ci":           "true",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, tsFiles, tsWithPkg, []string{
ciWorkflow,
		}),
		excludedFiles: MergeFiles(goFiles, dockerFiles, []string{releaseWorkflow, deployWorkflow}),
	})
}

func TestBlueprintTypescriptMobile(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "typescript-mobile",
		data: map[string]string{
			"project_name": "my-mobile",
			"language":     "typescript",
			"project_type": "mobile",
			"license":      "MIT",
			"ci":           "true",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, tsFiles, []string{
ciWorkflow,
		}),
		excludedFiles: MergeFiles(goFiles, dockerFiles, tsWithPkg, []string{releaseWorkflow, deployWorkflow}),
	})
}

func TestBlueprintTypescriptApiDeploy(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "typescript-api-deploy",
		data: map[string]string{
			"project_name": "my-api-deploy",
			"language":     "typescript",
			"project_type": "api",
			"license":      "MIT",
			"ci":           "true",
			"docker":       "true",
			"deploy":       "kubernetes",
		},
		wantFiles: MergeFiles(commonFiles, tsFiles, tsWithPkg, dockerFiles, []string{
ciWorkflow, deployWorkflow,
		}),
		excludedFiles: MergeFiles(goFiles, []string{releaseWorkflow}),
	})
}

// =============================================================================
// Go — language and project types
// =============================================================================

func TestBlueprintGoNone(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "go-none",
		data: map[string]string{
			"project_name": "my-go",
			"language":     "go",
			"project_type": "none",
			"license":      "MIT",
			"ci":           "false",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, goFiles),
		excludedFiles: MergeFiles(tsFiles, tsWithPkg, dockerFiles, []string{
ciWorkflow, releaseWorkflow, deployWorkflow,
		}),
	})
}

func TestBlueprintGoCli(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "go-cli",
		data: map[string]string{
			"project_name": "my-go-cli",
			"language":     "go",
			"project_type": "cli",
			"license":      "MIT",
			"ci":           "true",
			"docker":       "false",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, goFiles, []string{
ciWorkflow,
		}),
		excludedFiles: MergeFiles(tsFiles, tsWithPkg, dockerFiles, []string{releaseWorkflow, deployWorkflow}),
	})
}

func TestBlueprintGoApi(t *testing.T) {
	runBlueprintVariant(t, blueprintVariant{
		name: "go-api",
		data: map[string]string{
			"project_name": "my-go-api",
			"language":     "go",
			"project_type": "api",
			"license":      "MIT",
			"ci":           "true",
			"docker":       "true",
			"deploy":       "none",
		},
		wantFiles: MergeFiles(commonFiles, goFiles, dockerFiles, []string{
ciWorkflow,
		}),
		excludedFiles: MergeFiles(tsFiles, tsWithPkg, []string{releaseWorkflow, deployWorkflow}),
	})
}

// =============================================================================
// Content validation (reads from committed fixtures)
// =============================================================================

func TestBlueprintLicenseMIT(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(BlueprintFixtureDir("none-mit"), "LICENSE"))
	if err != nil {
		t.Fatalf("failed to read LICENSE: %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "MIT License") {
		t.Error("expected 'MIT License' in LICENSE")
	}
	if !strings.Contains(text, "Jean-Baptiste Terrazzoni") {
		t.Error("expected author name in LICENSE")
	}
}

func TestBlueprintLicenseProprietary(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(BlueprintFixtureDir("none-proprietary"), "LICENSE"))
	if err != nil {
		t.Fatalf("failed to read LICENSE: %v", err)
	}
	if !strings.Contains(string(content), "All Rights Reserved") {
		t.Error("expected 'All Rights Reserved' in LICENSE")
	}
}

func TestBlueprintTsconfigVariants(t *testing.T) {
	cases := []struct {
		fixture     string
		wantExtends string
	}{
		{"typescript-library", "@jterrazz/typescript/tsconfig/node"},
		{"typescript-api", "@jterrazz/typescript/tsconfig/node"},
		{"typescript-web", "@jterrazz/typescript/tsconfig/next.json"},
		{"typescript-mobile", "@jterrazz/typescript/tsconfig/expo"},
	}

	for _, tc := range cases {
		t.Run(tc.fixture, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(BlueprintFixtureDir(tc.fixture), "tsconfig.json"))
			if err != nil {
				t.Fatalf("failed to read tsconfig.json: %v", err)
			}
			if !strings.Contains(string(content), tc.wantExtends) {
				t.Errorf("expected tsconfig to extend %q, got:\n%s", tc.wantExtends, content)
			}
		})
	}
}

func TestBlueprintWorkflowRawEscaping(t *testing.T) {
	workflows, _ := filepath.Glob(filepath.Join(BlueprintFixtureDir("typescript-library"), ".github/workflows/*.yml"))
	if len(workflows) == 0 {
		t.Fatal("no workflow files in typescript-library fixture")
	}

	for _, wf := range workflows {
		content, err := os.ReadFile(wf)
		if err != nil {
			t.Fatalf("failed to read %s: %v", wf, err)
		}
		text := string(content)
		if strings.Contains(text, "{% raw %}") {
			t.Errorf("%s still contains {%% raw %%} tags", filepath.Base(wf))
		}
		if strings.Contains(text, "{% endraw %}") {
			t.Errorf("%s still contains {%% endraw %%} tags", filepath.Base(wf))
		}
	}
}
