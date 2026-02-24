# jterrazz shell configuration
# This file is sourced by ~/.zshrc

# Bun global binaries
export PATH="$HOME/.bun/bin:$PATH"

# Load j command completions
if command -v j &> /dev/null; then
    eval "$(j completion zsh)"
fi

# Tmux launcher commands
jt() {
    if ! command -v tmux &>/dev/null; then
        echo "tmux not found"
        return 1
    fi

    if tmux has-session -t main 2>/dev/null; then
        tmux new-session -t main \; set destroy-unattached on \; new-window
    else
        tmux new-session -s main
    fi
}

_jt_tool() {
    local window_name="$1"
    local tool_cmd="$2"

    if ! command -v tmux &>/dev/null; then
        echo "tmux not found"
        return 1
    fi
    if ! command -v "$tool_cmd" &>/dev/null; then
        echo "$tool_cmd not found"
        return 1
    fi

    if [[ -n "$TMUX" ]]; then
        if tmux has-session -t main 2>/dev/null; then
            tmux new-window -t main -n "$window_name" "$tool_cmd"
            tmux switch-client -t main
        else
            tmux new-session -ds main -n "$window_name" "$tool_cmd"
            tmux switch-client -t main
        fi
        return
    fi

    if tmux has-session -t main 2>/dev/null; then
        tmux new-session -t main \; set destroy-unattached on \; new-window -n "$window_name" "$tool_cmd"
    else
        tmux new-session -s main -n "$window_name" "$tool_cmd"
    fi
}

jc() { _jt_tool "claude" "claude"; }
jo() { _jt_tool "codex" "codex"; }
tg() { _jt_tool "gemini" "gemini"; }
