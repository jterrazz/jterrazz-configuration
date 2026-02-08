# jterrazz shell configuration
# This file is sourced by ~/.zshrc

# Bun global binaries
export PATH="$HOME/.bun/bin:$PATH"

# Load j command completions
if command -v j &> /dev/null; then
    eval "$(j completion zsh)"
fi
