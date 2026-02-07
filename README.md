# jterrazz-cli

A unified CLI tool (`j`) for development workflow automation, plus dotfiles and machine setup.

_Hey there – I'm Jean-Baptiste, just another developer doing weird things with code. All my projects live on [jterrazz.com](https://jterrazz.com) – complete with backstories and lessons learned. Feel free to poke around – you might just find something useful!_

## Installation

```bash
git clone https://github.com/jterrazz/jterrazz-cli.git
cd jterrazz-cli
make install
```

Requires Go 1.24+. Install Go via `brew install go` if needed. The binary is installed to `/usr/local/bin/j`.

## Usage

```bash
j help              # Show all commands
j status            # Show system status
j upgrade           # Update Homebrew + npm packages
j clean             # Clean caches, Docker, Multipass, trash
```

### Install (Packages)

```bash
j install                  # List available packages
j install brew             # Install Homebrew
j install go python node   # Install specific packages
j install copier           # Install copier template engine
```

### Setup (Configurations)

```bash
j setup              # Interactive TUI for system configuration
```

### Sync (Project Templates)

Sync configuration files across repositories using [Copier](https://github.com/copier-org/copier) templates stored in `dotfiles/blueprints/`.

```bash
j sync init          # Initialize project from template (auto-detects language)
j sync               # Update project from its template
j sync status        # Show template link status
j sync diff          # Preview changes before updating
j sync --all         # Update all projects in ~/Developer
```

**How it works:** Running `j sync init` asks a few questions (language, license, CI, etc.) and generates config files. A `.copier-answers.yml` file is created in the project to track the template version and your answers — commit this file. When templates are updated and tagged, run `j sync` in any linked project to pull the latest changes.

**Included templates:** .editorconfig, .gitattributes, .gitignore, LICENSE, plus conditional files for TypeScript (tsconfig, .nvmrc, package.json), Go (go.mod, Makefile, .golangci.yml), CI (GitHub Actions), Docker, and Claude Code skills.

### Run Commands

```bash
# Git shortcuts
j run git feat "message"   # git add . && git commit -m "feat: message"
j run git fix "message"    # git add . && git commit -m "fix: message"
j run git chore "message"  # git add . && git commit -m "chore: message"
j run git push             # git push -u origin HEAD
j run git sync             # git fetch -p && git pull
j run git wip              # git add --all && git commit -m "WIP"
j run git unwip            # Undo last commit and unstage

# Docker
j run docker rm            # Remove all containers
j run docker rmi           # Remove all images
j run docker clean         # docker system prune -af
j run docker reset         # Remove all containers and images
```

## Development

```bash
make build    # Build binary
make test     # Run unit tests
make install  # Build and install to /usr/local/bin/j
```

### E2E Tests

```bash
go test ./tests/e2e/ -v -timeout 120s           # Run all e2e tests
go test ./tests/e2e/ -run TestBlueprint -v       # Blueprint tests only
go test ./tests/e2e/ -run TestBlueprint -args -update  # Regenerate fixtures
```

Blueprint tests use committed fixtures in `tests/e2e/output/`. Each fixture isolates one feature:

| Fixture                 | Tests                                              |
| ----------------------- | -------------------------------------------------- |
| `none-mit`              | MIT license, common files                          |
| `none-proprietary`      | Proprietary license                                |
| `typescript-none`       | TypeScript language files                          |
| `typescript-library`    | Library CI + release workflow + architecture skill |
| `typescript-api`        | App CI + Docker + API architecture skill           |
| `typescript-web`        | Web architecture skill + Next.js tsconfig          |
| `typescript-mobile`     | Mobile architecture skill                          |
| `typescript-api-deploy` | Kubernetes deploy workflow                         |
| `go-none`               | Go language files                                  |
| `go-cli`                | Go CI + CLI architecture skill                     |
| `go-api`                | Go CI + Docker + API architecture skill            |

## Structure

```
.
├── src/
│   ├── cmd/j/main.go          # Entry point
│   └── internal/
│       ├── commands/           # CLI commands (Cobra)
│       ├── config/             # Tool/script/command registry
│       ├── domain/             # Business logic
│       └── presentation/       # TUI components and views
├── dotfiles/
│   ├── applications/           # App configs (Cursor, Ghostty, Zed, Zsh)
│   └── blueprints/             # Copier project templates
│       ├── copier.yml          # Template configuration
│       └── template/           # Template files
├── tests/e2e/
│   ├── blueprint_test.go       # Blueprint snapshot tests
│   ├── cli_test.go             # CLI command tests
│   ├── e2e_specification.go    # Shared test framework
│   └── output/                 # Committed fixture directories
├── go.mod
└── Makefile
```
