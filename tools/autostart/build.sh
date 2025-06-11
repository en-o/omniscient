#!/bin/bash

# 项目自启管理工具构建脚本

set -e

PROJECT_NAME="autostart"
VERSION="0.0.1"

echo "Building $PROJECT_NAME v$VERSION..."

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# 创建构建目录结构
BUILD_DIR="build/$VERSION"
mkdir -p "$BUILD_DIR/amd64"
mkdir -p "$BUILD_DIR/arm64"

# 构建 AMD64 版本
echo "Compiling for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "$BUILD_DIR/amd64/$PROJECT_NAME" main.go

# 构建 ARM64 版本
echo "Compiling for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o "$BUILD_DIR/arm64/$PROJECT_NAME" main.go

# 检查构建结果
if [ -f "$BUILD_DIR/amd64/$PROJECT_NAME" ] && [ -f "$BUILD_DIR/arm64/$PROJECT_NAME" ]; then
    echo "✓ Build successful!"
    echo "Binary locations:"
    echo "- AMD64: $BUILD_DIR/amd64/$PROJECT_NAME"
    echo "- ARM64: $BUILD_DIR/arm64/$PROJECT_NAME"

    # 显示文件信息
    echo -e "\nAMD64 binary:"
    ls -lh "$BUILD_DIR/amd64/$PROJECT_NAME"
    echo -e "\nARM64 binary:"
    ls -lh "$BUILD_DIR/arm64/$PROJECT_NAME"

    echo -e "\nInstallation options:"
    echo "For AMD64 systems:"
    echo "  sudo cp $BUILD_DIR/amd64/$PROJECT_NAME /usr/local/bin/"
    echo "For ARM64 systems:"
    echo "  sudo cp $BUILD_DIR/arm64/$PROJECT_NAME /usr/local/bin/"
else
    echo "✗ Build failed!"
    exit 1
fi