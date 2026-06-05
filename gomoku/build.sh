#!/bin/bash

# 输出目录
BUILD_DIR="bin"

# 入口文件
MAIN_PKG="./cmd/game"

# 创建输出目录
mkdir -p $BUILD_DIR

echo "Cleaning old builds..."
rm -rf $BUILD_DIR/*

echo "==================================="
echo "Building for macOS (Apple Silicon)"
echo "==================================="
GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o $BUILD_DIR/gomoku-mac-arm64 $MAIN_PKG

echo "==================================="
echo "Building for Windows (64-bit)"
echo "==================================="
# 使用 -H=windowsgui 可以在 Windows 上隐藏运行时的黑色命令行窗口
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w -H=windowsgui" -o $BUILD_DIR/gomoku-windows-amd64.exe $MAIN_PKG

echo "==================================="chmod +x build.sh
echo "Build complete! Check the '$BUILD_DIR' directory."
echo "==================================="
