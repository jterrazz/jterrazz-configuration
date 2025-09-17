#!/bin/bash

# Dock module for jterrazz command system
# This file defines all dock-related commands

# Main dock command handler
j_dock() {
    if [ $# -eq 0 ]; then
        j_dock_help
        return 1
    fi

    local subcommand="$1"
    shift

    case "$subcommand" in
        "add-spacer")
            echo "🔧 Adding spacer to macOS Dock..."
            defaults write com.apple.dock persistent-apps -array-add '{"tile-type"="small-spacer-tile";}'
            killall Dock
            echo "✅ Dock spacer added and restarted"
            ;;
        "reset")
            echo "🔧 Resetting macOS Dock to defaults..."
            defaults delete com.apple.dock
            killall Dock
            echo "✅ Dock reset to defaults"
            ;;
        "help"|"-h"|"--help")
            j_dock_help
            ;;
        *)
            echo "❌ Unknown dock subcommand: $subcommand"
            j_dock_help
            return 1
            ;;
    esac
}

# Dock help function
j_dock_help() {
    echo "🖥️  macOS Dock Commands"
    echo ""
    echo "Usage: j dock <subcommand>"
    echo ""
    echo "Subcommands:"
    echo "  add-spacer    Add a small spacer tile to the dock"
    echo "  reset         Reset dock to system defaults"
    echo "  help          Show this help"
    echo ""
    echo "Examples:"
    echo "  j dock add-spacer    # Add spacer between apps"
    echo "  j dock reset         # Reset dock layout"
}

# Auto-completion for dock subcommands
j_dock_completion() {
    echo "add-spacer reset help"
}

# Module metadata
J_MODULE_NAME="dock"
J_MODULE_DESCRIPTION="Manage macOS Dock (add spacers, reset)"
J_MODULE_COMMANDS="add-spacer reset help"
