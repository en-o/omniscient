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

# 构建二进制文件
echo "Compiling..."
go build -ldflags "-s -w" -o $BUILD_DIR/$PROJECT_NAME main.go

# 检查构建结果
if [ -f "$BUILD_DIR/$PROJECT_NAME" ]; then
    echo "✓ Build successful!"
    echo "Binary location: $BUILD_DIR/$PROJECT_NAME"

    # 显示文件信息
    ls -lh $BUILD_DIR/$PROJECT_NAME

    echo ""
    echo "Installation options:"
    echo "1. Use directly: ./$BUILD_DIR/$PROJECT_NAME"
    echo "2. Install to system: sudo cp $BUILD_DIR/$PROJECT_NAME /usr/local/bin/"
    echo "3. Install to /usr/bin: sudo cp $BUILD_DIR/$PROJECT_NAME /usr/bin/"

    echo ""
    echo "Usage examples:"
    echo "  ./$BUILD_DIR/$PROJECT_NAME help"
    echo "  ./$BUILD_DIR/$PROJECT_NAME list"
    echo "  sudo ./$BUILD_DIR/$PROJECT_NAME add myapp jar /path/to/app.jar"
else
    echo "✗ Build failed!"
    exit 1
fi