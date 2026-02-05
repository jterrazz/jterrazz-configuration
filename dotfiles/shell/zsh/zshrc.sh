# jterrazz shell configuration
# This file is sourced by ~/.zshrc

# Load j command completions
if command -v j &> /dev/null; then
    eval "$(j completion zsh)"
fi
