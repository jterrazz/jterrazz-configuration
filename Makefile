.PHONY: build install uninstall check test clean help

# Configuration
BINARY_NAME := j
INSTALL_PATH := /usr/local/bin/$(BINARY_NAME)
ZSHRC_SOURCE := dotfiles/shell/zsh/zshrc.sh
ZSHRC_LINE := source $(PWD)/$(ZSHRC_SOURCE)

help: ## Show this help message
	@echo "jterrazz-cli"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the j binary
	@echo "Building $(BINARY_NAME)..."
	@cd src && go build -o ../$(BINARY_NAME) ./cmd/j
	@echo "✅ Built ./$(BINARY_NAME)"

install: build ## Build and install j to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)
	@sudo chmod +x $(INSTALL_PATH)
	@rm $(BINARY_NAME)
	@echo "Setting up shell completions..."
	@if [ -f "$$HOME/.zshrc" ]; then \
		if ! grep -q "$(ZSHRC_SOURCE)" "$$HOME/.zshrc"; then \
			echo '\n# jterrazz-cli' >> "$$HOME/.zshrc"; \
			echo 'source $(PWD)/$(ZSHRC_SOURCE)' >> "$$HOME/.zshrc"; \
			echo "✅ Added shell completions to ~/.zshrc"; \
		else \
			echo "✅ Shell completions already configured in ~/.zshrc"; \
		fi \
	else \
		echo "⚠️  ~/.zshrc not found. Add this line manually:"; \
		echo "    source $(PWD)/$(ZSHRC_SOURCE)"; \
	fi
	@echo "✅ Installed! Run 'source ~/.zshrc' then 'j help' to get started."

uninstall: ## Remove j from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@if [ -f "$(INSTALL_PATH)" ]; then \
		sudo rm $(INSTALL_PATH); \
		echo "✅ Uninstalled $(BINARY_NAME)"; \
	else \
		echo "⚠️  $(BINARY_NAME) not found at $(INSTALL_PATH)"; \
	fi

check: ## Check if j is installed
	@if command -v $(BINARY_NAME) >/dev/null 2>&1; then \
		echo "✅ $(BINARY_NAME) is installed at $$(which $(BINARY_NAME))"; \
	else \
		echo "❌ $(BINARY_NAME) is not installed"; \
	fi

test: ## Run tests
	@cd src && go test ./...

clean: ## Remove build artifacts
	@rm -f $(BINARY_NAME)
	@echo "✅ Cleaned"
