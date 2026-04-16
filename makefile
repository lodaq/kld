.PHONY: build build-opt install clean help run

APP_NAME := kld
BINARY := $(APP_NAME)

# Colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

help:
	@echo "$(GREEN)KLD Build Targets$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@echo "  $(GREEN)make build$(NC)        - Build the app"
	@echo "  $(GREEN)make build-opt$(NC)    - Build optimized (smaller) binary"
	@echo "  $(GREEN)make install$(NC)      - Build and install to /usr/local/bin"
	@echo "  $(GREEN)make run$(NC)          - Build and run the app"
	@echo "  $(GREEN)make clean$(NC)        - Remove built binaries"
	@echo "  $(GREEN)make help$(NC)         - Show this help message"

build:
	@echo "$(YELLOW)Building $(APP_NAME)...$(NC)"
	go mod download
	go build -o $(BINARY)
	@echo "$(GREEN)✓ Build complete$(NC)"
	@echo "Run with: ./$(BINARY)"

build-opt:
	@echo "$(YELLOW)Building optimized binary...$(NC)"
	go mod download
	go build -ldflags="-s -w" -o $(BINARY)-opt
	@echo "$(GREEN)✓ Build complete$(NC)"
	@ls -lh $(BINARY)-opt

install: build
	@echo "$(YELLOW)Installing to /usr/local/bin/...$(NC)"
	sudo cp $(BINARY) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY)
	@echo "$(GREEN)✓ Installed successfully$(NC)"
	@echo "Run with: $(BINARY)"

run: build
	./$(BINARY)

clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -f $(BINARY)
	rm -f $(BINARY)-opt
	@echo "$(GREEN)✓ Cleaned$(NC)"

deps:
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"