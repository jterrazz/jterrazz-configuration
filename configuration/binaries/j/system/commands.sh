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
        "clean")
            echo "🧹 Cleaning system..."
            if command -v brew >/dev/null 2>&1; then
                brew cleanup
            fi
            echo "🗑️  Emptying trash..."
            rm -rf ~/.Trash/*
            ;;
        "info")
            echo "ℹ️  System information:"
            echo "OS: $(uname -s) $(uname -r)"
            echo "User: $(whoami)"
            echo "Shell: $SHELL"
            if command -v brew >/dev/null 2>&1; then
                echo "Homebrew: $(brew --version | head -1)"
            fi
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

# System help function
j_system_help() {
    echo "⚙️  System Commands"
    echo ""
    echo "Usage: j system <command>"
    echo ""
    echo "Commands:"
    echo "  update    Update system packages (Homebrew + npm global)"
    echo "  clean     Clean system caches and trash"
    echo "  info      Show system information"
    echo "  help      Show this help"
}

# Auto-completion for system commands
j_system_completion() {
    echo "update clean info help"
}

# Module metadata
J_MODULE_NAME="system"
J_MODULE_DESCRIPTION="System maintenance and information"
J_MODULE_COMMANDS="update clean info help"