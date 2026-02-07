# jterrazz-cli

A unified CLI tool (`j`) for development workflow automation, plus dotfiles and machine setup.

_Hey there – I'm Jean-Baptiste, just another developer doing weird things with code. All my projects live on [jterrazz.com](https://jterrazz.com) – complete with backstories and lessons learned. Feel free to poke around – you might just find something useful!_

## Installation

```bash
git clone https://github.com/jterrazz/jterrazz-cli.git
cd jterrazz-cli
make install
```

Requires Go 1.21+. Install Go via `brew install go` if needed. The binary is installed to `/usr/local/bin/j`.

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

Sync configuration files (.gitignore, tsconfig, LICENSE, etc.) across repositories using [Copier](https://github.com/copier-org/copier) templates stored in `dotfiles/blueprints/`.

```bash
j sync init          # Initialize project from template (auto-detects language)
j sync               # Update project from its template
j sync status        # Show template link status
j sync diff          # Preview changes before updating
j sync --all         # Update all projects in ~/Developer
```

**How it works:** Running `j sync init` asks a few questions (language, license, CI, etc.) and generates config files. A `.copier-answers.yml` file is created in the project to track the template version and your answers — commit this file. When templates are updated and tagged, run `j sync` in any linked project to pull the latest changes.

**Included templates:** .editorconfig, .gitattributes, .gitignore, LICENSE, plus conditional files for TypeScript (tsconfig, .npmrc, .nvmrc), Go (Makefile, .golangci.yml), CI (GitHub Actions), and Docker.

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
make test     # Run tests
make clean    # Remove build artifacts
```

## Structure

```
.
├── src/
│   ├── cmd/j/main.go          # Entry point
│   ├── internal/
│   │   ├── commands/           # CLI commands (Cobra)
│   │   │   ├── root.go
│   │   │   ├── status.go
│   │   │   ├── install.go
│   │   │   ├── upgrade.go
│   │   │   ├── clean.go
│   │   │   ├── setup.go
│   │   │   ├── run.go
│   │   │   └── sync.go
│   │   ├── config/             # Tool/script/command registry
│   │   ├── domain/             # Business logic
│   │   └── presentation/       # TUI components and views
│   └── go.mod
├── dotfiles/
│   ├── applications/           # App configs (Cursor, Ghostty, Zed, Zsh)
│   └── blueprints/             # Copier project templates
├── tests/e2e/                  # End-to-end CLI tests
└── Makefile
```
