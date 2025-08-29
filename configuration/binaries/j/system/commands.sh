#!/bin/bash

# System module for jterrazz command system
# This file defines all system-related commands

# Main system command handler
j_system() {
    if [ $# -eq 0 ]; then
        j_system_help
        return 1
    fi

    local command="$1"
    shift

    case "$command" in
        "update")
            echo "🔄 Updating system packages..."
            
            # Update Homebrew packages
            if command -v brew >/dev/null 2>&1; then
                echo "🍺 Updating Homebrew packages..."
                brew update && brew upgrade
            else
                echo "❌ Homebrew not found"
                return 1
            fi
            
            # Update npm global packages
            if command -v npm >/dev/null 2>&1; then
                echo "📦 Updating npm global packages..."
                npm update -g
            else
                echo "⚠️  npm not found, skipping global package updates"
            fi
            ;;
        "install")
            if [ $# -eq 0 ]; then
                j_system_install_help
                return 1
            fi
            j_system_install "$@"
            ;;
        "clean")
            echo "🧹 Cleaning system..."
            
            # Clean Homebrew
            if command -v brew >/dev/null 2>&1; then
                echo "🍺 Cleaning Homebrew cache..."
                brew cleanup
            fi
            
            # Clean Docker
            if command -v docker >/dev/null 2>&1; then
                echo "🐳 Cleaning Docker..."
                
                # Remove stopped containers
                echo "  🗑️  Removing stopped containers..."
                docker container prune -f
                
                # Remove unused images
                echo "  🗑️  Removing unused images..."
                docker image prune -f
                
                # Remove unused volumes
                echo "  🗑️  Removing unused volumes..."
                docker volume prune -f
                
                # Remove unused networks
                echo "  🗑️  Removing unused networks..."
                docker network prune -f
                
                # Clean build cache
                echo "  🗑️  Cleaning build cache..."
                docker builder prune -f
                
                echo "  ✅ Docker cleanup completed"
            else
                echo "  ⚠️  Docker not found, skipping Docker cleanup"
            fi
            
            # Clean Multipass
            if command -v multipass >/dev/null 2>&1; then
                echo "🖥️  Cleaning Multipass..."
                
                # List and delete all instances
                echo "  🗑️  Removing all instances..."
                multipass delete --all 2>/dev/null || true
                
                # Purge deleted instances
                echo "  🗑️  Purging deleted instances..."
                multipass purge
                
                echo "  ✅ Multipass cleanup completed"
            else
                echo "  ⚠️  Multipass not found, skipping Multipass cleanup"
            fi
            
            # Empty trash
            echo "🗑️  Emptying trash..."
            rm -rf ~/.Trash/*
            
            echo "✅ System cleanup completed"
            ;;
        "help"|"-h"|"--help")
            j_system_help
            ;;
        *)
            echo "❌ Unknown system command: $command"
            j_system_help
            return 1
            ;;
    esac
}

# System install function with subcommands
j_system_install() {
    if [ $# -eq 0 ]; then
        j_system_install_help
        return 1
    fi

    local subcommand="$1"
    shift

    case "$subcommand" in
        "brew")
            j_system_install_brew
            ;;
        "ohmyzsh")
            j_system_install_ohmyzsh
            ;;
        "nvm")
            j_system_install_nvm
            ;;
        "git-ssh")
            j_system_install_git_ssh
            ;;
        "all")
            echo "🚀 Installing full development environment..."
            j_system_install_brew && \
            j_system_install_ohmyzsh && \
            j_system_install_nvm && \
            j_system_install_git_ssh
            ;;
        "help"|"-h"|"--help")
            j_system_install_help
            ;;
        *)
            echo "❌ Unknown install subcommand: $subcommand"
            j_system_install_help
            return 1
            ;;
    esac
}

# Install Homebrew
j_system_install_brew() {
    echo "🍺 Installing Homebrew..."
    if command -v brew >/dev/null 2>&1; then
        echo "✅ Homebrew already installed"
    else
        echo "📥 Downloading and installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        
        if command -v brew >/dev/null 2>&1; then
            echo "✅ Homebrew installed successfully"
        else
            echo "❌ Homebrew installation failed"
            return 1
        fi
    fi
    
    # Install essential development packages
    echo "📦 Installing essential development packages..."
    
    # ansible-lint
    if command -v ansible-lint >/dev/null 2>&1; then
        echo "  ✅ ansible-lint already installed"
    else
        echo "  📥 Installing ansible-lint..."
        brew install ansible-lint
    fi
    
    # ansible
    if command -v ansible >/dev/null 2>&1; then
        echo "  ✅ ansible already installed"
    else
        echo "  📥 Installing ansible..."
        brew install ansible
    fi
    
    # terraform
    if command -v terraform >/dev/null 2>&1; then
        echo "  ✅ terraform already installed"
    else
        echo "  📥 Installing terraform..."
        brew install terraform
    fi
    
    # kubectl
    if command -v kubectl >/dev/null 2>&1; then
        echo "  ✅ kubectl already installed"
    else
        echo "  📥 Installing kubectl..."
        brew install kubectl
    fi
    
    # multipass
    if command -v multipass >/dev/null 2>&1; then
        echo "  ✅ multipass already installed"
    else
        echo "  📥 Installing multipass..."
        brew install multipass
    fi
    
    # biome
    if command -v biome >/dev/null 2>&1; then
        echo "  ✅ biome already installed"
    else
        echo "  📥 Installing biome..."
        brew install biome
    fi
    
    echo "✅ Development packages check completed"
}

# Install Oh My Zsh
j_system_install_ohmyzsh() {
    echo "🐚 Installing Oh My Zsh..."
    if [ -d "$HOME/.oh-my-zsh" ]; then
        echo "✅ Oh My Zsh already installed"
        return 0
    fi
    
    echo "📥 Downloading and installing Oh My Zsh..."
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
    
    echo "⚙️  Configuring zshrc..."
    if ! grep -q "source ~/Developer/jterrazz-configuration/configuration/binaries/zsh/zshrc.sh" ~/.zshrc; then
        echo "" >> ~/.zshrc
        echo "# jterrazz configuration" >> ~/.zshrc
        echo "source ~/Developer/jterrazz-configuration/configuration/binaries/zsh/zshrc.sh" >> ~/.zshrc
        echo "✅ Added jterrazz configuration to ~/.zshrc"
    else
        echo "✅ jterrazz configuration already in ~/.zshrc"
    fi
    
    if [ -d "$HOME/.oh-my-zsh" ]; then
        echo "✅ Oh My Zsh installed successfully"
    else
        echo "❌ Oh My Zsh installation failed"
        return 1
    fi
}

# Install NVM
j_system_install_nvm() {
    echo "📦 Installing NVM (Node Version Manager)..."
    
    if ! command -v brew >/dev/null 2>&1; then
        echo "❌ Homebrew required for NVM installation"
        echo "💡 Run: j system install brew"
        return 1
    fi
    
    if command -v nvm >/dev/null 2>&1; then
        echo "✅ NVM already installed"
        return 0
    fi
    
    echo "📥 Installing NVM via Homebrew..."
    brew install nvm
    
    echo "⚙️  Setting up NVM..."
    # Create nvm directory if it doesn't exist
    mkdir -p ~/.nvm
    
    # Add NVM configuration to shell profile if not already present
    local nvm_config='
# NVM Configuration
export NVM_DIR="$HOME/.nvm"
[ -s "/opt/homebrew/opt/nvm/nvm.sh" ] && \. "/opt/homebrew/opt/nvm/nvm.sh"
[ -s "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm" ] && \. "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm"'
    
    if ! grep -q "NVM Configuration" ~/.zshrc; then
        echo "$nvm_config" >> ~/.zshrc
        echo "✅ Added NVM configuration to ~/.zshrc"
    else
        echo "✅ NVM configuration already in ~/.zshrc"
    fi
    
    # Source the NVM script for immediate use
    export NVM_DIR="$HOME/.nvm"
    [ -s "/opt/homebrew/opt/nvm/nvm.sh" ] && \. "/opt/homebrew/opt/nvm/nvm.sh"
    
    if command -v nvm >/dev/null 2>&1; then
        echo "📥 Installing Node.js stable..."
        nvm install stable
        nvm alias default stable
        nvm use stable
        echo "✅ NVM and Node.js installed successfully"
    else
        echo "❌ NVM installation failed - restart terminal and try again"
        return 1
    fi
}

# Install Git SSH setup
j_system_install_git_ssh() {
    echo "🔑 Setting up Git SSH..."
    
    local ssh_key="$HOME/.ssh/id_github"
    local email="contact@jterrazz.com"
    
    # Check if SSH key already exists
    if [ -f "$ssh_key" ]; then
        echo "✅ SSH key already exists at $ssh_key"
    else
        echo "🔐 Generating SSH key..."
        ssh-keygen -t ed25519 -C "$email" -f "$ssh_key" -N ""
        echo "✅ SSH key generated"
    fi
    
    # Configure SSH
    echo "⚙️  Configuring SSH..."
    local ssh_config="$HOME/.ssh/config"
    
    if ! grep -q "Host github.com" "$ssh_config" 2>/dev/null; then
        cat >> "$ssh_config" << EOF

Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_github
EOF
        echo "✅ SSH config updated"
    else
        echo "✅ SSH config already configured"
    fi
    
    # Add key to SSH agent
    echo "🔗 Adding key to SSH agent..."
    eval "$(ssh-agent -s)"
    ssh-add --apple-use-keychain "$ssh_key"
    
    echo "📋 Your public key (add this to GitHub):"
    echo "----------------------------------------"
    cat "${ssh_key}.pub"
    echo "----------------------------------------"
    echo "💡 Copy the above key and add it to: https://github.com/settings/ssh/new"
    
    echo "✅ Git SSH setup completed"
}

# Install help function
j_system_install_help() {
    echo "📦 System Install Commands"
    echo ""
    echo "Usage: j system install <subcommand>"
    echo ""
    echo "Subcommands:"
    echo "  brew      Install Homebrew + development packages (ansible, terraform, kubectl, multipass, biome)"
    echo "  ohmyzsh   Install Oh My Zsh and configure shell"
    echo "  nvm       Install NVM and Node.js stable"
    echo "  git-ssh   Generate SSH key and configure Git"
    echo "  all       Install everything above in order"
    echo "  help      Show this help"
    echo ""
    echo "Examples:"
    echo "  j system install all       # Full development environment"
    echo "  j system install brew      # Homebrew + dev tools (ansible, terraform, kubectl, multipass, biome)"
    echo "  j system install git-ssh   # Just Git SSH setup"
}

# System help function
j_system_help() {
    echo "⚙️  System Commands"
    echo ""
    echo "Usage: j system <command>"
    echo ""
    echo "Commands:"
    echo "  update    Update system packages (Homebrew + npm global)"
    echo "  install   Install development tools (brew, ansible, terraform, kubectl, biome, etc.)"
    echo "  clean     Clean system caches, Docker, Multipass, and trash"
    echo "  help      Show this help"
    echo ""
    echo "⚠️  WARNING: 'clean' removes all Multipass instances and Docker containers"
    echo "💡 For comprehensive system status, use: j status"
}

# Auto-completion for system commands
j_system_completion() {
    echo "update install clean help"
}

# Auto-completion for system install subcommands
j_system_install_completion() {
    echo "brew ohmyzsh nvm git-ssh all help"
}

# Module metadata
J_MODULE_NAME="system"
J_MODULE_DESCRIPTION="System maintenance and development tools installation"
J_MODULE_COMMANDS="update install clean help"