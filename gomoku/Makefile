.PHONY: all build-mac build-windows clean

# 输出目录
BUILD_DIR = bin

# 游戏入口文件
MAIN_PKG = ./cmd/game

all: clean build-mac build-windows

build-mac:
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/gomoku-mac-arm64 $(MAIN_PKG)
	@echo "macOS builds completed."

build-windows:
	@echo "Building for Windows (64-bit)..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/gomoku-windows-amd64.exe $(MAIN_PKG)
	@echo "Windows builds completed."

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
