.PHONY: install uninstall help check

# Configuration
ZSHRC_PATH := $(HOME)/.zshrc
SOURCE_LINE := source $(PWD)/configuration/binaries/zsh/zshrc.sh

help: ## Show this help message
	@echo "Jterrazz Configuration - Installation"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install by adding source line to ~/.zshrc
	@echo "Installing jterrazz-configuration..."
	@if [ ! -f "$(ZSHRC_PATH)" ]; then \
		echo "Creating ~/.zshrc..."; \
		touch "$(ZSHRC_PATH)"; \
	fi
	@if ! grep -Fxq "$(SOURCE_LINE)" "$(ZSHRC_PATH)"; then \
		echo "Adding source line to ~/.zshrc..."; \
		echo "" >> "$(ZSHRC_PATH)"; \
		echo "# jterrazz-configuration" >> "$(ZSHRC_PATH)"; \
		echo "$(SOURCE_LINE)" >> "$(ZSHRC_PATH)"; \
		echo "✅ Installation complete! Please restart your terminal or run: source ~/.zshrc"; \
	else \
		echo "✅ Already installed - source line exists in ~/.zshrc"; \
	fi

uninstall: ## Remove source line from ~/.zshrc
	@echo "Uninstalling jterrazz-configuration..."
	@if [ -f "$(ZSHRC_PATH)" ]; then \
		if grep -Fxq "$(SOURCE_LINE)" "$(ZSHRC_PATH)"; then \
			echo "Removing source line from ~/.zshrc..."; \
			sed -i.bak '/# jterrazz-configuration/d' "$(ZSHRC_PATH)"; \
			sed -i.bak '\|$(SOURCE_LINE)|d' "$(ZSHRC_PATH)"; \
			rm "$(ZSHRC_PATH).bak"; \
			echo "✅ Uninstallation complete! Please restart your terminal."; \
		else \
			echo "⚠️  Source line not found in ~/.zshrc - nothing to remove"; \
		fi; \
	else \
		echo "⚠️  ~/.zshrc not found - nothing to remove"; \
	fi

check: ## Check if the package is currently installed
	@echo "Checking installation status..."
	@if [ -f "$(ZSHRC_PATH)" ] && grep -Fxq "$(SOURCE_LINE)" "$(ZSHRC_PATH)"; then \
		echo "✅ jterrazz-configuration is installed"; \
	else \
		echo "❌ jterrazz-configuration is not installed"; \
	fi
