#!/bin/bash

# jterrazz unified modular command system
# This script auto-discovers and loads command modules from subdirectories

# Get the directory where this script is located
J_ROOT_DIR="$(dirname "$0")"

# Array to store loaded modules
declare -A J_LOADED_MODULES
declare -A J_MODULE_DESCRIPTIONS

# Load all available modules
_j_load_modules() {
    for module_dir in "$J_ROOT_DIR"/*/; do
        if [ -d "$module_dir" ]; then
            local module_name=$(basename "$module_dir")
            local commands_file="$module_dir/commands.sh"
            
            if [ -f "$commands_file" ]; then
                # Source the module
                source "$commands_file"
                
                # Register the module if it has the required metadata
                if declare -f "j_${module_name}" > /dev/null 2>&1; then
                    J_LOADED_MODULES[$module_name]=1
                    
                    # Try to get module description from metadata
                    if [ -n "$J_MODULE_DESCRIPTION" ]; then
                        J_MODULE_DESCRIPTIONS[$module_name]="$J_MODULE_DESCRIPTION"
                    else
                        J_MODULE_DESCRIPTIONS[$module_name]="$module_name commands"
                    fi
                    
                    # Clean up the metadata variable for next module
                    unset J_MODULE_DESCRIPTION
                fi
            fi
        fi
    done
}

# Main j command function
j() {
    # Load modules if not already loaded
    if [ ${#J_LOADED_MODULES[@]} -eq 0 ]; then
        _j_load_modules
    fi

    if [ $# -eq 0 ]; then
        _j_help
        return 1
    fi

    local category="$1"
    shift

    case "$category" in
        "help"|"-h"|"--help")
            _j_help
            ;;
        *)
            # Check if module exists and call it
            if [ "${J_LOADED_MODULES[$category]}" = "1" ]; then
                "j_${category}" "$@"
            else
                echo "❌ Unknown category: $category"
                echo "💡 Run 'j help' to see available commands"
                _j_list_modules
                return 1
            fi
            ;;
    esac
}

# Help function
_j_help() {
    echo "🚀 jterrazz unified command system"
    echo ""
    echo "Usage: j <category> <command> [args...]"
    echo ""
    echo "Categories:"
    
    # List all loaded modules
    for module in ${(k)J_LOADED_MODULES}; do
        local description="${J_MODULE_DESCRIPTIONS[$module]}"
        printf "  %-10s %s\n" "$module" "$description"
    done
    
    echo ""
    echo "Examples:"
    echo "  j docker rm              # Remove all containers"
    echo "  j git feat 'new feature' # Add and commit with feat prefix"
    echo ""
    echo "For category-specific help:"
    for module in ${(k)J_LOADED_MODULES}; do
        echo "  j $module help"
    done
}

# List available modules (used in error messages)
_j_list_modules() {
    echo ""
    echo "Available categories:"
    for module in ${(k)J_LOADED_MODULES}; do
        echo "  - $module"
    done
}

# Completion function for zsh
_j_completion() {
    local state
    _arguments \
        '1: :->category' \
        '*: :->command'
    
    case $state in
        category)
            local modules=()
            for module in ${(k)J_LOADED_MODULES}; do
                modules+=("$module")
            done
            modules+=("help")
            _values 'category' "${modules[@]}"
            ;;
        command)
            # Try to get completion from the specific module
            local category="$words[2]"
            if [ "${J_LOADED_MODULES[$category]}" = "1" ]; then
                if declare -f "j_${category}_completion" > /dev/null 2>&1; then
                    local commands=$(j_${category}_completion)
                    _values "${category} command" ${=commands}
                fi
            fi
            ;;
    esac
}

# Register completion if in zsh
if [ -n "$ZSH_VERSION" ]; then
    compdef _j_completion j
fi

# Initialize modules on script load
_j_load_modules