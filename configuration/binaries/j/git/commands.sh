#!/bin/bash

# Git module for jterrazz command system
# This file defines all git-related commands

# Main git command handler
j_git() {
    if [ $# -eq 0 ]; then
        j_git_help
        return 1
    fi

    local command="$1"
    shift

    case "$command" in
        "feat")
            if [ $# -eq 0 ]; then
                echo "❌ Please provide a commit message"
                echo "💡 Usage: j git feat 'your message'"
                return 1
            fi
            git add . && git commit -m "feat: $*"
            ;;
        "fix")
            if [ $# -eq 0 ]; then
                echo "❌ Please provide a commit message"
                echo "💡 Usage: j git fix 'your message'"
                return 1
            fi
            git add . && git commit -m "fix: $*"
            ;;
        "chore")
            if [ $# -eq 0 ]; then
                echo "❌ Please provide a commit message"
                echo "💡 Usage: j git chore 'your message'"
                return 1
            fi
            git add . && git commit -m "chore: $*"
            ;;
        "push")
            git push -u origin HEAD
            ;;
        "sync")
            echo "🔄 Syncing with remote..."
            git fetch -p && git pull
            ;;
        "wip")
            git add --all && git commit -m "WIP"
            ;;
        "unwip")
            git reset --soft HEAD~1 && git reset HEAD
            ;;
        "status")
            git status
            ;;
        "log")
            git log --oneline -10
            ;;
        "branches")
            echo "🌿 Local branches:"
            git branch
            ;;
        "help"|"-h"|"--help")
            j_git_help
            ;;
        *)
            echo "❌ Unknown git command: $command"
            j_git_help
            return 1
            ;;
    esac
}

# Git help function
j_git_help() {
    echo "🌱 Git Commands"
    echo ""
    echo "Usage: j git <command> [args...]"
    echo ""
    echo "Commands:"
    echo "  feat <msg>    Add all and commit with 'feat:' prefix"
    echo "  fix <msg>     Add all and commit with 'fix:' prefix"
    echo "  chore <msg>   Add all and commit with 'chore:' prefix"
    echo "  push          Push current branch to origin"
    echo "  sync          Fetch and pull from remote"
    echo "  wip           Add all and commit as 'WIP'"
    echo "  unwip         Undo last commit and unstage"
    echo "  status        Show git status"
    echo "  log           Show recent commits"
    echo "  branches      List local branches"
    echo "  help          Show this help"
}

# Auto-completion for git commands
j_git_completion() {
    echo "feat fix chore push sync wip unwip status log branches help"
}

# Module metadata
J_MODULE_NAME="git"
J_MODULE_DESCRIPTION="Git workflow shortcuts"
J_MODULE_COMMANDS="feat fix chore push sync wip unwip status log branches help"