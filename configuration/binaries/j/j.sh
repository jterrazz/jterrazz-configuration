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
        "status")
            _j_status
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

# Status function
_j_status() {
    echo "📊 jterrazz command system status"
    echo ""
    
    # System Information
    echo "🖥️  System Information:"
    echo "  OS: $(uname -s) $(uname -r)"
    echo "  User: $(whoami)"
    echo "  Shell: $SHELL"
    echo "  Hostname: $(hostname)"
    echo ""
    
    # Development Tools
    echo "🛠️  Development Tools:"
    
    # Homebrew
    if command -v brew >/dev/null 2>&1; then
        echo "  ✅ Homebrew: $(brew --version | head -1)"
    else
        echo "  ❌ Homebrew: Not installed"
    fi
    
    # Node.js and npm
    if command -v node >/dev/null 2>&1; then
        echo "  ✅ Node.js: $(node --version)"
    else
        echo "  ❌ Node.js: Not installed"
    fi
    
    if command -v npm >/dev/null 2>&1; then
        echo "  ✅ npm: $(npm --version)"
    else
        echo "  ❌ npm: Not installed"
    fi
    
    # NVM
    if command -v nvm >/dev/null 2>&1; then
        echo "  ✅ NVM: $(nvm --version)"
    else
        echo "  ❌ NVM: Not installed"
    fi
    
    # Git
    if command -v git >/dev/null 2>&1; then
        echo "  ✅ Git: $(git --version | cut -d' ' -f3)"
    else
        echo "  ❌ Git: Not installed"
    fi
    
    # Docker
    if command -v docker >/dev/null 2>&1; then
        echo "  ✅ Docker: $(docker --version | cut -d' ' -f3 | sed 's/,//')"
    else
        echo "  ❌ Docker: Not installed"
    fi
    
    echo ""
    
    # Configuration Status
    echo "⚙️  Configuration Status:"
    
    # Oh My Zsh
    if [ -d "$HOME/.oh-my-zsh" ]; then
        echo "  ✅ Oh My Zsh: Installed"
    else
        echo "  ❌ Oh My Zsh: Not installed"
    fi
    
    # jterrazz configuration
    if grep -q "source ~/Developer/jterrazz-configuration/configuration/binaries/zsh/zshrc.sh" ~/.zshrc 2>/dev/null; then
        echo "  ✅ jterrazz config: Loaded in shell"
    else
        echo "  ❌ jterrazz config: Not loaded in shell"
    fi
    
    # SSH key
    if [ -f "$HOME/.ssh/id_github" ]; then
        echo "  ✅ GitHub SSH key: Configured"
    else
        echo "  ❌ GitHub SSH key: Not configured"
    fi
    
    echo ""
    
    # Development Packages
    echo "📦 Development Packages:"
    
    # ansible-lint
    if command -v ansible-lint >/dev/null 2>&1; then
        echo "  ✅ ansible-lint: $(ansible-lint --version | head -1)"
    else
        echo "  ❌ ansible-lint: Not installed"
    fi
    
    # ansible
    if command -v ansible >/dev/null 2>&1; then
        echo "  ✅ ansible: $(ansible --version | head -1 | cut -d' ' -f3)"
    else
        echo "  ❌ ansible: Not installed"
    fi
    
    # terraform
    if command -v terraform >/dev/null 2>&1; then
        echo "  ✅ terraform: $(terraform version | head -1 | cut -d' ' -f2)"
    else
        echo "  ❌ terraform: Not installed"
    fi
    
    # kubectl
    if command -v kubectl >/dev/null 2>&1; then
        echo "  ✅ kubectl: $(kubectl version --client 2>/dev/null | head -1 | cut -d' ' -f3)"
    else
        echo "  ❌ kubectl: Not installed"
    fi
    
    # multipass
    if command -v multipass >/dev/null 2>&1; then
        echo "  ✅ multipass: $(multipass version | head -1 | awk '{print $NF}')"
    else
        echo "  ❌ multipass: Not installed"
    fi
    
    # biome
    if command -v biome >/dev/null 2>&1; then
        echo "  ✅ biome: $(biome --version)"
    else
        echo "  ❌ biome: Not installed"
    fi
    
    # bun
    if command -v bun >/dev/null 2>&1; then
        echo "  ✅ bun: $(bun --version)"
    else
        echo "  ❌ bun: Not installed"
    fi
    
    # python
    if command -v python3 >/dev/null 2>&1; then
        echo "  ✅ python: $(python3 --version | cut -d' ' -f2)"
    else
        echo "  ❌ python: Not installed"
    fi
    
    # neohtop
    if brew list --cask neohtop >/dev/null 2>&1; then
        echo "  ✅ neohtop: $(brew list --cask --versions neohtop 2>/dev/null || echo "installed")"
    else
        echo "  ❌ neohtop: Not installed"
    fi
    
    # codex
    if command -v codex >/dev/null 2>&1; then
        echo "  ✅ codex: $(codex --version 2>/dev/null | head -1 || echo "installed")"
    else
        echo "  ❌ codex: Not installed"
    fi
    
    # claude
    if command -v claude >/dev/null 2>&1; then
        echo "  ✅ claude: $(claude --version 2>/dev/null | head -1 || echo "installed")"
    else
        echo "  ❌ claude: Not installed"
    fi
    
    echo ""
    
    # Loaded Modules
    echo "🔧 Loaded Command Modules:"
    for module in ${(k)J_LOADED_MODULES}; do
        local description="${J_MODULE_DESCRIPTIONS[$module]}"
        printf "  ✅ %-10s %s\n" "$module" "$description"
    done
}

# Help function
_j_help() {
    echo "🚀 jterrazz unified command system"
    echo ""
    echo "Usage: j <category> <command> [args...]"
    echo ""
    echo "Global Commands:"
    echo "  status    Show comprehensive system status"
    echo "  help      Show this help"
    echo ""
    echo "Categories:"
    
    # List all loaded modules
    for module in ${(k)J_LOADED_MODULES}; do
        local description="${J_MODULE_DESCRIPTIONS[$module]}"
        printf "  %-10s %s\n" "$module" "$description"
    done
    
    echo ""
    echo "Examples:"
    echo "  j status                 # Show system status"
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
    local state line
    local -A opt_args
    
    _arguments -C \
        '1: :->category' \
        '2: :->command' \
        '3: :->subcommand' \
        '*: :->args'
    
    case $state in
        category)
            local modules=()
            for module in ${(k)J_LOADED_MODULES}; do
                modules+=("$module")
            done
            modules+=("help" "status")
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
        subcommand)
            # Handle nested subcommands (e.g., j system install <subcommand>)
            local category="$words[2]"
            local command="$words[3]"
            
            if [ "${J_LOADED_MODULES[$category]}" = "1" ]; then
                # Check if module has nested completion function
                local nested_func="j_${category}_${command}_completion"
                if declare -f "$nested_func" > /dev/null 2>&1; then
                    local subcommands=$($nested_func)
                    _values "${category} ${command} subcommand" ${=subcommands}
                fi
            fi
            ;;
        args)
            # Future: handle command arguments
            ;;
    esac
}

# Register completion if in zsh and compdef is available
if [ -n "$ZSH_VERSION" ] && command -v compdef > /dev/null 2>&1; then
    compdef _j_completion j
fi

# Initialize modules on script load
_j_load_modules