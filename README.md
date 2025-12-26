*Hey there – I'm Jean-Baptiste, just another developer doing weird things with code. All my projects live on [jterrazz.com](https://jterrazz.com) – complete with backstories and lessons learned. Feel free to poke around – you might just find something useful!*

# jterrazz-configuration

A unified CLI tool (`j`) for development workflow automation, plus dotfiles and machine setup.

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

### Setup

```bash
j setup brew         # Install Homebrew + dev packages
j setup ohmyzsh      # Install Oh My Zsh
j setup nvm          # Install NVM + Node.js
j setup git-ssh      # Setup Git SSH keys
j setup all          # Full dev environment setup
j setup dock-spacer  # Add spacer to dock
j setup dock-reset   # Reset dock to defaults
```

### Git Shortcuts

```bash
j git feat "message"   # git add . && git commit -m "feat: message"
j git fix "message"    # git add . && git commit -m "fix: message"
j git chore "message"  # git add . && git commit -m "chore: message"
j git push             # git push -u origin HEAD
j git sync             # git fetch -p && git pull
j git wip              # git add --all && git commit -m "WIP"
j git unwip            # Undo last commit and unstage
```

### Docker

```bash
j docker ps            # List containers
j docker images        # List images
j docker rm            # Remove all containers
j docker rmi           # Remove all images
j docker clean         # docker system prune -af
j docker reset         # Remove all containers and images
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
