# Setup of Binaries

## Unified Command System (`j`) - Modular

A unified, modular command interface for all jterrazz tools and aliases. Each category is a separate module that auto-loads.

### Usage

```shell
j <category> <command> [args...]
```

### Available Commands

**Docker:**

```shell
j docker rm        # Remove all containers
j docker rmi       # Remove all images
j docker clean     # Clean up system (prune)
j docker reset     # Remove containers + images
j docker ps        # List containers
j docker images    # List images
```

**Git:**

```shell
j git feat "msg"   # Add all and commit with feat: prefix
j git fix "msg"    # Add all and commit with fix: prefix
j git chore "msg"  # Add all and commit with chore: prefix
j git push         # Push current branch to origin
j git sync         # Fetch and pull from remote
j git wip          # Quick WIP commit
j git unwip        # Undo WIP commit
j git status       # Git status
j git log          # Recent commits
j git branches     # List branches
```

**System:**

```shell
j system update    # Update system packages (Homebrew + npm global)
j system clean     # Clean system caches and trash
j system info      # Show system information
```

**Help:**

```shell
j help             # Show all categories
j docker help      # Show docker commands
j git help         # Show git commands
j system help      # Show system commands
```

### Modular Architecture

The system automatically discovers modules from `configuration/binaries/j/*/commands.sh`:

```
configuration/binaries/j/
├── j.sh                    # Main orchestrator
├── docker/
│   └── commands.sh         # Docker module
├── git/
│   └── commands.sh         # Git module
└── system/
    └── commands.sh         # System module
```

### Adding New Modules

1. Create `configuration/binaries/j/mymodule/commands.sh`
2. Implement the required functions (see `EXTENDING.md`)
3. Make it executable: `chmod +x commands.sh`
4. Reload: the module auto-loads on next `j` command

### Setup

Add to your `~/.zshrc`:

```shell
source ~/Developer/jterrazz-configuration/configuration/binaries/zsh/zshrc.sh
```

## Oh My Zsh

> Command line tools.

https://ohmyz.sh/#install

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

And use these settings

```shell
# Add to ~/.zshrc
source ~/Developer/jterrazz-configuration/scripts/zshrc.sh
```

## Brew

> Package manager for MacOS.

https://brew.sh

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## Nvm

> Node version manager.

https://github.com/nvm-sh/nvm

```shell
brew install nvm
nvm alias default stable
```

## Git

### Create the SSH key

```shell
ssh-keygen -t ed25519 -C "contact@jterrazz.com"
eval "$(ssh-agent -s)"
```

Link this key to the GitHub host.

```shell
# vim ~/.ssh/config

Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_github
```

### Use the SSH key

```shell
ssh-add --apple-use-keychain ~/.ssh/id_github
```
