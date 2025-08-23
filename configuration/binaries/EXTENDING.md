# Extending the Modular `j` Command System

## Prerequisites

Ensure the jterrazz-configuration system is installed:

```bash
# From project root
make install
# OR manually: source ~/Developer/jterrazz-configuration/configuration/binaries/zsh/zshrc.sh
```

## Adding a New Module

The modular system automatically discovers new modules. To add a new category (e.g., `j mymodule`), simply create a new module directory and commands file.

### Step 1: Create Module Directory

```bash
# From project root: ~/Developer/jterrazz-configuration/
mkdir configuration/binaries/j/mymodule
```

### Step 2: Create Commands File

Create `configuration/binaries/j/mymodule/commands.sh`:

```bash
#!/bin/bash

# MyModule module for jterrazz command system
# This file defines all mymodule-related commands

# Main mymodule command handler
j_mymodule() {
    if [ $# -eq 0 ]; then
        j_mymodule_help
        return 1
    fi

    local command="$1"
    shift

    case "$command" in
        "hello")
            echo "👋 Hello from mymodule!"
            ;;
        "world")
            echo "🌍 World command executed!"
            ;;
        "help"|"-h"|"--help")
            j_mymodule_help
            ;;
        *)
            echo "❌ Unknown mymodule command: $command"
            j_mymodule_help
            return 1
            ;;
    esac
}

# MyModule help function
j_mymodule_help() {
    echo "🔧 MyModule Commands"
    echo ""
    echo "Usage: j mymodule <command>"
    echo ""
    echo "Commands:"
    echo "  hello     Say hello"
    echo "  world     Execute world command"
    echo "  help      Show this help"
}

# Auto-completion for mymodule commands
j_mymodule_completion() {
    echo "hello world help"
}

# Module metadata
J_MODULE_NAME="mymodule"
J_MODULE_DESCRIPTION="My custom module"
J_MODULE_COMMANDS="hello world help"
```

### Step 3: Make Executable

```bash
chmod +x configuration/binaries/j/mymodule/commands.sh
```

### Step 4: Test

The module auto-loads on next use:

```bash
j help              # Should show mymodule
j mymodule help     # Show mymodule commands
j mymodule hello    # Test the command
```

## Module Structure Requirements

### Required Functions

1. **Main Handler**: `j_<modulename>()`

   - Entry point for your module
   - Must handle help and error cases

2. **Help Function**: `j_<modulename>_help()`
   - Display usage and commands
   - Follow the emoji + description format

### Optional Functions

3. **Completion**: `j_<modulename>_completion()`

   - Return space-separated command list
   - Used for shell tab completion

4. **Nested Completion**: `j_<modulename>_<command>_completion()`
   - For commands with subcommands (like `j system install <subcommand>`)
   - Return space-separated subcommand list
   - Automatically detected and used

### Required Metadata

```bash
J_MODULE_NAME="modulename"
J_MODULE_DESCRIPTION="Short description"
J_MODULE_COMMANDS="cmd1 cmd2 help"
```

## Best Practices

### Naming Conventions

- **Module directory**: lowercase, no spaces
- **Function names**: `j_<modulename>_<function>`
- **Variables**: `J_MODULE_*` for metadata

### Error Handling

```bash
# Check for required parameters
if [ $# -eq 0 ]; then
    echo "❌ Please provide a parameter"
    echo "💡 Usage: j mymodule command 'parameter'"
    return 1
fi

# Check for external dependencies
if ! command -v some_tool >/dev/null 2>&1; then
    echo "❌ some_tool not found"
    return 1
fi
```

### User Feedback

- Use emojis for visual clarity
- Provide clear error messages with usage hints
- Show progress for long-running operations
- Use consistent formatting

### Example Commands

```bash
case "$command" in
    "start")
        echo "🚀 Starting service..."
        # your code here
        ;;
    "stop")
        echo "🛑 Stopping service..."
        # your code here
        ;;
    "status")
        echo "📊 Service status:"
        # your code here
        ;;
esac
```

### Autocomplete Examples

**Basic completion:**

```bash
j_mymodule_completion() {
    echo "hello world help"
}
```

**Nested completion (for subcommands):**

```bash
# For: j mymodule install <subcommand>
j_mymodule_install_completion() {
    echo "package1 package2 all help"
}
```

## Real Examples

Study the existing modules for implementation patterns:

**`configuration/binaries/j/system/commands.sh`** - Complex module with:

- Parameter validation and external dependency checking
- Nested subcommands (`j system install brew`)
- Multi-level autocomplete support
- Idempotent installation logic

**`configuration/binaries/j/git/commands.sh`** - Simple module with:

- Streamlined git workflow commands
- Parameter handling for commit messages
- Basic command structure

**`configuration/binaries/j/docker/commands.sh`** - Clean utility module with:

- System cleanup and management commands
- Clear command organization
- Consistent error handling

## Testing Your Module

```bash
# Reload shell or source the zsh integration
source ~/Developer/jterrazz-configuration/configuration/binaries/zsh/zshrc.sh

# Test commands
j help                  # Should show your module
j mymodule help         # Show your module's commands
j mymodule hello        # Test specific commands

# Test autocomplete
j <TAB>                 # Should include your module
j mymodule <TAB>        # Should show your commands
```

## Module Ideas

- **dev**: Development workflow (build, test, deploy)
- **net**: Network utilities (ping, scan, etc.)
- **file**: File operations (compress, backup, etc.)
- **config**: Configuration management
- **media**: Media conversion and processing
- **cloud**: Cloud service management

The modular design makes it easy to organize and maintain your custom commands!
