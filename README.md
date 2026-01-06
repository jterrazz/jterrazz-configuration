# jterrazz-configuration

A unified CLI tool (`j`) for development workflow automation, plus dotfiles and machine setup.

*Hey there – I'm Jean-Baptiste, just another developer doing weird things with code. All my projects live on [jterrazz.com](https://jterrazz.com) – complete with backstories and lessons learned. Feel free to poke around – you might just find something useful!*

## Installation

```bash
git clone https://github.com/jterrazz/jterrazz-configuration.git
cd jterrazz-configuration
make install
```

Requires Go 1.21+. Install Go via `brew install go` if needed. The binary is installed to `/usr/local/bin/j`.

## Usage

```bash
j help              # Show all commands
j status            # Show system status
j update            # Update Homebrew + npm packages
j clean             # Clean caches, Docker, Multipass, trash
```

### Install (Packages)

```bash
j install                  # List available packages
j install --all            # Install all packages
j install brew             # Install Homebrew
j install go python node   # Install specific packages
```

### Setup (Configurations)

```bash
j setup ohmyzsh      # Configure Oh My Zsh
j setup git-ssh      # Setup Git SSH keys
j setup dock-spacer  # Add spacer to dock
j setup dock-reset   # Reset dock to defaults
j setup all          # Setup all configurations
```

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
│   ├── cmd/j/main.go       # Entry point
│   ├── internal/commands/  # Command implementations
│   │   ├── root.go
│   │   ├── status.go
│   │   ├── update.go
│   │   ├── clean.go
│   │   ├── setup.go
│   │   ├── git.go
│   │   └── docker.go
│   └── go.mod
├── configuration/
│   ├── applications/       # App configs (Cursor, etc.)
│   └── binaries/zsh/       # Shell config sourced by ~/.zshrc
└── Makefile
```
