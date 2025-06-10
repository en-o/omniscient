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

# 创建构建目录
BUILD_DIR="build"
mkdir -p $BUILD_DIR

# 构建 AMD64 版本
echo "Compiling for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $BUILD_DIR/${PROJECT_NAME}_amd64 main.go

# 构建 ARM64 版本
echo "Compiling for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $BUILD_DIR/${PROJECT_NAME}_arm64 main.go

# 检查构建结果
if [ -f "$BUILD_DIR/${PROJECT_NAME}_amd64" ] && [ -f "$BUILD_DIR/${PROJECT_NAME}_arm64" ]; then
    echo "✓ Build successful!"
    echo "Binary locations:"
    echo "- AMD64: $BUILD_DIR/${PROJECT_NAME}_amd64"
    echo "- ARM64: $BUILD_DIR/${PROJECT_NAME}_arm64"

    # 显示文件信息
    ls -lh $BUILD_DIR/${PROJECT_NAME}_*

    echo ""
    echo "Installation options:"
    echo "For AMD64 systems:"
    echo "  sudo cp $BUILD_DIR/${PROJECT_NAME}_amd64 /usr/local/bin/$PROJECT_NAME"
    echo "For ARM64 systems:"
    echo "  sudo cp $BUILD_DIR/${PROJECT_NAME}_arm64 /usr/local/bin/$PROJECT_NAME"
else
    echo "✗ Build failed!"
    exit 1
fi