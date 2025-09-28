#!/bin/bash

# Docker module for jterrazz command system
# This file defines all docker-related commands

# Main docker command handler
j_docker() {
    if [ $# -eq 0 ]; then
        j_docker_help
        return 1
    fi

    local command="$1"
    shift

    case "$command" in
        "rm")
            echo "🧹 Removing all Docker containers..."
            docker rm -vf $(docker ps -aq)
            ;;
        "rmi")
            echo "🗑️  Removing all Docker images..."
            docker rmi -f $(docker images -aq)
            ;;
        "clean")
            echo "🧹 Cleaning up Docker system..."
            docker system prune -af
            ;;
        "reset")
            echo "🔄 Resetting Docker (removing containers and images)..."
            j_docker rm
            j_docker rmi
            ;;
        "ps")
            echo "📋 Docker containers:"
            docker ps -a
            ;;
        "images")
            echo "📋 Docker images:"
            docker images
            ;;
        "help"|"-h"|"--help")
            j_docker_help
            ;;
        *)
            echo "❌ Unknown docker command: $command"
            j_docker_help
            return 1
            ;;
    esac
}

# Docker help function
j_docker_help() {
    echo "🐳 Docker Commands"
    echo ""
    echo "Usage: j docker <command>"
    echo ""
    echo "Commands:"
    echo "  rm       Remove all containers"
    echo "  rmi      Remove all images"
    echo "  clean    Clean up Docker system (prune)"
    echo "  reset    Remove all containers and images"
    echo "  ps       List all containers"
    echo "  images   List all images"
    echo "  help     Show this help"
}

# Auto-completion for docker commands
j_docker_completion() {
    echo "rm rmi clean reset ps images help"
}

# Module metadata
J_MODULE_NAME="docker"
J_MODULE_DESCRIPTION="Docker container and image management"
J_MODULE_COMMANDS="rm rmi clean reset ps images help"
